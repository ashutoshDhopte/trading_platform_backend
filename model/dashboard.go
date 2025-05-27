package model

import "trading_platform_backend/orm"

type Dashboard struct {
	User     orm.Users
	Stocks   []orm.Stocks
	Holdings []orm.Holdings
}
