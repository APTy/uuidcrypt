package main

import (
	"io"

	"github.com/google/uuid"
)

type UUIDCrypt interface {
	Run() error
}

type UUIDCryptOptions func(*uuidCrypt)

func WithColumns(columns ...int) UUIDCryptOptions {
	return func(u *uuidCrypt) {
		u.columns = columns
	}
}

func NewUUIDCrypt(input File, output File, processor Processor, options ...UUIDCryptOptions) UUIDCrypt {
	u := &uuidCrypt{
		input:       input,
		output:      output,
		processor:   processor,
		columns:     []int{1},
		headerError: false,
	}
	for _, opt := range options {
		opt(u)
	}
	return u
}

type uuidCrypt struct {
	input       File
	output      File
	processor   Processor
	columns     []int
	headerError bool
}

func (u *uuidCrypt) Run() error {
	defer u.input.Close()
	defer u.output.Close()
	for {
		if err := u.runOnce(); err != nil {
			return errIfNotEOF(err)
		}
	}
	return nil
}

func (u *uuidCrypt) runOnce() error {
	record, err := u.input.Read()
	if err != nil {
		return err
	}
	var rowErr error
	for _, column := range u.columns {
		col := column - 1
		newUUID, err := u.processUUID(record[col])
		if err != nil {
			rowErr = err
			continue
		}
		record[col] = newUUID
	}
	if rowErr != nil {
		if u.headerError {
			return rowErr
		}
		u.headerError = true
	}
	if err := u.output.Write(record); err != nil {
		return err
	}
	return nil
}

func (u *uuidCrypt) processUUID(preUUID string) (string, error) {
	preProc, err := uuidToBytes(preUUID)
	if err != nil {
		return "", err
	}
	postProc := u.processor.Process(preProc)
	postUUID, err := uuidFromBytes(postProc)
	if err != nil {
		return "", err
	}
	return postUUID, nil
}

func uuidToBytes(stringUUID string) ([]byte, error) {
	u, err := uuid.Parse(stringUUID)
	if err != nil {
		return nil, err
	}
	b, err := u.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func uuidFromBytes(bytes []byte) (string, error) {
	u := new(uuid.UUID)
	if err := u.UnmarshalBinary(bytes); err != nil {
		return "", err
	}
	return u.String(), nil
}

func errIfNotEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}
