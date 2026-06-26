// Package captcha 图形验证码：支持内存 / Redis 存储
package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

var (
	driver  = base64Captcha.NewDriverDigit(60, 180, 5, 0.7, 80) // 高度/宽度/位数/扭曲/最大噪点
	store   base64Captcha.Store
	captcha *base64Captcha.Captcha
)

// Captcha 验证码配置
type Captcha struct {
	// Drive 存储驱动：memory / redis，默认 memory
	Drive string `mapstructure:"drive" json:"drive" yaml:"drive"`
}

func init() {
	// 默认使用内存存储
	store = newMemStore()
	captcha = base64Captcha.NewCaptcha(driver, store)
}

// SetStore 替换验证码存储实现；由 initialize 在启动时按配置调用
func SetStore(s base64Captcha.Store) {
	store = s
	captcha = base64Captcha.NewCaptcha(driver, store)
}

// SetupRedisStore 切换到 Redis 存储；cli 为已连接的 Redis 客户端
func SetupRedisStore(cli *redis.Client) {
	SetStore(newRedisStore(cli))
}

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
