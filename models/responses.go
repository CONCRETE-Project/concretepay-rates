package models

import "github.com/shopspring/decimal"

type BaseResponse struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
	Error  string      `json:"error"`
}

type Rate struct {
	Code string          `json:"code"`
	Rate decimal.Decimal `json:"rate"`
}

type RatesResponse []Rate
