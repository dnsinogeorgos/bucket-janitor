# Bucket Janitor
##### This is work in progress - experimental

Use on S3 buckets to verify the state of objects.  
More specifically, to make sure files have certain file extensions
and contain data of the appropriate type.  

### Note:
This is meant to be an excercise in concurrent coding.  
I'm using the inverted worker pool pattern to list bucket object.
I will be using a custom HTTP client to enable persistent connections to S3.  
I intend to implement a resource pool of persistent connections to S3 in order to speed
up fetching object headers.  

```
|-- LICENSE
|-- README.md
|-- cmd
|   `-- bucket-janitor.go
|-- configs
|   `-- example.json
|-- go.mod
|-- go.sum
`-- internal
    |-- config
    |   `-- config.go
    |-- load
    |   |-- bucket.go
    |   |-- main.go
    |   `-- object.go
    `-- verify
        `-- header.go
```

### Use:
It's still not in a useful state, but if you're inclined...  
copy and edit the example json from the configs folder, then execute as below:
```
go build cmd/bucket-janitor.go
./bucket-janitor.go -c config.json
```

### Progress:
- [ ] complete imperative pieces
- [ ] restructure for concurrency
- [ ] retrieve values from flags/env but also json file
- [ ] export/apply diff functionality

### What it will do:
##### receive:
* aws credentials
* a list of buckets
* (soon) a set of file extensions and magic types
* (soon) concurrency/rate limiting options
* (optional) a json file with corrective actions to take on bucket objects (exported by
  the same tool a-la-terraform plan)

##### return:
* the object keys that do not comply with the allowed file/magic types
* the corrective actions required to take on the objects

### Use case:  
When you provide 3rd party access to a bucket, there are
certain risks related to the content that is uploaded.  
This tool will use the magic library to verify the type
of the file and the extension.
