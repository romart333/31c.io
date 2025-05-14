package currencyUtil

import (
	"github.com/shopspring/decimal"
)

const (
	SatoshiInBTC = 100000000
	BTCPrecision = 8
)

func FromSatoshi(v interface{}) decimal.Decimal {
	value := toDecimal(v)
	return value.Div(decimal.NewFromInt(SatoshiInBTC)).Round(BTCPrecision)
}

func ToSatoshi(v interface{}) decimal.Decimal {
	value := toDecimal(v)
	return value.Mul(decimal.NewFromInt(SatoshiInBTC)).Round(0)
}

func RoundValue(value interface{}) decimal.Decimal {
	bf := toDecimal(value)
	return bf.Round(BTCPrecision)
}

func toDecimal(v interface{}) decimal.Decimal {
	switch val := v.(type) {
	case string:
		d, err := decimal.NewFromString(val)
		if err != nil {
			return decimal.Zero
		}
		return d
	case int, int64:
		return decimal.NewFromInt(val.(int64))
	case float64:
		return decimal.NewFromFloat(val)
	default:
		return decimal.Zero
	}
}
