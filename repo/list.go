package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/kevin-luvian/moonlay-test/model"
)

type ListRepo struct {
	db *gorm.DB
}

func NewListRepo(db *gorm.DB) *ListRepo {
	return &ListRepo{
		db: db,
	}
}

func (r *ListRepo) Create(l model.List) (model.List, error) {
	err := r.db.Create(&l).Error
	return l, err
}

func (r *ListRepo) Update(l *model.List) (err error) {
	return r.db.Model(&model.List{}).Update(l).Error
}

func (r *ListRepo) GetByID(id uint) (model.List, error) {
	var m model.List
	if err := r.db.First(&m, id).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (r *ListRepo) GetAll(offset, limit int) ([]model.List, int, error) {
	var (
		lists []model.List
		count int
	)

	r.db.Model(&lists).Count(&count)
	q := r.db.Order("created_at desc")

	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}

	q.Find(&lists)

	return lists, count, nil
}

func (r *ListRepo) GetRoots(offset, limit int, preload bool) ([]model.List, int, error) {
	var (
		lists []model.List
		count int
	)

	r.db.Model(&lists).Where("level0_id IS NULL").Count(&count)
	q := r.db.Where("level0_id IS NULL").
		Order("created_at desc")

	if limit > 0 {
		q = q.Offset(offset).Limit(limit)
	}

	if preload {
		q = q.Preload("Sublists")
	}

	q.Find(&lists)

	return lists, count, nil
}

func (r *ListRepo) DeleteByID(id uint) error {
	if err := r.db.Delete(&model.List{}, id).Error; err != nil {
		return err
	}

	return nil
}
