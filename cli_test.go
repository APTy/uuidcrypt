package main

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"testing"
)

const (
	testInputFile   = "testfile.csv"
	testOutputFile  = ".testfile.csv"
	testOutputFile2 = ".testfile.csv2"
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
	// setup
	_, err := ioutil.TempFile("", testOutputFile)
	failIfError(t, err)
	defer os.Remove(testOutputFile)
	_, err = ioutil.TempFile("", testOutputFile2)
	failIfError(t, err)
	defer os.Remove(testOutputFile2)

	// encrypt
	runCLIWithMockConfig(Config{
		inputFile:  testInputFile,
		outputFile: testOutputFile,
	})

	// test encrypt success
	input := getRecordsFromCSV(t, testInputFile)
	output := getRecordsFromCSV(t, testOutputFile)
	assert(t, len(input) == len(output), "num input rows should match num output rows")
	for i := range input {
		in := input[i]
		out := output[i]
		assert(t, in[0] != out[0], "input uuid should not match output uuid")
		assert(t, in[1] == out[1], "other input data should match output data")
		assert(t, in[2] == out[2], "other input data should match output data")
	}

	// decrypt
	runCLIWithMockConfig(Config{
		inputFile:  testOutputFile,
		outputFile: testOutputFile2,
		decrypt:    true,
	})

	// test decrypt success
	input = getRecordsFromCSV(t, testInputFile)
	output = getRecordsFromCSV(t, testOutputFile2)
	assert(t, len(input) == len(output), "num input rows should match num output rows")
	for i := range input {
		in := input[i]
		out := output[i]
		assert(t, in[0] == out[0], "input uuid should match output uuid")
		assert(t, in[1] == out[1], "other input data should match output data")
		assert(t, in[2] == out[2], "other input data should match output data")
	}
}

func runCLIWithMockConfig(config Config) {
	cli := NewCLI(newMockConfig(config))
	cli.Run()
}

func getRecordsFromCSV(t *testing.T, filename string) [][]string {
	f, err := os.Open(filename)
	failIfError(t, err)
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	failIfError(t, err)
	return records
}
