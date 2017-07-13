package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
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
	reader    CSVReader
	writer    *csv.Writer
	filename  string
	delimiter rune
	numLines  uint
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
	f.reader = NewCSVReader(file, f.delimiter)
	f.r = file
	return nil
}

func (f *csvFile) openFileOrStdin(filename string) (io.ReadCloser, error) {
	if filename == stdPipe {
		return os.Stdin, nil
	}
	return os.Open(filename)
}

func (f *csvFile) Write(row []string) error {
	if f.writer == nil {
		if err := f.createWriter(); err != nil {
			return err
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
