// Package installer 在线初始化模块：检测 DB / 注册 Model / AutoMigrate / 写入种子数据
package installer

import "sync"

// registry 待迁移的全部 Model 指针（核心 + 插件）
var (
	regMu    sync.RWMutex
	registry []interface{}
	// seeds 注册的种子数据回调（按注册顺序执行）
	seedsMu sync.RWMutex
	seeds   []SeedFn
)

// SeedFn 种子数据写入函数
type SeedFn func(deps SeedDeps) error

// Register 注册 Model（启动期 / 插件加载期调用）
func Register(models ...interface{}) {
	regMu.Lock()
	defer regMu.Unlock()
	registry = append(registry, models...)
}

// RegisterSeed 注册种子数据回调
func RegisterSeed(fn SeedFn) {
	seedsMu.Lock()
	defer seedsMu.Unlock()
	seeds = append(seeds, fn)
}

// AllModels 返回已注册的全部 Model（拷贝，防并发）
func AllModels() []interface{} {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]interface{}, len(registry))
	copy(out, registry)
	return out
}

// AllSeeds 返回所有种子回调
func AllSeeds() []SeedFn {
	seedsMu.RLock()
	defer seedsMu.RUnlock()
	out := make([]SeedFn, len(seeds))
	copy(out, seeds)
	return out
}
