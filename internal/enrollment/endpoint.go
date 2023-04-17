package enrollment

import (
	"context"
	"errors"

	"github.com/jonathannavas/go_lib_response/response"
	"github.com/jonathannavas/gocourse_meta/meta"

	courseSdk "github.com/jonathannavas/go_course_sdk/course"
	userSdk "github.com/jonathannavas/go_course_sdk/user"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
	Endpoints  struct {
		Create Controller
		GetAll Controller
		Update Controller
	}

	CreateRequest struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

	GetAllRequest struct {
		UserID   string
		CourseID string
		Limit    int
		Page     int
	}

	UpdateRequest struct {
		ID     string
		Status *string `json:"status"`
	}

	Config struct {
		LimitPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		enrollmentBody := request.(CreateRequest)

		if enrollmentBody.UserID == "" {
			return nil, response.BadRequest(errUserIDRequired.Error())
		}

		if enrollmentBody.CourseID == "" {
			return nil, response.BadRequest(errCourseIDRequired.Error())
		}

		enrollment, err := s.Create(ctx, enrollmentBody.UserID, enrollmentBody.CourseID)

		if err != nil {
			if errors.As(err, &userSdk.ErrNotFound{}) ||
				errors.As(err, &courseSdk.ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enrollment, nil, 200), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllRequest)

		filters := Filters{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimitPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		enrollments, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", enrollments, meta, 201), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		enrollmentBody := request.(UpdateRequest)

		if enrollmentBody.Status != nil && *enrollmentBody.Status == "" {
			return nil, response.BadRequest(errStatusRequired.Error())
		}

		if err := s.Update(ctx, enrollmentBody.ID, enrollmentBody.Status); err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			if errors.As(err, &ErrInvalidStatus{}) {
				return nil, response.BadRequest(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", nil, nil, 200), nil
	}
}
