package dao

import (
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"mytask/model"
	"time"
)

// NodeQuery the filter
type NodeQuery struct {
	ID       int
	Master   string
	IsWeight bool
	Host     string
	Limit    int
	Offset   int
	Order    string
}

// NodeCondition all condition maybe used
func NodeCondition(db *gorm.DB, query *NodeQuery, total *int) *gorm.DB {
	if query.ID > 0 {
		db = db.Where("id = ?", query.ID)
	}
	if query.IsWeight == true {
		db = db.Where("weight > ?", 0)
	}
	if query.Host != "" {
		db = db.Where("host = ?", query.Host)
	}
	if query.Master != "" {
		db = db.Where("master = ?", query.Master)
	}
	if total != nil {
		db = db.Count(total)
	}
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}
	if query.Offset > 0 {
		db = db.Offset(query.Offset)
	}
	if query.Order != "" {
		db = db.Order(query.Order)
	}
	return db
}

// DescribeNodeByHost .
func DescribeNodeByHost(ip string) (node model.Node, err error) {
	if ip == "" {
		return
	}
	list, err := DescribeAllNodes(&NodeQuery{Host: ip})
	if err != nil {
		return
	}
	if len(list) > 0 {
		node = *list[0]
	}
	return
}

// DescribeMasterNode .
func DescribeMasterNode() (node model.Node, err error) {
	query := &NodeQuery{Master: model.NodeMaster.Yes}
	list, err := DescribeAllNodes(query)
	if err != nil {
		return
	}
	if len(list) > 0 {
		node = *list[0]
	}
	return
}

// GetMaxWeightNode .
func GetMaxWeightNode() (node model.Node, err error) {
	var nodes []*model.Node
	err = db.Model(&model.Node{}).Select("id,max(weight) as weight").Group("id").Find(&nodes).Error
	if err != nil {
		logs.Error("get max weight node err:%s", err.Error())
	}
	if len(nodes) == 0 {
		return node, nil
	}
	return *nodes[0], nil
}

// DescribeDisableNodes get the disable
func DescribeDisableNodes() (list []*model.Node, err error) {
	t := time.Now().Add(-time.Second * 60).Format("2006-01-02 15:04:05")
	err = db.Model(&model.Node{}).Where("weight != 0 and ping_time < ?", t).Find(&list).Error
	return
}

// DisableNode .
func DisableNode(id int) (err error) {
	err = db.Model(&model.Node{}).Where("id = ?", id).
		Update(map[string]interface{}{"master": model.NodeMaster.No, "weight": 0}).Error
	return
}

// DescribeAllNodes .
func DescribeAllNodes(query *NodeQuery) (list []*model.Node, err error) {
	err = NodeCondition(db.Model(&model.Node{}), query, nil).Find(&list).Error
	if err != nil {
		logs.Error("describe all nodes err:%s", err.Error())
	}
	return
}
