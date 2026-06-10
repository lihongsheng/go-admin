// Package example example 插件 DTO
package example

// NoteCreateReq 新建笔记
type NoteCreateReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
