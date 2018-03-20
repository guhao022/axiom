package bot

import "axiom/utils/logger"

// robot结构体，包含所有的内部相关数据
type Robot struct {
	name   string      // 机器人名字
	logger *logger.Log // 日志
}
