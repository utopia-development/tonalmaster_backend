package calendars

// mod returns a non-negative remainder, including for negative values.
func mod(value, divisor int64) int {
	if divisor <= 0 {
		panic("mod divisor must be positive")
	}
	remainder := value % divisor
	if remainder < 0 {
		remainder += divisor
	}
	return int(remainder)
}
