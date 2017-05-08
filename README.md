# uuidcrypt
A tool for encrypting files with UUIDs.

## Install

```
go get github.com/APTy/uuidcrypt
```

## Usage

```
uuidcrypt -s 'my secret password' myfile.csv
```

By default, it will parse the CSV as comma-delimited (',') and encrypt all UUIDs in the first column only.
