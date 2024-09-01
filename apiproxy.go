// Package apiproxy 
// This package implements a simple HTTP reverse proxy server, 
// supporting background running and environment variable configuration. 
// Author: [Michael lai] 
// Email: [gm2068@gmail]
// Version: 1.0.0 
// Created: [2024-8-25]
// Last Updated: [2024-9-1] 
// 
// Features: 
// - Supports background running 
// - Supports environment variable configuration 
// - Supports HTTP reverse proxy 
// - Supports log recording
// Usage: 
// 1. Set environment variables PORT, TARGET, and APIKEY 
// 2. Run the program 
//   `go run apiproxy. -daemon`
// 3.check daemon running
//   `ps -ef | grep go`
// 4.check api.log file
// 5.stop running
//   `kill PID`

package main

import (
    "flag"
    "log"
    "net/http"
    "net/http/httputil"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "fmt"
    "github.com/joho/godotenv"
)

const (
    daemonFlag = "daemon"
)

var (
    runDaemon bool
)

func init() {
    if err := godotenv.Load(); err!= nil {
        log.Fatal("Error loading .env file")
    }
}

func getEnvVar(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("%s environment variable is not set", key)
    }
    return value
}

func getPort() int {
    portStr := getEnvVar("PORT")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        log.Fatal("Invalid port number")
    }
    return port
}

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("[1] receive a request from %s, request header: %s: \n", r.RemoteAddr, r.Header)
    target := getEnvVar("target")
    apiKey := getEnvVar("apiKey")
    director := func(req *http.Request) {
        req.URL.Scheme = "https"
        req.URL.Host = target
        req.Host = target
        authHeader := fmt.Sprintf("Bearer %s", apiKey)
        req.Header.Set("Authorization", authHeader)
    }
    proxy := &httputil.ReverseProxy{Director: director}
    proxy.ServeHTTP(w, r)
    logPrintResponseHeaders(w)
}

func logPrintResponseHeaders(w http.ResponseWriter) {
    headers := w.Header().Clone()
    for key, values := range headers {
        log.Printf("[*] %s: %s\n", key, strings.Join(values, ", "))
    }
}

func stripSlice(slice []string, element string) []string {
    for i := len(slice) - 1; i >= 0; i-- {
        if slice[i] == element {
            slice = append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}

func subProcess(args []string) *exec.Cmd {
    cmd := exec.Command(args[0])
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Start(); err!= nil {
        log.Printf("[-] Error: %s\n", err)
    }
    return cmd
}

func main() {
    flag.BoolVar(&runDaemon, daemonFlag, false, "Run as daemon")
    flag.Parse()

    // Check if the file exists
	if _, err := os.Stat("api.log"); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		logFile, err := os.Create("api.log")
		if err != nil {
			log.Fatal(err)
		}
		logFile.Close()
	}
    logFile, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    log.Printf("[1] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)

    if runDaemon {
        args := stripSlice(os.Args, "-"+daemonFlag)
        cmd := subProcess(args)
        defer cmd.Process.Kill()
        log.Printf("[2] Daemon running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
        os.Exit(0)
    }

    log.Printf("[3] Forever running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())

    port := getPort()
    log.Printf("[4] Starting server at port %v\n", port)
    log.Printf("[5] start writing api.log")

    http.HandleFunc("/", ReverseProxyHandler)
    if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err!= nil {
        log.Fatal(err)
    }
}
