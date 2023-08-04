package util

const (
	USD = "USD"
	EUR = "EUR"
	INR = "INR"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, EUR, INR:
		return true
	}

	return false
}
