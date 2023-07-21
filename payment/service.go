package payment

import (
	"bwastartup/user"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
}

func NewService() *service{
	return &service{}
}

func(s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error){

	midtrans.ServerKey = "SB-Mid-server-fYSW5gTaJPsGDkyGQ2ItGKBr"
	midtrans.Environment = midtrans.Sandbox

	snapReq := &snap.Request{
		CustomerDetail: &midtrans.CustomerDetails{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}
	snapResp, err := snap.CreateTransaction(snapReq)
	if err != nil {
		return "", err
	}
	return snapResp.RedirectURL, nil
}
