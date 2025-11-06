package response

type E interface {
	IsAppError() bool
	Code() int
	error
}
