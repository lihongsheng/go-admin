package utils

import (
	"crypto/md5"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

var GenDeviceID = NewGenID(1758357486530)

// genID 雪花算法生成器（时钟回滚标识在最后一位，修复序列号解析bug）
type genID struct {
	mu            sync.Mutex // 保证并发安全
	epoch         int64      // 起始时间戳（毫秒）
	lastTimestamp int64      // 上一次生成ID的时间戳
	sequence      int64      // 序列号（0-511）
	random        *rand.Rand // 随机数生成器
}

// NewGenID 创建雪花算法实例
// epoch: 起始时间戳（毫秒），默认2020-01-01 00:00:00
func NewGenID(epoch int64) *genID {
	return &genID{
		epoch:  epoch,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate0X 生成唯一ID并返回16进制字符串
func (s *genID) Generate0X() string {
	return fmt.Sprintf("%x", s.Generate())
}

// Generate 生成唯一ID
func (s *genID) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixMilli()
	rollbackFlag := int64(0)

	// 处理时钟回滚
	if timestamp < s.lastTimestamp {
		rollbackFlag = 1
		random := int64(s.random.Intn(128)) // 7位随机数（0-127）
		s.sequence++
		if s.sequence > 511 {
			s.sequence = 0
		}
		// 组合ID：时间戳(46) | 随机数(7) | 序列号(9) | 回滚标识(1)
		id := ((s.lastTimestamp - s.epoch) << 17) | // 46位时间戳差左移17位（7+9+1=17）
			(random << 10) | // 7位随机数左移10位（9+1=10）
			(s.sequence << 1) | // 9位序列号左移1位（避开回滚标识）
			rollbackFlag // 最后1位回滚标识
		return id
	}

	// 处理同一毫秒
	if timestamp == s.lastTimestamp {
		s.sequence++
		if s.sequence > 511 { // 9位最大为511（2^9-1）
			// 等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixMilli()
			}
			s.sequence = 0
		}
	} else {
		s.sequence = 0 // 新毫秒重置序列号
	}

	s.lastTimestamp = timestamp
	random := int64(s.random.Intn(128))

	// 正常生成ID
	id := ((timestamp - s.epoch) << 17) | // 46位时间戳差
		(random << 10) | // 7位随机数
		(s.sequence << 1) | // 9位序列号（左移1位避开回滚标识）
		rollbackFlag // 最后1位回滚标识
	return id
}

// ParseID 解析ID各部分信息（修复序列号解析逻辑）
func ParseID(id int64) map[string]int64 {
	return map[string]int64{
		// 时间戳差：取高46位（右移17位，再与46位掩码）
		"timestamp_diff": (id >> 17) & ((1 << 46) - 1),
		// 随机数：取中间7位（右移10位，再与7位掩码）
		"random": (id >> 10) & ((1 << 7) - 1),
		// 序列号：取中间9位（右移1位，再与9位掩码）
		"sequence": (id >> 1) & ((1 << 9) - 1),
		// 回滚标识：取最后1位
		"rollback_flag": id & 1,
	}
}

func GetGenIDTimestamp(id int64) time.Time {
	t := (id>>17)&((1<<46)-1) + 1758357486530
	return time.UnixMilli(t)
}

// CRC32Pool CRC32对象池（复用哈希对象，提升高并发性能）
var CRC32Pool = sync.Pool{
	New: func() interface{} {
		return crc32.NewIEEE()
	},
}

// HashTo32 通用哈希函数，输出严格在 0-32 之间
func HashTo32(data []byte) int {
	// 从池获取CRC32对象
	h := CRC32Pool.Get().(hash.Hash32)
	defer func() {
		// 重置对象并放回池
		h.Reset()
		CRC32Pool.Put(h)
	}()
	// 步骤1：计算CRC32哈希（IEEE标准，分布均匀，性能高）
	_, _ = h.Write(data)
	crc := h.Sum32()
	return int(crc % 33) // 输出0-32
}

// HashStringTo32 便捷封装：针对不同输入类型的重载（可选）
func HashStringTo32(s string) int {
	return HashTo32([]byte(s))
}

var nodeID int
var nodeOnce sync.Once

func GetNodeID() int {
	nodeOnce.Do(func() {
		// 创建 Snowflake 节点
		host := os.Getenv("HOSTNAME")
		if host == "" {
			host = RandomString(16)
		}
		nodeID = HashStringTo32(host)
	})
	return nodeID
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenTradeNo(mchNO string) (string, int64) {
	gId := GenDeviceID.Generate()
	return fmt.Sprintf("%d%s%s%s",
		gId,
		fmt.Sprintf("%02d", HashStringTo32(mchNO)),
		fmt.Sprintf("%02d", GetNodeID()),
		fmt.Sprintf("%04d", random.Intn(1000))), gId
}

func GenRandNumToStr() string {
	return fmt.Sprintf("%04d", random.Intn(1000))
}

func RandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}
