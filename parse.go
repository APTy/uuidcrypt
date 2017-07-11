package main

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

var (
	ErrParseEmpty = errors.New("parse: empty string")
)

func parseStringToRune(str string) (rune, error) {
	delim, err := strconv.Unquote("'" + str + "'")
	if err != nil {
		return 0, err
	}
	if delim == "" {
		return 0, ErrParseEmpty
	}
	r, _ := utf8.DecodeRuneInString(delim)
	return r, nil
}
