package utils

import "errors"

// RoundUpToQuotient 向上取整的商
func RoundUpToQuotient(dividend, divisor int) (int, error) {
	if dividend == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	quotient := dividend / divisor
	if dividend%divisor != 0 {
		quotient++
	}
	return quotient, nil
}
