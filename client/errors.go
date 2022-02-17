package client

import (
	"errors"
	"fmt"

	"github.com/hamburghammer/grcon"
)

func newGrconUtilError(act grcon.Action, err error) GrconUtilError {
	return GrconUtilError{
		Act: act,
		Err: err,
	}
}

// GrconUtilError is a generic error that provides default implementations for the GrconError interface in the util module.
type GrconUtilError struct {
	Err error
	Act grcon.Action
}

func (grue GrconUtilError) Error() string {
	return fmt.Sprintf("grcon-client: on %s: %s", grue.Action(), grue.Err.Error())
}

func (grue GrconUtilError) Action() grcon.Action {
	return grue.Act
}

func newInvalidResponseTypeError(expected, actual grcon.PacketType) InvalidResponseTypeError {
	return InvalidResponseTypeError{
		GrconUtilError: newGrconUtilError(grcon.Read, fmt.Errorf("invalid response type: expected %d but got %d", expected, actual)),
		Expected:       expected,
		Actual:         actual,
	}
}

type InvalidResponseTypeError struct {
	GrconUtilError
	Expected grcon.PacketType
	Actual   grcon.PacketType
}

func newAuthFailedError() AuthFailedError {
	return AuthFailedError{
		newGrconUtilError(grcon.Read, errors.New("authentication failed")),
	}
}

type AuthFailedError struct {
	GrconUtilError
}

func newResponseIdMismatchError(expected, actual grcon.PacketId) ResponseIdMismatchError {
	return ResponseIdMismatchError{
		GrconUtilError: newGrconUtilError(grcon.Read, errors.New("invalid response type")),
		Expected:       expected,
		Actual:         actual,
	}
}

type ResponseIdMismatchError struct {
	GrconUtilError
	Expected grcon.PacketId
	Actual   grcon.PacketId
}

func newResponseBodyError(expected, actual string) ResponseBodyError {
	return ResponseBodyError{
		newGrconUtilError(
			grcon.Read,
			fmt.Errorf("response body error: expected '%s' got '%s'", expected, actual),
		),
	}
}

type ResponseBodyError struct {
	GrconUtilError
}
