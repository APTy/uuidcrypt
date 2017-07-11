package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

// RunConfig decides how uuidcrypt will run.
type RunConfig interface {
	Load() error
	Config() Config
	Done() error
}

// Config provides user input to the CLI invocation.
// TODO: a lot of config could live with UUIDCrypt as opposed to
//       being so tightly coupled to the CLI.
type Config struct {
	inputFile       string
	outputFile      string
	secret          string
	namespace       string
	delimiter       string
	delimiterOutput string
	columns         []int
	inPlace         bool
	decrypt         bool
	showVersion     bool
}

type flagConfig struct {
	config Config
}

func NewFlagConfig() RunConfig {
	return &flagConfig{}
}

func (c *flagConfig) Load() error {
	cfg := defaultFlagsFromEnv()
	var columns string
	stringVarIfNoDefault(&cfg.secret, "s", "Secret key used to generate all encryption keys")
	stringVarIfNoDefault(&cfg.namespace, "n", "Namespace to generate an entity-specific encryption key")
	flag.StringVar(&cfg.delimiter, "F", "", "Custom delimiter for CSV file (default: ',')")
	flag.StringVar(&cfg.delimiterOutput, "OF", "", "Custom output delimiter for CSV file (default: ',')")
	flag.StringVar(&columns, "c", "", "Comma-separated list of columns to encrypt/decrypt (default: 1)")
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
	if intColumns, err := parseColumns(columns); err != nil {
		return err
	} else {
		cfg.columns = intColumns
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

func parseColumns(columns string) ([]int, error) {
	var intColumns []int
	strippedColumns := strings.Replace(columns, " ", "", -1)
	splitColumns := strings.Split(strippedColumns, ",")
	for i, col := range splitColumns {
		if len(splitColumns) == 1 && col == "" {
			break
		}
		if col == "" {
			continue
		}
		if intColumns == nil {
			intColumns = make([]int, len(splitColumns))
		}
		intCol, err := strconv.Atoi(col)
		if err != nil {
			return nil, err
		}
		intColumns[i] = intCol
	}
	return intColumns, nil
}
