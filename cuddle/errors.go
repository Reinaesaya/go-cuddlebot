package cuddle

type Error struct {
	OK      bool   `json:"ok"`
	Message string `json:"error,omitempty"`
}

var (
	InvalidAddressError  = &Error{Message: "InvalidAddressError"}
	InvalidMessageError  = &Error{Message: "InvalidMessageError"}
	InvalidSetpointError = &Error{Message: "InvalidSetpointError"}
	MethodNotAllowed     = &Error{Message: "MethodNotAllowed"}
	MissingFieldError    = &Error{Message: "MissingFieldError"}
	NotImplementedError  = &Error{Message: "NotImplementedError"}
)

func (e *Error) Error() string {
	return e.Message
}
