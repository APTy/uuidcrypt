package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

const stdPipe = "-"

type File interface {
	Read() ([]string, error)
	Write([]string) error
	Close() error
}

type CSVOptions func(*csvFile)

func WithDelimiter(delimiter string) CSVOptions {
	runeDelimiter, err := parseStringToRune(delimiter)
	if err != nil {
		runeDelimiter = defaultDelimiter
	}
	return func(f *csvFile) {
		f.delimiter = runeDelimiter
	}
}

func WithFixQuotes(fixQuotes bool) CSVOptions {
	return func(f *csvFile) {
		f.fixQuotes = fixQuotes
	}
}

const defaultDelimiter = ','

func NewCSVFile(filename string, options ...CSVOptions) File {
	f := &csvFile{
		filename:  filename,
		delimiter: defaultDelimiter,
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

type csvFile struct {
	r         io.ReadCloser
	w         io.WriteCloser
	bw        *bufio.Writer
	reader    *csv.Reader
	writer    *csv.Writer
	filename  string
	delimiter rune
	numLines  uint
	fixQuotes bool
}

func (f *csvFile) Read() ([]string, error) {
	if f.reader == nil {
		if err := f.createReader(); err != nil {
			return nil, err
		}
	}
	row, err := f.reader.Read()
	f.numLines++
	return row, err
}

func (f *csvFile) createReader() error {
	if f.writer != nil {
		return fmt.Errorf("file object is already a writer")
	}
	file, err := f.openFileOrStdin(f.filename)
	if err != nil {
		return err
	}
	f.reader = csv.NewReader(file)
	f.reader.Comma = f.delimiter
	f.r = file
	return nil
}

func (f *csvFile) openFileOrStdin(filename string) (io.ReadCloser, error) {
	var file io.ReadCloser = os.Stdin
	if filename != stdPipe {
		var err error
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}
	if f.fixQuotes {
		return readerWithUnicodeQuotes(file), nil
	}
	return file, nil
}

func readerWithUnicodeQuotes(file io.ReadCloser) io.ReadCloser {
	r, w := io.Pipe()
	br := bufio.NewReader(file)
	bw := bufio.NewWriter(w)
	go func() {
		defer file.Close()
		defer w.Close()
		// replace the standard ascii quote rune with the unicode quote
		streamReplace(w, br, bw, '"', '\u201E')
	}()
	return r
}

func streamReplace(w *io.PipeWriter, br *bufio.Reader, bw *bufio.Writer, find, replace rune) {
	for {
		r1, _, err := br.ReadRune()
		if err != nil {
			bw.Flush()
			w.CloseWithError(err)
			return
		}
		if r1 == find {
			r1 = replace
		}
		if _, err := bw.WriteRune(r1); err != nil {
			bw.Flush()
			w.CloseWithError(err)
			return
		}
	}
}

func (f *csvFile) Write(row []string) error {
	if f.writer == nil {
		if err := f.createWriter(); err != nil {
			return err
		}
	}
	if f.fixQuotes {
		for i := range row {
			// replace the unicode quote rune with the standard ascii quote
			row[i] = strings.Replace(row[i], "\u201E", "\"", -1)
		}
	}
	return f.writer.Write(row)
}

func (f *csvFile) createWriter() error {
	if f.reader != nil {
		return fmt.Errorf("file object is already a reader")
	}
	file, err := createFileOrStdout(f.filename)
	if err != nil {
		return fmt.Errorf("file create error: %v", err)
	}
	f.writer = csv.NewWriter(file)
	f.writer.Comma = f.delimiter
	f.w = file
	return nil
}

func createFileOrStdout(filename string) (io.WriteCloser, error) {
	if filename == stdPipe {
		return os.Stdout, nil
	}
	return os.Create(filename)
}

func (f *csvFile) Close() error {
	if f.writer != nil {
		f.writer.Flush()
	}
	if f.w != nil {
		return f.w.Close()
	}
	if f.r != nil {
		return f.r.Close()
	}
	return nil
}
