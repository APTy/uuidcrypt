package main

import (
	"fmt"
	"os"
)

type CLI struct {
	cfg RunConfig
}

func NewCLI(cfg RunConfig) CLI {
	return CLI{cfg: cfg}
}

func (c CLI) Run() {
	if err := c.run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func (c CLI) run() error {
	if err := c.cfg.Load(); err != nil {
		return err
	}
	cfg := c.cfg.Config()
	if cfg.showVersion {
		fmt.Fprintf(os.Stdout, "uuidcrypt %s\n", Version)
		return nil
	}
	uuidCrypt := NewUUIDCrypt(
		NewCSVFile(cfg.inputFile, WithDelimiter(cfg.delimiter)),
		NewCSVFile(cfg.outputFile, WithDelimiter(cfg.delimiter)),
		NewCrypterProcessor(toBytes(cfg.secret), toBytes(cfg.namespace), toCryptType(cfg.decrypt)),
		WithColumns(cfg.columns...),
	)
	if err := uuidCrypt.Run(); err != nil {
		return err
	}
	if err := c.cfg.Done(); err != nil {
		return err
	}
	return nil
}
