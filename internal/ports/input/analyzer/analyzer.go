package dto

import (
	"encoding/json"
	"math"
	"strings"
	"time"
)

type AnalyzerRequest struct {
	InitialDate *Date `json:"initialDate" validate:"required" example:"2021-03-01 01:01:01"`
	FinalDate   *Date `json:"finalDate" validate:"required" example:"2021-03-30 01:01:01"`
	Interval    int64 `json:"interval" validate:"required"`
}

type AnalyzersResponse struct {
	Registers []Register `json:"registers"`
}

type Register struct {
	HeritageEvolution        float64 `json:"heritageEvolution"`
	AccumulatedProfitability float64 `json:"accumulatedProfitability"`
	Timestamp                *Date   `json:"timestamp"`
}

type Money float64

func (m Money) Round(precision int) float64 {
	fator := math.Pow(10, float64(precision))
	return math.Round(float64(m)*fator) / fator
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	t := strings.Trim(string(data), "\"")
	d.Time, err = time.Parse(time.DateTime, t)
	return
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(time.DateTime))
}
