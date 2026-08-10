// Package example example 插件仓储层
package repo

import (
	exampleModel "go-admin/server/plugin/example/model"

	"gorm.io/gorm"
)

// NoteRepo 笔记仓储接口
type NoteRepo interface {
	Create(n *exampleModel.Note) error
	Delete(id uint) error
	List() ([]exampleModel.Note, error)
}

// NewNoteRepo 构造 NoteRepo
func NewNoteRepo(db *gorm.DB) NoteRepo { return &noteRepo{db: db} }

type noteRepo struct{ db *gorm.DB }

func (r *noteRepo) Create(n *exampleModel.Note) error {
	return r.db.Create(n).Error
}

func (r *noteRepo) Delete(id uint) error {
	return r.db.Delete(&exampleModel.Note{}, id).Error
}

func (r *noteRepo) List() ([]exampleModel.Note, error) {
	var list []exampleModel.Note
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
