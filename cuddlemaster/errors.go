package cuddlemaster

import "errors"

var AddressError = errors.New("Invalid address")
var MissingFieldsError = errors.New("Missing fields")
var InvalidSetpointError = errors.New("Invalid setpoint(s)")

type jsonError struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}
