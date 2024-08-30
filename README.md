# llmapiproxy
Easy llm api call from local go through vps proxy

# HTTP Reverse Proxy Server

## Overview

This is a simple HTTP reverse proxy server written in Go. It can forward received HTTP requests to a specified target server and add an authorization header to the request.

## Features

*   **Environment Variable Support**: The server can load environment variables from a `.env` file using the `github.com/joho/godotenv` package.
*   **Daemon Mode**: The server can run in daemon mode, allowing it to run in the background and not be terminated when the terminal is closed.
*   **HTTP Reverse Proxy**: The server can forward received HTTP requests to a specified target server and add an authorization header to the request.
*   **Response Header Logging**: The server can log the response headers of the target server.

## Usage

To use the server, simply compile the code and run the resulting executable. The server will start listening on the port specified in the `PORT` environment variable.

You can also run the server in daemon mode by passing the `-daemon` flag when running the executable.

## Configuration

The server can be configured using environment variables. The following environment variables are supported:

*   `PORT`: The port number to listen on.
*   `TARGET`: The URL of the target server.
*   `API_KEY`: The API key to use for authorization.

## Requirements

*   Go 1.14 or later
*   `github.com/joho/godotenv` package
*   `net/http` package
*   `net/http/httputil` package
*   `os` package
*   `os/exec` package
*   `strconv` package
*   `strings` package
*   `fmt` package

## Installation

To install the server, simply clone the repository and compile the code using the `go build` command.

## License

This software is licensed under the Apache-2.0 license
