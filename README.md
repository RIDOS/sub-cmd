[![Go](https://github.com/RIDOS/sub-cmd/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/RIDOS/sub-cmd/actions/workflows/go.yml)

# SUB-CMD

## Get Started

This application is built with Go version 1.23.2 and demonstrates the usage of subcommands.

To build the application and view help documentation, run the following commands:

```bash
go build -o sub-cmd
./sub-cmd http -h
```

## Usage

The basic usage of the application is as follows:

```bash
sub-cmd [http|grpc] -h
```

### Subcommands

- **http**: A simple HTTP client.
- **grpc**: A gRPC client.

#### HTTP Subcommand

To get help on the HTTP subcommand, use:

```bash
./sub-cmd http --help
```

**Example Usage**:

```bash
./sub-cmd http -verb POST http://localhost
```

This will execute the HTTP command with the specified method and URL. For the HTTP subcommand, the available methods are `GET`, `POST`, and `HEAD`.

**HTTP Options**:
- `-verb string`: Specifies the HTTP method (default is "GET").

#### gRPC Subcommand

To get help on the gRPC subcommand, use:

```bash
./sub-cmd grpc --help
```

**Example Usage**:

```bash
./sub-cmd grpc -method YourMethodName -body "Your request body"
```

**gRPC Options**:
- `-method string`: Specifies the method to call.
- `-body string`: Specifies the body of the request.

## Links

https://github.com/practicalgo/code