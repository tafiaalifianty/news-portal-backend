package models

type UserSpending struct {
	UserID        int64 `json:"user_id" gorm:"primaryKey"`
	Month         int   `json:"month" gorm:"primaryKey"`
	Year          int   `json:"year" gorm:"primaryKey"`
	TotalSpending int   `json:"total_spending"`
}

type UserSpendingTotalAggregates struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Fullname      string `json:"fullname"`
	TotalSpending int    `json:"total_spending"`
}
