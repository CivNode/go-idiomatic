package anyoveremptybad

type Box struct {
	V interface{} // want `use any instead of interface\{\}`
}

func Identity(v interface{}) interface{} { // want `use any instead of interface\{\}` `use any instead of interface\{\}`
	return v
}
