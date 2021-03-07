package code

type ErrCode int32

const (
	OK ErrCode = iota
	InvalidArgument
	NotFound
	Internal
)
