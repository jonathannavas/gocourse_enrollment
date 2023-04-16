package enrollment

import (
	"errors"
	"fmt"
)

var errUserIDRequired = errors.New("user id is required")
var errCourseIDRequired = errors.New("course id is required")
var errStatusRequired = errors.New("status is required")

type ErrNotFound struct {
	courseID string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Enrollment '%s' doesn't exist", e.courseID)
}
