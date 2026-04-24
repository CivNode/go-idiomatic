package anyoveremptyok

// any is the idiomatic form.
type Box struct {
	V any
}

// Non-empty interfaces must not flag.
type Stringer interface {
	String() string
}

func Identity(v any) any { return v }
