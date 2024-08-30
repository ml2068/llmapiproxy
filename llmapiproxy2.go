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
	port     int
	runDaemon bool
)

func init() {
	if err := godotenv.Load(); err!= nil {
		log.Fatal("Error loading .env file")
	}
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT environment variable is not set")
	}
	port, err := strconv.Atoi(portStr)
	if err!= nil {
		log.Fatal("Invalid port number")
	}
	flag.IntVar(&port, "port", port, "监听端口")
	flag.BoolVar(&runDaemon, daemonFlag, false, "是否后台运行")
}

func ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[*] receive a request from %s, request header: %s: \n", r.RemoteAddr, r.Header)
	target := os.Getenv("target")
	if target == "" {
		log.Fatal("target environment variable is not set")
	}
	apiKey := os.Getenv("apiKey")
	if apiKey == "" {
		log.Fatal("apiKey environment variable is not set")
	}
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
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err!= nil {
		log.Printf("[-] Error: %s\n", err)
	}
	return cmd
}

func main() {
	flag.Parse()
	log.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	if runDaemon {
		args := stripSlice(os.Args, "-"+daemonFlag)
		cmd := subProcess(args)
		defer cmd.Process.Kill()
		log.Printf("[*] Daemon running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
		os.Exit(0)
	}
	log.Printf("[*] Forever running in PID: %d PPID: %d\n", os.Getpid(), os.Getppid())
		if port == 0 {
		portStr := os.Getenv("PORT")
		if portStr == "" {
			log.Fatal("PORT environment variable is not set")
		}
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatal("Invalid port number")
		}
	}
	log.Printf("[*] Starting server at port %v\n", port)
	server := http.Server{Addr: ":" + strconv.Itoa(port)}
	defer server.Close()
	if err := http.ListenAndServe(":"+strconv.Itoa(port), http.HandlerFunc(ReverseProxyHandler)); err!= nil {
		log.Fatal(err)
	}
}