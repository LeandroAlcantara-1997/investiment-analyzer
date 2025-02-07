package models

import "time"

type Trader struct {
	Quantity    int
	Price       float64
	Side        string
	ActionTime  time.Time
	CompanyName string
}

type CompanyPrice struct {
	PriceTime time.Time
	Price     float64
}
