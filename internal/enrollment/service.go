package enrollment

import (
	"context"
	"log"

	"github.com/jonathannavas/gocourse_domain/domain"

	courseSdk "github.com/jonathannavas/go_course_sdk/course"
	userSdk "github.com/jonathannavas/go_course_sdk/user"
)

type (
	Service interface {
		Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
		Update(ctx context.Context, id string, status *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	service struct {
		log         *log.Logger
		repo        Repository
		userTrans   userSdk.Transport
		courseTrans courseSdk.Transport
	}

	Filters struct {
		UserID   string
		CourseID string
	}
)

func NewService(log *log.Logger, repo Repository, courseTrans courseSdk.Transport, userTrans userSdk.Transport) Service {
	return &service{
		log:         log,
		repo:        repo,
		userTrans:   userTrans,
		courseTrans: courseTrans,
	}
}

func (s service) Create(ctx context.Context, userID, courseID string) (*domain.Enrollment, error) {

	enrollment := &domain.Enrollment{
		UserID:   userID,
		CourseId: courseID,
		Status:   domain.Pending,
	}

	if _, err := s.userTrans.Get(userID); err != nil {
		return nil, err
	}

	if _, err := s.courseTrans.Get(courseID); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	enrollments, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (s service) Update(ctx context.Context, id string, status *string) error {

	if status != nil {
		switch domain.EnrollmentStatus(*status) {
		case domain.Pending, domain.Active, domain.Studying, domain.Inactive:
		default:
			return ErrInvalidStatus{*status}
		}
	}

	if err := s.repo.Update(ctx, id, status); err != nil {
		return err
	}
	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
