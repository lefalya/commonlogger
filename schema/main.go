package schema

type CommonError struct {
	Id          string
	Context     string
	Err         error
	ErrResponse error
}
