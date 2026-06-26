// Package example example 插件业务服务
package service

import (
	dtoExample "github.com/lihongsheng/go-admin/server/plugin/example/dto"
	exampleModel "github.com/lihongsheng/go-admin/server/plugin/example/model"
	repoExample "github.com/lihongsheng/go-admin/server/plugin/example/repo"
)

// NoteService 笔记业务接口
type NoteService interface {
	Create(req dtoExample.NoteCreateReq) (*exampleModel.Note, error)
	Delete(id uint) error
	List() ([]exampleModel.Note, error)
}

// NewNoteService 构造 NoteService
func NewNoteService(repo repoExample.NoteRepo) NoteService {
	return &noteService{repo: repo}
}

type noteService struct {
	repo repoExample.NoteRepo
}

// DefaultNote 包级单例
var DefaultNote NoteService

func (s *noteService) Create(req dtoExample.NoteCreateReq) (*exampleModel.Note, error) {
	n := &exampleModel.Note{Title: req.Title, Content: req.Content}
	if err := s.repo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *noteService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *noteService) List() ([]exampleModel.Note, error) {
	return s.repo.List()
}
