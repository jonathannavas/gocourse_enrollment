package enrollment

import (
	"errors"
	"fmt"
)

var errUserIDRequired = errors.New("user id is required")
var errCourseIDRequired = errors.New("course id is required")
var errStatusRequired = errors.New("status invalid")

type ErrNotFound struct {
	enrollmentID string
}

type ErrInvalidStatus struct {
	Status string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Enrollment '%s' doesn't exist", e.enrollmentID)
}

func (e ErrInvalidStatus) Error() string {
	return fmt.Sprintf("Invalid '%s' status", e.Status)
}
