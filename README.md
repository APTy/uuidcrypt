# uuidcrypt
A tool for two-way encryption of UUIDs in CSV files as a layer of obfuscation. It is not entirely secure, but helps prevent a class of UUID-enumeration attack vectors.

## Install

```
go get github.com/APTy/uuidcrypt
```

## Examples

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

$ echo 558ece65-c7c8-4ad2-83dd-f696b2c540a4 | uuidcrypt -s 'my secret' -n 'my namespace'
0fe55a05-c4d4-6858-2ef7-9176b3d1312f

# decrypt it back

$ echo 0fe55a05-c4d4-6858-2ef7-9176b3d1312f | uuidcrypt -s 'my secret' -n 'my namespace' -d
558ece65-c7c8-4ad2-83dd-f696b2c540a4

# the result should be the same as the original UUID
```

### Environment Variables
You can set `secret` and `namespace` configuration using environment variables.

```
$ export UUIDCRYPT_SECRET="my secret"
$ export UUIDCRYPT_NAMESPACE="my namespace"
$ echo 558ece65-c7c8-4ad2-83dd-f696b2c540a4 | uuidcrypt
0fe55a05-c4d4-6858-2ef7-9176b3d1312f

$ echo 0fe55a05-c4d4-6858-2ef7-9176b3d1312f | uuidcrypt -d
558ece65-c7c8-4ad2-83dd-f696b2c540a4
```

## Usage
```
$ uuidcrypt -help
Usage of uuidcrypt:
  -F string
        Custom delimiter for CSV file (default: ',')
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
