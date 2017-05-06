package main

import (
	"io"

	"github.com/google/uuid"
)

type UUIDCrypt interface {
	Run() error
}

func NewUUIDCrypt(input File, output File, processor Processor, columns ...int) UUIDCrypt {
	return &uuidCrypt{
		input:     input,
		output:    output,
		processor: processor,
		columns:   columns,
	}
}

type uuidCrypt struct {
	input     File
	output    File
	processor Processor
	columns   []int
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
	for _, column := range u.columns {
		col := column - 1
		newUUID, err := u.processUUID(record[col])
		if err != nil {
			return err
		}
		record[col] = newUUID
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
