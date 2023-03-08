package lib

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type CoreApi struct {
	ca coreapi.Client
}

func NewMidtransDriver() CoreApi {
	return CoreApi{ca: coreapi.Client{}}
}

func (c *CoreApi) HandleNotification(id string) (*coreapi.TransactionStatusResponse, error) {
	c.ca.New(os.Getenv("AUTHORIZATION_VALUE"), midtrans.Sandbox)

	midtransReport, err := c.ca.CheckTransaction(id)
	if err != nil {
		return midtransReport, err
	}

	return midtransReport, nil
}

// func (c *CoreApi) GetDetailAction(id string) (*coreapi.PaymentAccountResponse, error) {
// 	c.ca.New(os.Getenv("AUTHORIZATION_VALUE"), midtrans.Sandbox)

// 	midtrans, err := c.ca.ChargeTransaction()
// }