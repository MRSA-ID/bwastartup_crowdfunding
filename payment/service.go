package payment

import (
	"bwastartup/campaign"
	"bwastartup/transaction"
	"bwastartup/user"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
	transactionRepository transaction.Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService(transactionRepository transaction.Repository, campaignRepository campaign.Repository) *service{
	return &service{transactionRepository, campaignRepository}
}

func(s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error){

	midtrans.ServerKey = ""
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

func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) error{
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.transactionRepository.GetByID(transaction_id)
	if err != nil {
		return err
	}

	if (input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept"){
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement"{
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel"{
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid"{
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}
	return nil
}