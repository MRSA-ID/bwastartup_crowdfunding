package transaction

import "time"

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     string
	Status     string
	Code       string
	CreatedAt  time.Time
	UpdateAt   time.Time
}