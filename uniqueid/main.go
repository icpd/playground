package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
	"github.com/teris-io/shortid"
)

// google/uuid
//
// 特性
// 唯一性：基于随机数和时间戳组合，生成全局唯一的 ID。
// 无序性：UUID 是随机生成的，不具有时间排序性。
// 性能：UUID 的生成速度较快，适合高并发环境。
//
// 优点
// 非常适合分布式系统中的全局唯一 ID 生成。
// 无需依赖中央节点，生成 ID 简单易用。
//
// 缺点
// 长度较长（36 个字符），不适合用于显示在 UI 上。
// 不具备时间排序特性。
//
// 适用场景
// 分布式环境中生成全局唯一 ID，例如用户 ID、订单 ID 等。
func generateUUID() {
	id := uuid.New()
	fmt.Println(id.String()) // b03e789c-62e0-4d4b-b4e8-5629422919ba
}

// bwmarrin/snowflake
//
// 特性
// 唯一性：由时间戳、机器 ID 和序列号组成，保证了全局唯一性。
// 排序性：ID 基于时间戳生成，具备时间排序特性。
// 性能：生成速度非常快，适用于高并发环境。
//
// 优点
// ID 紧凑，适合用于数据库的主键或分布式系统的唯一 ID。
// 生成的 ID 有序，方便插入有序数据库。
//
// 缺点
// 需要配置机器节点 ID，部署时稍微复杂。
//
// 适用场景
// 分布式系统中需要有序的唯一 ID，如订单号、用户 ID 等。
func generateSnowflakeID() {
	node, _ := snowflake.NewNode(1)
	id := node.Generate()
	fmt.Println(id)         // 1234567890123456789
	fmt.Println(id.Int64()) // 获取 int64 类型的 ID
}

// oklog/ulid
//
// 特性
// 唯一性：基于时间戳和随机数的组合，保证全局唯一性。
// 排序性：ID 基于时间戳生成，具备排序特性。
// 可读性：使用 26 个字符的 Base32 编码，较 UUID 更紧凑。
//
// 优点
// 具备时间有序性，方便插入有序数据库。
// 紧凑且易读，适合展示和日志记录。
//
// 缺点
// 在极高并发场景下，存在理论上的碰撞可能。
//
// 适用场景
// 需要排序且唯一的 ID，如消息队列、日志 ID 等。
func generateULID() {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	fmt.Println(id.String()) // 01FHZB1E8B9B8YZW6CVDBKMT6T
}

// segmentio/ksuid
//
// 特性
// 唯一性：基于时间戳、随机数生成，全局唯一。
// 排序性：具备时间排序特性，适合检索和排序。
// 可读性：27 字符的 Base62 字符串，适合展示和日志记录。
//
// 适用场景
// 需要唯一且具备时间排序特性的场景，如日志、消息队列等。
func generateKSUID() {
	id := ksuid.New()
	fmt.Println(id.String()) // 1CHzCEl82ZZe2r2KS1qEz3XF8Ve
}

// teris-io/shortid
//
// 特性
// 唯一性：基于随机数生成，具有较高的唯一性。
// 无序性：生成的 ID 无时间排序特性。
// 可读性：短 ID 的长度适中，便于展示。
//
// 优点
// 生成的 ID 短小，适合在 UI 和用户交互中展示。
// 不依赖中央节点，适合单机使用。
//
// 缺点
// 不具备排序特性，不能用于需要时间顺序的场景。
//
// 适用场景
// 短链接、验证码、临时文件名等场景。
func generateShortID() {
	id, _ := shortid.Generate()
	fmt.Println(id) // dppUrjK3
}

// rs/xid
//
// 特性
// 唯一性：通过时间戳、机器 ID 和计数器组合，保证全局唯一性。
// 排序性：ID 具备时间排序性，适合检索和排序。
// 性能：生成 ID 速度非常快。
//
// 适用场景
// 分布式系统中需要紧凑唯一 ID，如数据库主键、消息队列等。
func generateXID() {
	id := xid.New()
	fmt.Println(id.String()) // 9m4e2mr0ui3e8a215n4g
}

// sony/sonyflake
//
// 特性
// 唯一性：基于时间戳、机器 ID 和序列号，确保唯一性。
// 排序性：具备时间排序特性。
// 可扩展性：支持自定义机器 ID 位数和时间戳精度。
//
// 适用场景
// 高性能分布式环境，如订单号、日志序列。
func generateSonyflakeID() {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, _ := sf.NextID()
	fmt.Println(id) // 173301874793540608
}

func main() {
	generateUUID()
	generateSnowflakeID()
	generateULID()
	generateKSUID()
	generateShortID()
	generateXID()
	generateSonyflakeID()
}
