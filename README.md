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
1. Set environment variables PORT, TARGET, and APIKEY

2. Run the program 

3.check daemon running

4.check api.log file

5.stop running（kill PID)

`go run apiproxy. -daemon`

`ps -ef | grep go`

`kill 732924`

## Configuration

The server can be configured using environment variables. The following environment variables are supported:

*   `PORT`: The port number to listen on.
*   `TARGET`: The URL of the target llm server.
*   `API_KEY`: The API key to use for llm authorization.

## License

This software is licensed under the Apache-2.0 license
