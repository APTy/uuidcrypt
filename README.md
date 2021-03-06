# uuidcrypt
Implements two-way [format-preserving encryption](https://en.wikipedia.org/wiki/Format-preserving_encryption) for UUIDs in CSV files.
It may be used when sending data to third parties to help prevent a class of UUID-enumeration attacks in case of data leakage.

## Install

```
go get github.com/APTy/uuidcrypt
```

## Simple example with a CSV file

### Sample input (plaintext)
A sample CSV file with UUIDs in the first column.
``` bash
$ cat testdata/testfile.csv
4a1981ca-94af-481d-8266-58d86cc8199a,other,data
37abbed5-e81e-45d6-a6d4-3548685203cc,other,data
c3263b22-ed7b-45b8-9e3f-f55b16b3a37f,other,data
2cdca796-68cf-45b7-90ce-fad929209d3d,other,data
```

### Encrypt UUIDs
Encrypt the UUIDs in the CSV and make a temp copy.
``` bash
$ uuidcrypt -s 'my secret password' -n 'namespace-foo' testdata/testfile.csv | tee /tmp/testfile.csv.enc
aeee3314-701a-0e2e-ae26-050735153353,other,data
69d0f027-bec9-1bdc-966b-930f0766367c,other,data
ef98ae92-eed6-f41e-bdc2-b9ce61ce6b59,other,data
9f0ac2ce-44d3-bb5e-d41b-3aa5d18d0242,other,data
```

### Decrypt UUIDs
Decrypt the UUIDs in the CSV from the temp copy and verify that its the same as the plaintext input.
``` bash
$ uuidcrypt -d -s 'my secret password' -n 'namespace-foo' /tmp/testfile.csv.enc
4a1981ca-94af-481d-8266-58d86cc8199a,other,data
37abbed5-e81e-45d6-a6d4-3548685203cc,other,data
c3263b22-ed7b-45b8-9e3f-f55b16b3a37f,other,data
2cdca796-68cf-45b7-90ce-fad929209d3d,other,data
```

By default, it will parse the CSV as comma-delimited (`','`) and encrypt/decrypt all UUIDs in the first column only.
See Usage for configuring the field delimiter and which columns are transformed.

## Usage
``` bash
$ uuidcrypt -help
Usage of uuidcrypt:
  -F string
        Field separator for CSV file (default: ',')
  -OF string
        Field separator for output CSV file (default: ',')
  -c string
        Comma-separated list of columns to encrypt/decrypt (default: 1)
  -d    Set operation to DECRYPT (default: ENCRYPT)
  -i    Operate on the file in-place
  -n string
        Namespace to generate an entity-specific encryption key
  -o string
        Output file (default "-")
  -s string
        Secret key used to generate all encryption keys
  -version
        Display version information
```

### Environment Variables
You can set `secret` and `namespace` configuration using environment variables.

``` bash
$ export UUIDCRYPT_SECRET="my secret"
$ export UUIDCRYPT_NAMESPACE="my namespace"
$ echo 558ece65-c7c8-4ad2-83dd-f696b2c540a4 | uuidcrypt
0fe55a05-c4d4-6858-2ef7-9176b3d1312f

$ echo 0fe55a05-c4d4-6858-2ef7-9176b3d1312f | uuidcrypt -d
558ece65-c7c8-4ad2-83dd-f696b2c540a4
```

### Custom CSV field separator/delimiter

Delimit input by a tab (`\t`) and delimit output by a space (` `).
``` bash
$ echo -e 'd13d625c-f451-40b8-91e6-7b56589b91f1\t123\t456' | uuidcrypt -F '\t' -OF ' '
66281a1f-eb55-59fd-7676-c9e50560ca42 123 456
```

### Custom columns

Operate on columns `2` and `3`.
``` bash
$ echo -e 'foo,d13d625c-f451-40b8-91e6-7b56589b91f1,d13d625c-f451-40b8-91e6-7b56589b91f1,123,456' | uuidcrypt -c 2,3
foo,66281a1f-eb55-59fd-7676-c9e50560ca42,66281a1f-eb55-59fd-7676-c9e50560ca42,123,456
```
