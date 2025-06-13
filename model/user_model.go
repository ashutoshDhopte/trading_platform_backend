package model

type UserModel struct {
	UserID             int64
	Username           string
	Email              string
	CashBalanceDollars float64
	CreatedAt          string
	UpdatedAt          string
	NotificationsOn    bool
}
