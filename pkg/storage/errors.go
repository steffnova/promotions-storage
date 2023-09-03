package storage

import "fmt"

var (
	ErrOptionPeriod   = fmt.Errorf("period can't be less then 0")
	ErrLoadPromotions = fmt.Errorf("failed to load promotions")
)
