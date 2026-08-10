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

	"go-admin/server/utils/genid"
)

// GenDeviceID 向后兼容别名
var GenDeviceID = genid.GenDeviceID

// ParseID 解析ID各部分信息
func ParseID(id int64) map[string]int64 {
	return map[string]int64{
		"timestamp_diff": (id >> 17) & ((1 << 46) - 1),
		"random":         (id >> 10) & ((1 << 7) - 1),
		"sequence":       (id >> 1) & ((1 << 9) - 1),
		"rollback_flag":  id & 1,
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
