# uuidcrypt
A tool for two-way encryption of UUIDs in CSV files as a layer of obfuscation. It is not entirely secure, but helps prevent a class of UUID-enumeration attack vectors.

## Install

```
go get github.com/APTy/uuidcrypt
```

## Usage

### Encrypt UUIDs
```
uuidcrypt -s 'my secret password' -n 'namespace-foo' myfile.csv
```

### Decrypt UUIDs
```
uuidcrypt -d -s 'my secret password' -n 'namespace-foo' myfile.csv
```

By default, it will parse the CSV as comma-delimited (',') and encrypt/decrypt all UUIDs in the first column only.


### Testing
```
# encrypt a random uuid

$ echo -n 558ece65-c7c8-4ad2-83dd-f696b2c540a4 | uuidcrypt -s -secret -n namespace
4d4054d3-8360-97b4-bb45-44ab833b79f0

# decrypt it back

$ echo -n 4d4054d3-8360-97b4-bb45-44ab833b79f0 | uuidcrypt -d -s -secret -n namespace
558ece65-c7c8-4ad2-83dd-f696b2c540a4

# the result should be the same as the original UUID
```

