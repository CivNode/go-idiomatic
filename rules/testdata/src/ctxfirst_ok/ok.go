package ctxfirstok

import "context"

// ctx first: fine.
func Do(ctx context.Context, id string) error {
	_ = ctx
	_ = id
	return nil
}

// No context at all: fine.
func NoCtx(id string) error {
	_ = id
	return nil
}

// Method with receiver plus ctx first: fine.
type Service struct{}

func (s *Service) Fetch(ctx context.Context, key string) (string, error) {
	_ = ctx
	return key, nil
}
