package client

import (
	"errors"
	"fmt"

	"github.com/hamburghammer/grcon"
)

func newGrconClientError(act grcon.Action, err error) GrconClientError {
	return GrconClientError{
		Act: act,
		Err: err,
	}
}

// GrconClientError is a generic error that provides default implementations for the GrconError interface in the util module.
type GrconClientError struct {
	Err error
	Act grcon.Action
}

func (grue GrconClientError) Error() string {
	return fmt.Sprintf("grcon-client: on %s: %s", grue.Action(), grue.Err.Error())
}

func (grue GrconClientError) Action() grcon.Action {
	return grue.Act
}

func newInvalidResponseTypeError(expected, actual grcon.PacketType) InvalidResponseTypeError {
	return InvalidResponseTypeError{
		GrconClientError: newGrconClientError(grcon.Read, fmt.Errorf("invalid response type: expected %d but got %d", expected, actual)),
		Expected:         expected,
		Actual:           actual,
	}
}

type InvalidResponseTypeError struct {
	GrconClientError
	Expected grcon.PacketType
	Actual   grcon.PacketType
}

func newAuthFailedError() AuthFailedError {
	return AuthFailedError{
		newGrconClientError(grcon.Read, errors.New("authentication failed")),
	}
}

type AuthFailedError struct {
	GrconClientError
}

func newResponseIdMismatchError(expected, actual grcon.PacketId) ResponseIdMismatchError {
	return ResponseIdMismatchError{
		GrconClientError: newGrconClientError(grcon.Read, errors.New("invalid response type")),
		Expected:         expected,
		Actual:           actual,
	}
}

type ResponseIdMismatchError struct {
	GrconClientError
	Expected grcon.PacketId
	Actual   grcon.PacketId
}

func newResponseBodyError(expected, actual string) ResponseBodyError {
	return ResponseBodyError{
		newGrconClientError(
			grcon.Read,
			fmt.Errorf("response body error: expected '%s' got '%s'", expected, actual),
		),
	}
}

type ResponseBodyError struct {
	GrconClientError
}
