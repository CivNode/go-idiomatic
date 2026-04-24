package preferrangeintok

// Countdown must not flag.
func Countdown(xs []int) {
	for i := len(xs) - 1; i >= 0; i-- {
		_ = xs[i]
	}
}

// Step greater than 1 must not flag.
func EveryOther(xs []int) {
	for i := 0; i < len(xs); i += 2 {
		_ = xs[i]
	}
}

// Already using range.
func Ranged(xs []int) int {
	sum := 0
	for _, v := range xs {
		sum += v
	}
	return sum
}
