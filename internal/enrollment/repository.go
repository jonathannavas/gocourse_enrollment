package enrollment

import (
	"context"
	"log"

	"github.com/jonathannavas/gocourse_domain/domain"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, enrollment *domain.Enrollment) error
	GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error)
	Update(ctx context.Context, id string, status *string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(ctx context.Context, enrollment *domain.Enrollment) error {

	if err := repo.db.Debug().WithContext(ctx).Create(enrollment).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("Enrollment created with id: ", enrollment.ID)
	return nil
}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment

	tx := repo.db.WithContext(ctx).Model(&enrollments)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&enrollments)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return nil, result.Error
	}
	return enrollments, nil
}

func (repo *repo) Update(ctx context.Context, id string, status *string) error {
	values := make(map[string]interface{})
	if status != nil {
		values["status"] = *status
	}

	result := repo.db.WithContext(ctx).Model(&domain.Enrollment{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound{id}
	}
	return nil

}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.UserID != "" {
		tx = tx.Where("user_id = ?", filters.UserID)
	}

	if filters.CourseID != "" {
		tx = tx.Where("course_id = ?", filters.CourseID)
	}
	return tx
}

func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.Enrollment{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		repo.log.Println(err)
		return 0, nil
	}
	return int(count), nil
}
