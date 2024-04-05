# Agent

This is a Go project that interacts with an API at `https://phalerum.stickybits.red/api/v1/agents/test`.

## Requirements

- Go 1.21 or higher

## Usage

To run the project, use the following command:

```bash
go run main.go
```

## Building

To see the available options to build for run:
```bash
go tool dist list
```

To build on linux for a certain platform run:

```bash
GOOS=linux GOARCH=amd64 go build main.go
```

To build on windows for a certain platform run: (in elevated powershell)
```bash
$Env:GOOS = "linux"; $Env:GOARCH = "amd64"
```
Check if GOOS and GOARCH is set correctly

```bash
go env
```
Run go build

```bash
go build main.go
```

This will start the application and it will begin interacting with the API.
