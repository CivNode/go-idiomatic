package ctxfirstbad

import "context"

// ctx second: should flag.
func Do(id string, ctx context.Context) error { // want `context.Context must be the first parameter`
	_ = ctx
	_ = id
	return nil
}

type Service struct{}

// ctx second on a method: should flag.
func (s *Service) Fetch(key string, ctx context.Context) (string, error) { // want `context.Context must be the first parameter`
	_ = ctx
	return key, nil
}
