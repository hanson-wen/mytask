package model

import "time"

// Node 机器节点表
type Node struct {
	ID          int       `gorm:"primaryKey;column:id;type:int;not null" json:"id"`
	Host        string    `gorm:"column:host;type:varchar(64)" json:"host"`
	Weight      int       `gorm:"column:weight;type:int" json:"weight"`                             // 选举权重值0表示不可用
	Master      string    `gorm:"column:master;type:enum('0','1');default:null" json:"master"`      // 是否master
	CPU         int8      `gorm:"column:cpu;type:tinyint" json:"cpu"`                               // cpu使用率
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`            // 节点加入时间
	MasterTime  time.Time `gorm:"column:master_time;type:datetime;default:null" json:"master_time"` // 成为master的时间
	PingTime    time.Time `gorm:"column:ping_time;type:datetime;default:null" json:"ping_time"`     // 最近一次心跳时间
}

// TableName get sql table name.获取数据库表名
func (m *Node) TableName() string {
	return "t_node"
}

type NodeMasterDefined struct {
	No  string
	Yes string
}

var (
	NodeMaster = NodeMasterDefined{
		No:  "0",
		Yes: "1",
	}
	NodeZeroWeight = 0
)
