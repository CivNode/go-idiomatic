package errorsisasbad

import "errors"

var ErrSentinel = errors.New("boom")

func StringCompare(err error) bool {
	return err.Error() == "boom" // want `err.Error\(\) against a string literal`
}

func SentinelCompare(err error) bool {
	return err == ErrSentinel // want `compare errors with errors.Is`
}
