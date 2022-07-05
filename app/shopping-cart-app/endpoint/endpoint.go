package endpoint

import (
	"br.com.thyagoluciano.poc/app/shopping-cart-app/endpoint/api"
	"br.com.thyagoluciano.poc/app/shopping-cart-app/service"
	"br.com.thyagoluciano.poc/domain"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	PurchaseOrderCheckedOut endpoint.Endpoint
}

func MakeEndpoints(s service.PurchaseOrderService) Endpoints {
	return Endpoints{
		PurchaseOrderCheckedOut: makePurchaseOrderCheckedOut(s),
	}
}

func makePurchaseOrderCheckedOut(svc service.PurchaseOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(api.PurchaseOrderCheckedOutRequest)

		command := toCommand(&req)

		err = svc.CheckedOut(ctx, command)

		if err != nil {

		}

		response = api.PurchaseOrderCheckedOutResponse{
			PurchaseOrderId: command.Id,
		}

		return response, nil
	}
}

func toCommand(request *api.PurchaseOrderCheckedOutRequest) domain.PurchaseOrderCheckedOut {
	return domain.PurchaseOrderCheckedOut{
		Id:         domain.GenerateId(),
		ModuleName: "shopping-cart",
		ExternalId: domain.GenerateId(),
		CustomerId: request.CustomerId,
		ProductId:  request.ProductId,
		Address: domain.Address{
			Cep:    request.Address.Cep,
			Number: request.Address.Number,
		},
		Payment: domain.Payment{
			Token: request.Payment.Token,
			Price: domain.Price{
				Amount:   request.Payment.Price.Amount,
				Currency: request.Payment.Price.Currency,
				Scale:    request.Payment.Price.Scale,
			},
		},
	}
}
