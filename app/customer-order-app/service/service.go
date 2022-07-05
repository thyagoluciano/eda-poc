package service

import (
	"br.com.thyagoluciano.poc/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type service struct {
	logger   log.Logger
	producer domain.Producer
}

func NewServiceTask(logger log.Logger, producer domain.Producer) (domain.Service, error) {
	return &service{
		logger:   logger,
		producer: producer,
	}, nil
}

func (s service) Consumer(message *domain.Message) error {
	event, err := domain.Event{}.ToEvent(message)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	payload, err := toPayload(event.Payload)
	if err != nil {
		fmt.Println(err)
	}

	s.Producer(ctx, payload)

	fmt.Println(event)

	return nil
}

func toPayload(payload interface{}) (domain.PurchaseOrderCheckedOut, error) {
	str := fmt.Sprintf("%v", payload)
	purchaseOrder := domain.PurchaseOrderCheckedOut{}
	err := json.Unmarshal([]byte(str), &purchaseOrder)
	if err != nil {
		return domain.PurchaseOrderCheckedOut{}, err
	}

	return purchaseOrder, nil
}

func (s service) Producer(ctx context.Context, cmd domain.Command) error {

	logger := log.With(s.logger, "method", "CUSTOMER_ORDER_APP")

	customerOrderCmd, ok := cmd.(domain.PurchaseOrderCheckedOut)
	if !ok {
		logger.Log("Error to cast type", customerOrderCmd.CustomerId)
	}

	event := domain.Event{
		EventId:   domain.GenerateId(),
		EventType: "CUSTOMER_ORDER_CREATED",
		Metadata: domain.Metadata{
			Context: domain.Context{
				SpanId:       "",
				TraceId:      "",
				Application:  "",
				Channel:      "",
				Organization: "",
				Custom:       "",
			},
			Domain:     customerOrderCmd.ModuleName,
			ExternalId: customerOrderCmd.ExternalId,
			Id:         customerOrderCmd.Id,
			Timestamp:  time.Now().UTC(),
		},
		Payload: customerOrderCmd,
	}

	if err := s.producer.Publish(customerOrderCmd.Id, event); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}

	logger.Log("create customer order", customerOrderCmd.Id)
	return nil
}
