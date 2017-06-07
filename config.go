package main

import (
	"flag"
	"os"
)

type RunConfig interface {
	Load() error
	Config() Config
	Done() error
}

type Config struct {
	inputFile   string
	outputFile  string
	secret      string
	namespace   string
	delimiter   string
	columns     string
	inPlace     bool
	decrypt     bool
	showVersion bool
}

type flagConfig struct {
	config Config
}

func NewFlagConfig() RunConfig {
	return &flagConfig{}
}

func (c *flagConfig) Load() error {
	cfg := defaultFlagsFromEnv()
	stringVarIfNoDefault(&cfg.secret, "s", "Secret key used to generate all encryption keys")
	stringVarIfNoDefault(&cfg.namespace, "n", "Namespace to generate an entity-specific encryption key")
	flag.StringVar(&cfg.delimiter, "F", "", "Custom delimiter for CSV file (default: ',')")
	// flag.StringVar(&f.columns, "c", "", "Comma-separated list of columns to encrypt/decrypt (default: 1)")
	flag.StringVar(&cfg.outputFile, "o", "-", "Output file")
	flag.BoolVar(&cfg.decrypt, "d", false, "Set operation to DECRYPT (default: ENCRYPT)")
	flag.BoolVar(&cfg.inPlace, "i", false, "Operate on the file in-place")
	flag.BoolVar(&cfg.showVersion, "version", false, "Display version information")
	flag.Parse()
	cfg.inputFile = flag.Arg(0)
	if cfg.inputFile == "" {
		cfg.inputFile = stdPipe
	}
	if err := setFilesIfInPlace(&cfg); err != nil {
		return err
	}
	c.config = cfg
	return nil
}

func (c *flagConfig) Config() Config {
	return c.config
}

func (c *flagConfig) Done() error {
	return removeBackupFileIfInPlace(c.config)
}

func defaultFlagsFromEnv() Config {
	var c Config
	c.secret = os.Getenv("UUIDCRYPT_SECRET")
	c.namespace = os.Getenv("UUIDCRYPT_NAMESPACE")
	return c
}

func stringVarIfNoDefault(s *string, name, description string) {
	currentValue := *s
	defaultValue := ""
	flag.StringVar(s, name, defaultValue, description)
	if *s == defaultValue {
		*s = currentValue
	}
}

func toBytes(str string) []byte {
	return []byte(str)
}

func toCryptType(shouldDecrypt bool) CryptType {
	if shouldDecrypt {
		return DecryptType
	}
	return EncryptType
}

func setFilesIfInPlace(c *Config) error {
	if !c.inPlace {
		return nil
	}
	backupFile, err := createBackupFile(c.inputFile)
	if err != nil {
		return err
	}
	c.outputFile = c.inputFile
	c.inputFile = backupFile
	return nil
}

func createBackupFile(filename string) (string, error) {
	if filename == stdPipe {
		return stdPipe, nil
	}
	newFilename := filename + ".tmp.uuidcrypt"
	if err := os.Rename(filename, newFilename); err != nil {
		return "", err
	}
	return newFilename, nil
}

func removeBackupFileIfInPlace(c Config) error {
	if !c.inPlace || c.inputFile == stdPipe {
		return nil
	}
	return os.Remove(c.inputFile)
}
