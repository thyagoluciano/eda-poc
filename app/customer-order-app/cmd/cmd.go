package cmd

import (
	"br.com.thyagoluciano.poc/app/customer-order-app/service"
	"br.com.thyagoluciano.poc/domain"
	"br.com.thyagoluciano.poc/infra/kafka/consumer"
	"br.com.thyagoluciano.poc/infra/kafka/producer"
	"br.com.thyagoluciano.poc/infra/kafka/types"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

func fatal(logger log.Logger, err error) {
	_ = level.Error(logger).Log("err", err)
	os.Exit(1)
}

func Start() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	{
		ctx = context.Background()
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	_ = logger.Log("msg", "initializing services")

	kafkaURL := "localhost:9092"
	topic := "PURCHASE_ORDER_CHECKED_OUT"
	groupId := "CUSTOMER_ORDER_APP"
	producerTopic := "CUSTOMER_ORDER"

	var srv domain.Service
	{
		producer, err := producer.New([]string{kafkaURL}, producerTopic, types.Sarama)
		if err != nil {
			level.Error(logger).Log("msg", "Error create producer: %v", err)
		}
		getProductOnCimService, err := service.NewServiceTask(
			log.With(logger, "component", "CustomerOrderApp"),
			producer,
		)
		if err != nil {
			fatal(logger, fmt.Errorf("failed to init storage: %w", err))
		}
		srv = getProductOnCimService
	}

	consumer, err := consumer.NewConsumer([]string{kafkaURL}, topic, groupId, types.Sarama)
	if err != nil {
		level.Error(logger).Log("msg", "Error to initialize producer: %v", err)
	}

	_ = logger.Log("msg", "initializing kafka consumer")

	go func() {
		//errs <- consumer.Subscribe(srv.Consumer)
		consumer.Subscribe(srv.Consumer)
	}()
}
