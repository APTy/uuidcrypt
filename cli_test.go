package main

import (
	"encoding/csv"
	"os"
	"testing"
)

const (
	testSecret    = "foo"
	testNamespace = "bar"

	testDir             = "testdata/"
	testInputFile       = testDir + "testfile.csv"
	testEncInputFile    = testDir + "testfile.csv.enc"
	testOutputFile      = testDir + ".testfile.csv"
	testOutputFile2     = testDir + ".testfile.csv2"
	testOutputFile3     = testDir + ".testfile.csv3"
	testBinaryDelimFile = testDir + "binary-delim.csv"
)

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("encountered error: %v\n", err)
	}
}

func assert(t *testing.T, condition bool, description string) {
	if !condition {
		t.Fatalf(description)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	defer os.Remove(testOutputFile)
	defer os.Remove(testOutputFile2)
	defer os.Remove(testOutputFile3)

	// encrypt
	runCLIWithMockConfig(Config{
		inputFile:  testInputFile,
		outputFile: testOutputFile,
		secret:     testSecret,
		namespace:  testNamespace,
	})

	// test encrypt success with regression comparison
	input := getRecordsFromCSV(t, testInputFile)
	encInput := getRecordsFromCSV(t, testEncInputFile)
	output := getRecordsFromCSV(t, testOutputFile)
	assert(t, len(input) == len(encInput), "num input rows should match num encrypted input rows")
	assert(t, len(input) == len(output), "num input rows should match num output rows")
	for i := range input {
		in := input[i]
		out := output[i]
		encIn := encInput[i]
		assert(t, in[0] != out[0], "input uuid should not match output uuid")
		assert(t, in[1] == out[1], "other input data should match output data")
		assert(t, in[2] == out[2], "other input data should match output data")
		assert(t, out[0] == encIn[0], "output uuid should match encrypted input uuid")
		assert(t, out[1] == encIn[1], "other output data should match encrypted input data")
		assert(t, out[2] == encIn[2], "other output data should match encrypted input data")
	}

	// decrypt with correct key
	runCLIWithMockConfig(Config{
		inputFile:  testOutputFile,
		outputFile: testOutputFile2,
		secret:     testSecret,
		namespace:  testNamespace,
		decrypt:    true,
	})

	// test decrypt success
	output = getRecordsFromCSV(t, testOutputFile2)
	assert(t, len(input) == len(output), "num input rows should match num output rows")
	for i := range input {
		in := input[i]
		out := output[i]
		assert(t, in[0] == out[0], "input uuid should match output uuid")
		assert(t, in[1] == out[1], "other input data should match output data")
		assert(t, in[2] == out[2], "other input data should match output data")
	}

	// decrypt with bad key
	runCLIWithMockConfig(Config{
		inputFile:  testOutputFile,
		outputFile: testOutputFile3,
		secret:     testSecret,
		namespace:  "wrong",
		decrypt:    true,
	})

	// test decrypt failure
	output = getRecordsFromCSV(t, testOutputFile3)
	assert(t, len(input) == len(output), "num input rows should match num output rows")
	for i := range input {
		in := input[i]
		out := output[i]
		assert(t, in[0] != out[0], "input uuid should not match output uuid")
		assert(t, in[1] == out[1], "other input data should match output data")
		assert(t, in[2] == out[2], "other input data should match output data")
	}
}

func runCLIWithMockConfig(config Config) {
	cli := NewCLI(newMockConfig(config))
	cli.Run()
}

func getRecordsFromCSV(t *testing.T, filename string) [][]string {
	return getRecordsFromCSVWithDelim(t, filename, ',')
}

func getRecordsFromCSVWithDelim(t *testing.T, filename string, delim rune) [][]string {
	f, err := os.Open(filename)
	failIfError(t, err)
	r := csv.NewReader(f)
	r.Comma = delim
	records, err := r.ReadAll()
	failIfError(t, err)
	return records
}

func TestBinaryDelim(t *testing.T) {
	testOutputFile := testBinaryDelimFile + ".1"
	expectedOutputFile := testBinaryDelimFile + ".enc"
	defer os.Remove(testOutputFile)

	// encrypt
	runCLIWithMockConfig(Config{
		inputFile:  testBinaryDelimFile,
		outputFile: testOutputFile,
		secret:     testSecret,
		namespace:  testNamespace,

		// simulate unescaped strings entered at the command line
		delimiter:       "\\x01",
		delimiterOutput: "\\x02",
	})

	// test encrypt success with regression comparison
	input := getRecordsFromCSVWithDelim(t, testBinaryDelimFile, '\x01')
	expectedOutput := getRecordsFromCSVWithDelim(t, expectedOutputFile, '\x02')
	testOutput := getRecordsFromCSVWithDelim(t, testOutputFile, '\x02')
	assert(t, len(input) == len(expectedOutput), "num input rows should match num encrypted input rows")
	assert(t, len(input) == len(testOutput), "num input rows should match num output rows")
	for i := range input {
		in := input[i]
		out := testOutput[i]
		encIn := expectedOutput[i]
		assert(t, in[0] != out[0], "input uuid should not match output uuid")
		assert(t, in[1] == out[1], "other input data should match output data")
		assert(t, in[2] == out[2], "other input data should match output data")
		assert(t, out[0] == encIn[0], "output uuid should match encrypted input uuid")
		assert(t, out[1] == encIn[1], "other output data should match encrypted input data")
		assert(t, out[2] == encIn[2], "other output data should match encrypted input data")
	}
}
