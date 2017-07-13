package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrBadColumnLength = errors.New("csv: bad column length")
)

// CSVReader encapsulates reading of CSV files to better customize
// how parsing occurs.
type CSVReader interface {
	Read() ([]string, error)
}

// NewCSVReader returns a CSVReader using a combination of the
// encoding/csv reader and the strings.split() method.
//
// It supports double quotes that are escaped by placing two
// double-quotes next to each other.
// 	e.g. `"i am ""tyler"""` is interpreted as `i am "tyler"`
//
// It also supports double quotes that are escaped by placing
// a backslash before the double quote character.
// 	e.g. `"i am \"tyler\""` is interpreted as `i am "tyler"`
func NewCSVReader(r io.Reader, delimiter rune) CSVReader {
	return &csvReader{
		r:         bufio.NewScanner(r),
		delimiter: delimiter,
	}
}

type csvReader struct {
	r          *bufio.Scanner
	delimiter  rune
	numColumns uint
}

func (r *csvReader) Read() ([]string, error) {
	if !r.r.Scan() {
		return nil, io.EOF
	}
	line := r.r.Text()
	columns := strings.Split(line, string(r.delimiter))
	if err := r.validateNumColumns(len(columns)); err != nil {
		return nil, err
	}
	for i := range columns {
		columns[i] = r.unquoteIfNeeded(columns[i])
	}
	return columns, nil
}

func (r *csvReader) validateNumColumns(lenColumns int) error {
	numColumns := uint(lenColumns)
	if r.numColumns == 0 {
		r.numColumns = numColumns
	}
	if r.numColumns != numColumns {
		return ErrBadColumnLength
	}
	return nil
}

// remove surrounding double-quotes from the string if there are any.
// this will also un-escape any double-quote characters.
func (r *csvReader) unquoteIfNeeded(str string) string {
	if !isEnclosedInDoubleQuotes(str) {
		return str
	}
	newStr, err := unquoteDoubleQuotes(strings.NewReader(str))
	if err == nil {
		return newStr
	}
	newStr, err = strconv.Unquote(str)
	if err != nil {
		return str
	}
	return newStr
}

func isEnclosedInDoubleQuotes(str string) bool {
	r := strings.NewReader(str)
	return isFirstCharDoubleQuote(r) && isLastCharDoubleQuote(r)
}

func isFirstCharDoubleQuote(r *strings.Reader) bool {
	fst, _, err := r.ReadRune()
	return err == nil && fst == '"'
}

func isLastCharDoubleQuote(r *strings.Reader) bool {
	if _, err := r.Seek(-1, io.SeekEnd); err != nil {
		return false
	}
	lst, _, err := r.ReadRune()
	return err == nil && lst == '"'
}

// use encoding/csv to unquote double double-quotes because its good at that
func unquoteDoubleQuotes(r *strings.Reader) (string, error) {
	cr := csv.NewReader(r)
	rec, err := cr.Read()
	if err != nil {
		return "", err
	}
	return rec[0], nil
}
