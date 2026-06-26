package captcha

import (
	"sync"
	"time"
)

// memStore 进程内验证码存储；带过期时间
type memStore struct {
	mu   sync.Mutex
	data map[string]item
	ttl  time.Duration
}

type item struct {
	answer string
	exp    time.Time
}

func newMemStore() *memStore {
	return &memStore{data: make(map[string]item), ttl: 3 * time.Minute}
}

func (s *memStore) Set(id, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = item{answer: val, exp: time.Now().Add(s.ttl)}
	// 顺手清理过期项
	for k, v := range s.data {
		if time.Now().After(v.exp) {
			delete(s.data, k)
		}
	}
	return nil
}

func (s *memStore) Get(id string, clear bool) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[id]
	if !ok || time.Now().After(v.exp) {
		delete(s.data, id)
		return ""
	}
	if clear {
		delete(s.data, id)
	}
	return v.answer
}

func (s *memStore) Verify(id, answer string, clear bool) bool {
	want := s.Get(id, clear)
	return want != "" && want == answer
}
