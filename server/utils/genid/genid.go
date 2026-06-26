// Package genid 雪花算法 ID 生成器，从 utils 中独立出来以避免循环依赖
package genid

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// GenDeviceID 全局雪花 ID 生成器
var GenDeviceID = NewGenID(1758357486530)

// genID 雪花算法生成器
type genID struct {
	mu            sync.Mutex
	epoch         int64
	lastTimestamp int64
	sequence      int64
	random        *rand.Rand
}

// NewGenID 创建雪花算法实例
func NewGenID(epoch int64) *genID {
	return &genID{
		epoch:  epoch,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate0X 生成唯一 ID 并返回 16 进制字符串
func (s *genID) Generate0X() string {
	return fmt.Sprintf("%x", s.Generate())
}

// Generate 生成唯一 ID
func (s *genID) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli()
	rollbackFlag := int64(0)

	if timestamp < s.lastTimestamp {
		rollbackFlag = 1
		random := int64(s.random.Intn(128))
		s.sequence++
		if s.sequence > 511 {
			s.sequence = 0
		}
		id := ((s.lastTimestamp - s.epoch) << 17) |
			(random << 10) |
			(s.sequence << 1) |
			rollbackFlag
		return id
	}

	if timestamp == s.lastTimestamp {
		s.sequence++
		if s.sequence > 511 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixMilli()
			}
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp
	random := int64(s.random.Intn(128))

	id := ((timestamp - s.epoch) << 17) |
		(random << 10) |
		(s.sequence << 1) |
		rollbackFlag
	return id
}
