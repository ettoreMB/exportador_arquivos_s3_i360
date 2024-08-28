# AWS S3 File Export System

This system downloads data from all tables described in `config.yaml` and exports them to an Amazon S3 bucket. After exporting, the files are saved in a ZIP archive.

## Requirements

- Go 1.23.0

## Build Instructions

To build the project, run the following commands:

```sh
go build -o executar
go build -o executar.exe
```


## Configuration
Ensure that config.yaml is properly configured with the necessary table descriptions and AWS S3 bucket details.