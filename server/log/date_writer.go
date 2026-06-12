package log

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// dateWriter 按日期分目录，每个目录下用 lumberjack 做大小轮转。
// 启动时清理超过 7 天的日期目录。
type dateWriter struct {
	dir        string
	maxSize    int
	maxBackups int
	maxAge     int

	mu          sync.Mutex
	currentDate string
	writer      io.WriteCloser
}

func NewDateWriter(dir string, maxSize, maxBackups, maxAge int) *dateWriter {
	dw := &dateWriter{dir: dir, maxSize: maxSize, maxBackups: maxBackups, maxAge: maxAge}
	dw.clean()
	return dw
}

func (w *dateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != w.currentDate {
		if w.writer != nil {
			w.writer.Close()
		}
		dayDir := filepath.Join(w.dir, today)
		os.MkdirAll(dayDir, 0o755)
		w.writer = &lumberjack.Logger{
			Filename:   filepath.Join(dayDir, "server.log"),
			MaxSize:    w.maxSize,
			MaxBackups: w.maxBackups,
			MaxAge:     w.maxAge,
			Compress:   false,
		}
		w.currentDate = today
		go w.clean()
	}
	return w.writer.Write(p)
}

func (w *dateWriter) clean() {
	entries, _ := os.ReadDir(w.dir)
	cutoff := time.Now().AddDate(0, 0, -7)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t, err := time.Parse("2006-01-02", e.Name())
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			os.RemoveAll(filepath.Join(w.dir, e.Name()))
		}
	}
}
