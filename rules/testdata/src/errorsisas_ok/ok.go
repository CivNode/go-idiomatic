package errorsisasok

import "errors"

var ErrSentinel = errors.New("boom")

// Comparing to nil is fine.
func NilCheck(err error) bool {
	return err == nil
}

// Using errors.Is is the whole point.
func IsCheck(err error) bool {
	return errors.Is(err, ErrSentinel)
}

// Comparing two non-error values is not our business.
func StringEq(a, b string) bool {
	return a == b
}
