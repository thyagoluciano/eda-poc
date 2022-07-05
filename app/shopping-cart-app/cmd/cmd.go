package cmd

import (
	"br.com.thyagoluciano.poc/app/shopping-cart-app/endpoint"
	"br.com.thyagoluciano.poc/app/shopping-cart-app/service"
	"br.com.thyagoluciano.poc/app/shopping-cart-app/transport"
	"br.com.thyagoluciano.poc/infra/kafka/producer"
	"br.com.thyagoluciano.poc/infra/kafka/types"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const topic = "CUSTOMER_ORDER"

func Start() {
	var httpAddr = flag.String("http", ":8080", "http listen address")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service", "customerOrder",
			"time:", log.DefaultTimestampUTC,
			"caller:", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	producer, err := producer.New([]string{"localhost:9092"}, topic, types.Sarama)
	if err != nil {
		level.Error(logger).Log("msg", "Error create producer: %v", err)
		panic(err)
	}

	flag.Parse()
	ctx := context.Background()

	srv := service.NewCommandApplicationService(producer, logger)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := endpoint.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := transport.NewHttpServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
