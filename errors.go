package grcon

import (
	"errors"
	"fmt"
)

// Action is a type to indicate the part in which an error occurred.
type Action string

// Actions for error checks.
const (
	// Write indicates that it happened on a write operation.
	Write Action = "write"

	// Read indicates that it happened on a read operation.
	Read Action = "read"
)

// GrconError is the interface all errors from this packet implement.
// You can use this interface for errors.Is checks.
type GrconError interface {
	// Default error interface.
	error
	// Action returns the action that produced the error.
	Action() Action
}

func newGrconGenericError(act Action, err error) GrconGenericError {
	return GrconGenericError{
		Act: act,
		Err: err,
	}
}

// GrconGenericError is a generic error that provides default implementations for the grconError interface.
type GrconGenericError struct {
	Err error
	Act Action
}

// Error returns the error in string format.
func (rge GrconGenericError) Error() string {
	return fmt.Sprintf("grcon: on %s: %s", rge.Action(), rge.Err.Error())
}

// Action returns the action where the error was thrown.
func (rge GrconGenericError) Action() Action {
	return rge.Act
}

func newUnexpectedFormatError() UnexpectedFormatError {
	return UnexpectedFormatError{
		newGrconGenericError(
			Read,
			errors.New("unexpected response format: the packet is smaller than the minimum size"),
		),
	}
}

// UnexpectedFormatError occurres when the packet size is smaller than the minimum size.
// This indicates a wrongly composed/formatted packet.
type UnexpectedFormatError struct {
	GrconGenericError
}

func newRequestTooLongError() RequestTooLongError {
	return RequestTooLongError{
		newGrconGenericError(
			Write,
			errors.New("request body is too long"),
		),
	}
}

// RequestTooLongError occurres when the length of a packet is to big.
// This indicates that the body is too long.
type RequestTooLongError struct {
	GrconGenericError
}

func newResponseTooLongError() ResponseTooLongError {
	return ResponseTooLongError{
		newGrconGenericError(
			Read,
			errors.New("response body is too long"),
		),
	}
}

// ResponseTooLongError occurres when the size of a packet is to big.
// This indicates a wrongly composed/formatted packet.
type ResponseTooLongError struct {
	GrconGenericError
}
