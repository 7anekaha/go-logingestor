package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	internal "github.com/7anekaha/go-logingestor/api/server"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	user     = "nico"
	pass     = "secret"
	hostName = "test-mongo"
	port     = "27017"
)

var minPool string
var maxPool string

func main() {

	flag.StringVar(&minPool, "minPool", "10", "Minimum number of connections in the connection pool")
	flag.StringVar(&maxPool, "maxPool", "50", "Maximum number of connections in the connection pool")
	flag.Parse()

	minPoolInt, err := strconv.Atoi(minPool)
	if err != nil {
		panic("error parsing minPool")
	}
	maxPoolInt, err := strconv.Atoi(maxPool)
	if err != nil {
		panic("error parsing maxPool")
	}

	mongoConfig := internal.NewMongoConfig(
		fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, hostName, port),
		uint64(minPoolInt),
		uint64(maxPoolInt),
		30,
	)

	client, cancel, err := internal.NewMongoClient(mongoConfig)
	if err != nil {
		panic("error connecting to mongodb")
	}

	defer func() {
		cancel()
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("mongodb disconnect error : %v", err)
		}
	}()

	repo := internal.NewMongoRepository(
		client,
	)



	// add instrumentation
	// reg := prometheus.NewRegistry()
	requestCountAdd := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "request_count_add",
		Help:      "Number of requests for add",
		Subsystem: "logingestor",
		Namespace: "api",
	})

	requestCountGet := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "request_count_get",
		Help:      "Number of requests for get",
		Subsystem: "logingestor",
		Namespace: "api",
	})

	requestLatency := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "request_latency",
		Help:      "Request latency",
		Subsystem: "logingestor",
		Namespace: "api",
	})
	reg := prometheus.NewRegistry()
	reg.MustRegister(requestCountAdd)
	reg.MustRegister(requestCountGet)
	reg.MustRegister(requestLatency)

	service := internal.NewService(repo)
	service = internal.NewInstrumentation(requestLatency, requestCountAdd, requestCountGet)(service)
	

	server := internal.NewServer(service, reg)


	go func() {
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("Server is shutting down...")
}
