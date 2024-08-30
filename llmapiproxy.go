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
	"syscall"
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
	flag.BoolVar(&runDaemon, daemonFlag, false, "set Background mode")
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

func main() {
	flag.Parse()
	if runDaemon {
		// Fork process and continue as daemon
		pid := os.Getpid()
		log.Printf("[*] Daemonizing, original PID: %d\n", pid)
		if pid > 0 {
			// Create a new process
			cmd := exec.Command(os.Args[0])
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
			_ = cmd.Start()
			log.Printf("[*] Daemon started, PID: %d\n", cmd.Process.Pid)
			// Exit parent process
			defer cmd.Process.Kill() // release resources when the program exits
			log.Printf("[*] Closing parent process, PID: %d\n", pid)
			os.Exit(0)
		}
	}
	log.Printf("[*] PID: %d PPID: %d ARG: %s\n", os.Getpid(), os.Getppid(), os.Args)
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT environment variable is not set")
	}
	var err error
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid port number")
	}
	log.Printf("[*] Starting server at port %v\n", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ReverseProxyHandler(w, r)
	})
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err!= nil {
		log.Fatal(err)
	}
}