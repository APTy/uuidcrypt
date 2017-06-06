package main

import (
	"flag"
	"os"
)

type flags struct {
	inputFile  string
	outputFile string
	secret     string
	namespace  string
	delimiter  string
	columns    string
	inPlace    bool
	decrypt    bool
}

func RunCLI() {
	if err := runCLI(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
	}
}

func runCLI() error {
	f, err := parseFlags()
	if err != nil {
		return err
	}
	uuidCrypt := NewUUIDCrypt(
		NewCSVFile(f.inputFile),
		NewCSVFile(f.outputFile),
		NewCrypterProcessor(toBytes(f.secret), toBytes(f.namespace), toCryptType(f.decrypt)),
		WithColumns(1),
	)
	if err := uuidCrypt.Run(); err != nil {
		return err
	}
	if err := removeBackupFileIfInPlace(f); err != nil {
		return err
	}
	return nil
}

func parseFlags() (flags, error) {
	f := defaultFlagsFromEnv()
	stringVarIfNoDefault(&f.secret, "s", "Secret key used to generate all encryption keys")
	stringVarIfNoDefault(&f.namespace, "n", "Namespace to generate an entity-specific encryption key")
	// flag.StringVar(&f.delimiter, "F", "", "Custom delimiter for CSV file (default: ',')")
	// flag.StringVar(&f.columns, "c", "", "Comma-separated list of columns to encrypt/decrypt (default: 1)")
	flag.StringVar(&f.outputFile, "o", "-", "Output file")
	flag.BoolVar(&f.decrypt, "d", false, "Set operation to DECRYPT (default: ENCRYPT)")
	flag.BoolVar(&f.inPlace, "i", false, "Operate on the file in-place")
	flag.Parse()
	f.inputFile = flag.Arg(0)
	if f.inputFile == "" {
		f.inputFile = stdPipe
	}
	if err := setFilesIfInPlace(&f); err != nil {
		return f, err
	}
	return f, nil
}

func defaultFlagsFromEnv() flags {
	var f flags
	f.secret = os.Getenv("UUIDCRYPT_SECRET")
	f.namespace = os.Getenv("UUIDCRYPT_NAMESPACE")
	return f
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

func setFilesIfInPlace(f *flags) error {
	if !f.inPlace {
		return nil
	}
	backupFile, err := createBackupFile(f.inputFile)
	if err != nil {
		return err
	}
	f.outputFile = f.inputFile
	f.inputFile = backupFile
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

func removeBackupFileIfInPlace(f flags) error {
	if !f.inPlace || f.inputFile == stdPipe {
		return nil
	}
	return os.Remove(f.inputFile)
}
