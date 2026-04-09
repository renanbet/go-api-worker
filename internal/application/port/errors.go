package port

const errorNotFound = "not found"

type ErrNotFound struct{}

func (e ErrNotFound) Error() string { return errorNotFound }