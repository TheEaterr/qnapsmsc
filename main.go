package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/TheEaterr/qnapsmsc/lib/notifications"
	"github.com/TheEaterr/qnapsmsc/lib/utils"
)

const (
	notificationEndpoint = "/notification"
)

type httpServerArgs struct {
	port     string
	logger   *log.Logger
	username string
	password string
}

func main() {
	runtime.GOMAXPROCS(0)

	port := flag.String("port", ":9094", "Port to serve at (e.g. :9094).")
	username := flag.String("username", "admin", "Username to connect to the notification server.")
	password := flag.String("password", "notsecure", "Password to connect to the notification server.")
	logFile := flag.String("log", "", "Log file path (defaults to empty, i.e. STDOUT).")
	handler := flag.String("handler", "log", "Handler to use for notifications (log or mail).")
	mailSender := flag.String("mail-sender", "", "Email address to use as sender.")
	mailReceiver := flag.String("mail-receiver", "", "Email address to use as receiver.")
	smtpHost := flag.String("smtp-host", "localhost", "SMTP host to use for sending emails.")
	smtpPort := flag.Int("smtp-port", 587, "SMTP port to use for sending emails.")
	smtpUsername := flag.String("smtp-username", "", "SMTP username to use for sending emails.")
	smtpPassword := flag.String("smtp-password", "", "SMTP password to use for sending emails.")
	defaultUsage := flag.Usage
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "qnapsmsc version %s (%s-%s) built on %s\n", utils.VERSION, utils.REVISION, utils.BRANCH, utils.BUILT)
		fmt.Fprintln(flag.CommandLine.Output(), "")
		defaultUsage()
	}
	flag.Parse()

	log.Println("Using options:")
	log.Printf("  - Port: %s\n", *port)
	log.Printf("  - Username: %s\n", *username)
	log.Printf("  - Log file: %s\n", *logFile)
	log.Printf("  - Handler: %s\n", *handler)
	log.Printf("  - Mail sender: %s\n", *mailSender)
	log.Printf("  - Mail receiver: %s\n", *mailReceiver)
	log.Printf("  - SMTP host: %s\n", *smtpHost)
	log.Printf("  - SMTP port: %d\n", *smtpPort)
	log.Printf("  - SMTP username: %s\n", *smtpUsername)

	if *handler != "log" && *handler != "mail" {
		log.Fatalf("Invalid handler: %s\n", *handler)
	}
	if *handler == "mail" && (*mailSender == "" || *mailReceiver == "") {
		log.Fatalf("Mail handler requires both sender and receiver to be set\n")
	}

	var logWriter io.Writer = os.Stderr
	if *logFile != "" {
		lf, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Error creating log file: %v\n", err)
		}
		defer lf.Close()

		logWriter = lf
	}
	logger := log.New(logWriter, "", log.LstdFlags)

	args := httpServerArgs{
		port:     *port,
		logger:   logger,
		username: *username,
		password: *password,
	}
	notifHandler := notifications.NewLogHandler(logger)
	if *handler == "mail" {
		notifHandler = notifications.NewMailHandler(
			logger,
			*mailSender,
			*mailReceiver,
			*smtpHost,
			*smtpPort,
			*smtpUsername,
			*smtpPassword,
		)
	}

	ctx, cancelFn := context.WithCancel(context.Background())

	// Setup our Ctrl+C handler
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		defer cancelFn()

		// Wait for program exit
		<-exitCh
	}()

	err := serveHTTP(ctx, args, notifHandler)
	if err != nil {
		log.Println(err.Error())
	}
	os.Exit(1)
}

func handleNotificationHTTPRequest(args httpServerArgs, w http.ResponseWriter, r *http.Request, annotator notifications.Handler) {
	log.Printf("Received notification request from %s\n", r.RemoteAddr)
	username := r.URL.Query().Get("username")
	if len(username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Returning 400 Bad Request, missing username")
		return
	}
	password := r.URL.Query().Get("password")
	if len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Returning 400 Bad Request, missing password")
		return
	}
	if username != args.username || password != args.password {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Returning 401 Unauthorized")
		return
	}

	notification := r.URL.Query().Get("text")
	if len(notification) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Returning 400 Bad Request, missing text")
		return
	}

	_, _ = annotator.Post(notification)
}

func serveHTTP(ctx context.Context, args httpServerArgs, handler notifications.Handler) error {
	http.HandleFunc(notificationEndpoint, func(w http.ResponseWriter, r *http.Request) {
		handleNotificationHTTPRequest(args, w, r, handler)
	})

	// listen to port
	server := http.Server{Addr: args.port}
	server.ErrorLog = args.logger
	go func() {
		log.Printf("Listening to HTTP requests at %s\n", args.port)

		// Wait for program exit
		<-ctx.Done()

		log.Println("Program aborted, exiting...")
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println(err.Error())
		}
	}()

	return server.ListenAndServe()
}
