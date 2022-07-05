package service

import (
	"br.com.thyagoluciano.poc/domain"
	"context"
	"github.com/go-kit/kit/log/level"
	"time"

	"github.com/go-kit/log"
)

type PurchaseOrderService interface {
	CheckedOut(ctx context.Context, checkedOutCmd domain.Command) error
}

// service implements the Command Application Service
type purchaseApplicationService struct {
	producer domain.Producer
	logger   log.Logger
}

func NewCommandApplicationService(producer domain.Producer, logger log.Logger) PurchaseOrderService {
	return &purchaseApplicationService{
		producer: producer,
		logger:   logger,
	}
}

func (s purchaseApplicationService) CheckedOut(
	ctx context.Context,
	cmd domain.Command,
) error {
	// TODO: Adicionar os headers no contexto para adicionar no payload do kafka

	logger := log.With(s.logger, "method", "CreateCustomerOrder")

	purchaseOrderCheckedOutCmd, ok := cmd.(domain.PurchaseOrderCheckedOut)
	if !ok {
		logger.Log("Error to cast type", purchaseOrderCheckedOutCmd.Id)
	}

	event := domain.Event{
		EventId:   domain.GenerateId(),
		EventType: "PURCHASE_ORDER_CHECKED_OUT",
		Metadata: domain.Metadata{
			Context: domain.Context{
				SpanId:       "",
				TraceId:      "",
				Application:  "",
				Channel:      "",
				Organization: "",
				Custom:       "",
			},
			Domain:     purchaseOrderCheckedOutCmd.ModuleName,
			ExternalId: purchaseOrderCheckedOutCmd.ExternalId,
			Id:         purchaseOrderCheckedOutCmd.Id,
			Timestamp:  time.Now().UTC(),
		},
		Payload: purchaseOrderCheckedOutCmd,
	}

	if err := s.producer.Publish(purchaseOrderCheckedOutCmd.Id, event); err != nil {
		level.Error(logger).Log("err", err)
		return err
	}

	logger.Log("create customer order", purchaseOrderCheckedOutCmd.Id)

	return nil
}
