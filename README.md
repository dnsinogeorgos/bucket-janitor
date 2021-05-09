# Bucket Janitor
##### This is work in progress - experimental

Use on S3 buckets to verify the state of objects.  
More specifically, to make sure files have certain file extensions
and contain data of the appropriate type.

```
|-- LICENSE
|-- README.md
|-- bucket-janitor
|-- cmd
|   `-- bucket-janitor.go
|-- config.json
|-- configs
|   `-- example.json
|-- git@github.com:dnsinogeorgos
|   `-- bucket-janitor.git
|-- go.mod
|-- go.sum
`-- internal
    |-- config
    |   `-- config.go
    |-- load
    |   |-- bucket.go
    |   |-- downloader.go
    |   |-- object.go
    |   `-- session.go
    `-- verify
        `-- header.go
```

### Use:
You can't really use it yet, but if you're inclined...  
copy the example json from the configs folder, then execute as below:
```
go build cmd/bucket-janitor.go
./bucket-janitor.go -c config.json
```

### Progress:
- [ ] complete imperative pieces
- [ ] restructure for concurrency
- [ ] retrieve values from flags/env but also json file

### What it does:
##### receives:
* aws credentials
* a list of buckets
* a struct of file extensions (key) and an array of magic types (value)
* concurrency/rate limiting options
* (optional) a json file with corrective actions to take on bucket objects

##### returns:
* the object paths that do not comply with the allowed file/magic types
* the corrective actions required to take on the objects

### Use case:  
When you provide 3rd party access to a bucket, there are
certain risks related to the content that is uploaded.  
This tool will use the magic library to verify the type
of the file and the extension.
