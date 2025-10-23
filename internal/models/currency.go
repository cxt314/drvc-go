package models

import (
	"fmt"
	"strconv"
	"strings"
)

// USD represents US dollar amount in terms of cents
type USD int64

// ToUSD converts a float64 to USD
// e.g. 1.23 to $1.23, 1.345 to $1.35
func ToUSD(f float64) USD {
	return USD((f * 100) + 0.5)
}

func StrToUSD(s string) USD {
	f, err := strconv.ParseFloat(strings.ReplaceAll(s, "$", ""), 64)
	if err != nil {
		return ToUSD(0.0)
	}

	return ToUSD(f)
}

// Float64 converts a USD to float64
func (m USD) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

// Multiply safely multiplies a USD value by a float64, rounding
// to the nearest cent.
func (m USD) Multiply(f float64) USD {
	x := (float64(m) * f) + 0.5
	return USD(x)
}

func (m USD) AddUSD(a USD) USD {
	x := m.Float64()
	y := a.Float64()

	return ToUSD(x + y)
}

// String returns a formatted USD value
func (m USD) String() string {
	x := float64(m)
	x = x / 100
	return fmt.Sprintf("$%.2f", x)
}
