package transaction

import (
	"bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     string
	Status     string
	Code       string
	User 			 user.User
	CreatedAt  time.Time
	UpdateAt   time.Time
}