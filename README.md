# Bucket Janitor
##### This is work in progress

Use on S3 buckets to verify the content.  
More specifically, to make sure files have certain file extensions
and contain data of the appropriate type.  

### Use case:
When you provide 3rd party access to a bucket, there are
certain risks related to the content that is uploaded.  
This tool will use the magic library to verify the type
of the files and their extensions.

### Note:
This is meant to be an excercise in concurrent coding.  
Semaphores are being used to limit concurrency, objects are pipelined through a process
of listing, downloading, detecting type and aggregating. For now a table is being
printed of the detected mime/magic combination and their counts.

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
    |-- aws
    |   `-- main.go
    `-- janitor
        |-- count.go
        |-- detect.go
        |-- list.go
        |-- main.go
        |-- prepare.go
        `-- retrieve.go
```

### Use:
Copy and edit the example json from the configs folder, then execute as below:
```
go build cmd/bucket-janitor.go
./bucket-janitor.go -c config.json
```
