// https://exiftool.org
// https://github.com/barasher/go-exiftool
// uses exiftool's stay_open feature to optimize performance

package exiftool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"myPhotos/logger"
	"myPhotos/tools"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var Et *Exiftool

func init() {
	et, err := NewExiftool()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
	Et = et
}

var exiftoolBinary = "exiftool"
var executeArg = "-execute"
var args = []string{"-stay_open", "True", "-@", "-",
	"-common_args", "-fast1", "-charset", "filename=utf8",
	"-FileType", "-FileSize", "-ImageSize", "-Make", "-Model", "-Megapixels", "-GPSLongitude#", "-GPSLatitude#",
	"-ISO", "-Flash", "-FocalLength", "-ShutterSpeed", "-Aperture", "-MediaGroupUUID", "-DateTimeOriginal",
	"-Duration#", "-ContentIdentifier", "-CreationDate", "-CreateDate", "-j"}

var closeArgs = []string{"-stay_open", "False", executeArg}

var waitTimeout = time.Second
var readyToken = []byte("{ready}\r\n")
var readyTokenLen = len(readyToken)

type Exiftool struct {
	lock          sync.Mutex
	stdin         io.WriteCloser
	stdMergedOut  io.ReadCloser
	scanMergedOut *bufio.Scanner
	cmd           *exec.Cmd
}

type FileMetadata struct {
	Fields map[string]interface{}
	Err    error
}

func NewExiftool() (*Exiftool, error) {
	e := Exiftool{}

	e.cmd = exec.Command(exiftoolBinary, args...)
	r, w := io.Pipe()
	e.stdMergedOut = r

	e.cmd.Stdout = w
	e.cmd.Stderr = w

	var err error
	if e.stdin, err = e.cmd.StdinPipe(); err != nil {
		return nil, fmt.Errorf("error when piping stdin: %w", err)
	}

	e.scanMergedOut = bufio.NewScanner(r)
	e.scanMergedOut.Split(splitReadyToken)

	if err = e.cmd.Start(); err != nil {
		return nil, fmt.Errorf("error when executing command: %w", err)
	}

	return &e, nil
}

// ErrBufferTooSmall is a sentinel error that is returned when the buffer used to store Exiftool's output is too small.
var ErrBufferTooSmall = errors.New("exiftool's buffer too small (see Buffer init option)")

func (e *Exiftool) ExtractMetadata(file string) FileMetadata {
	e.lock.Lock()
	defer e.lock.Unlock()

	fm := FileMetadata{}

	if _, err := fmt.Fprintln(e.stdin, file); err != nil {
		fm.Err = err
		return fm
	}
	if _, err := fmt.Fprintln(e.stdin, executeArg); err != nil {
		fm.Err = err
		return fm
	}

	scanOk := e.scanMergedOut.Scan()
	scanErr := e.scanMergedOut.Err()
	if scanErr != nil {
		if scanErr == bufio.ErrTooLong {
			fm.Err = ErrBufferTooSmall
			return fm
		}
		fm.Err = fmt.Errorf("error while reading stdMergedOut: %w", e.scanMergedOut.Err())
		return fm
	}
	if !scanOk {
		fm.Err = fmt.Errorf("error while reading stdMergedOut: EOF")
		return fm
	}

	var m []map[string]interface{}
	if err := json.Unmarshal(e.scanMergedOut.Bytes(), &m); err != nil {
		fm.Err = fmt.Errorf("error during unmarshaling (%v): %w)", string(e.scanMergedOut.Bytes()), err)
		return fm
	}

	fm.Fields = m[0]

	return fm
}

func (e *Exiftool) Close() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, v := range closeArgs {
		_, err := fmt.Fprintln(e.stdin, v)
		if err != nil {
			return err
		}
	}

	var errs []error
	if err := e.stdMergedOut.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error while closing stdMergedOut: %w", err))
	}

	if err := e.stdin.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error while closing stdin: %w", err))
	}

	ch := make(chan struct{})
	go func() {
		if e.cmd != nil {
			if err := e.cmd.Wait(); err != nil {
				errs = append(errs, fmt.Errorf("error while waiting for exiftool to exit: %w", err))
			}
		}
		ch <- struct{}{}
		close(ch)
	}()

	// Wait for wait to finish or timeout
	select {
	case <-ch:
	case <-time.After(waitTimeout):
		errs = append(errs, errors.New("Timed out waiting for exiftool to exit"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("error while closing exiftool: %v", errs)
	}

	return nil
}

func splitReadyToken(data []byte, atEOF bool) (int, []byte, error) {
	idx := bytes.Index(data, readyToken)
	if idx == -1 {
		if atEOF && len(data) > 0 {
			return 0, data, fmt.Errorf("no final token found")
		}

		return 0, nil, nil
	}

	return idx + readyTokenLen, data[:idx], nil
}

func (fm FileMetadata) GetString(k string) string {
	v, found := fm.Fields[k]
	if !found || v == nil {
		return ""
	}

	switch v.(type) {
	case string:
		return v.(string)
	default:
		return tools.NumberToString(v)
	}
}

// GetFloat returns a field value as float64 and an error if one occurred.
// KeyNotFoundError will be returned if the key can't be found.
func (fm FileMetadata) GetFloat(k string) float64 {
	v, found := fm.Fields[k]
	if !found || v == nil {
		return 0
	}

	switch v := v.(type) {
	case string:
		return toFloatFallback(v)
	case float64:
		return v
	case int64:
		return float64(v)
	default:
		str := fmt.Sprintf("%v", v)
		return toFloatFallback(str)
	}
}

func toFloatFallback(str string) float64 {
	f, err := strconv.ParseFloat(str, -1)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return 0
	}

	return f
}
