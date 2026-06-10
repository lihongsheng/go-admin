// Package captcha 图形验证码：基于 base64Captcha 的内存 store
package captcha

import (
	"sync"
	"time"

	"github.com/mojocn/base64Captcha"
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

var (
	store   = &memStore{data: make(map[string]item), ttl: 3 * time.Minute}
	driver  = base64Captcha.NewDriverDigit(60, 180, 5, 0.7, 80) // 高度/宽度/位数/扭曲/最大噪点
	captcha = base64Captcha.NewCaptcha(driver, store)
)

// Generate 生成新的验证码，返回 id 和 base64 图片
func Generate() (id, b64 string, err error) {
	id, b64, _, err = captcha.Generate()
	return
}

// Verify 校验答案；clear=true 用过即弃
func Verify(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	return store.Verify(id, answer, true)
}
