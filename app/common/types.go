package common

// A zero-sized type
type Void struct{}

var Voidance Void

type EmptyParamErr struct{}

func (e EmptyParamErr) Error() string {
	return "empty search param"
}
