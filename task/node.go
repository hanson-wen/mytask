package task

import (
	"errors"
	"mytask/dao"
	"mytask/model"
	"mytask/util"
	"time"
)

// nodeInit set node info
func nodeInit() {
	ip, err := util.GetIPV4()
	if err != nil {
		return
	}
	node, err := dao.DescribeNodeByHost(ip)
	if err != nil {
		return
	}
	if node.ID == 0 { //new node add
		node.Host = ip
		node.CreatedTime = time.Now()
		node.PingTime = time.Now()
		node.Weight = util.RandIntWithoutZero(100)
		node.CPU = util.GetCpuUsageRate()
		err = dao.CreateData(&node)
		if err != nil {
			util.Alarm("node join err:" + err.Error())
		}
	} else {
		monitorNode()
		nodePing(true)
	}
}

// nodePing report sth periodically to approve the node is active
func nodePing(r bool) {
	ip, err := util.GetIPV4()
	if err != nil {
		return
	}
	node, err := dao.DescribeNodeByHost(ip)
	if err != nil || node.ID == 0 {
		return
	}
	update := &model.Node{ID: node.ID, CPU: util.GetCpuUsageRate(), PingTime: time.Now()}
	if r == true {
		update.Weight = util.RandIntWithoutZero(100)
	}
	dao.ModifyData(update)
}

// monitorNode
// find node which is broken or offline and alarm,depend on the pingTime
func monitorNode() {
	// better use redis lock, confirm only one node run this in one moment
	// bool := redis.SetNx("monitor_key", val, 2)
	list, err := dao.DescribeDisableNodes()
	if err != nil {
		return
	}
	if len(list) == 0 {
		return
	}
	util.Alarm("some node is down")

	// set the node disable and trans the task serial
	for _, node := range list {
		_ = dao.DisableNode(node.ID)
		_ = dao.AssignTaskNoSerial(node.Host)
	}
}

// electMaster
// to select master node if master had not set yet
func electMaster() (node model.Node, err error) {
	// better use redis lock, confirm only one node run this in one moment
	// bool := redis.SetNx("elect_key", val, 2)
	node, err = dao.DescribeMasterNode()
	if err != nil {
		return
	}
	// if master had set, return
	if node.ID > 0 {
		return
	}

	// set a master
	node, err = dao.GetMaxWeightNode()

	// err or there is no node available
	if err != nil || node.ID == 0 || node.Weight == model.NodeZeroWeight {
		return node, errors.New("no available master node")
	}
	node.Master = model.NodeMaster.Yes
	node.MasterTime = time.Now()
	_ = dao.ModifyData(&node)
	return
}

// assignedTasks assign task to node
func assignedTasks() {
	ip, err := util.GetIPV4()
	if err != nil {
		return
	}
	for {
		node, err := dao.DescribeNodeByHost(ip)
		if err != nil || node.ID == 0 { //not the master
			time.Sleep(time.Second * 10)
			continue
		}

		// step1 get the un assign task by limit 10 per time
		query := &dao.TaskMetaQuery{IsNoExecutorSerial: true, Limit: 10}
		list, err := dao.DescribeAllTaskMeta(query)
		if err != nil {
			util.Alarm(err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		if len(list) == 0 { // no task , go on waiting
			time.Sleep(time.Second * 10)
			continue
		}

		// step2 get the able nodes
		nodes, err := dao.DescribeAllNodes(&dao.NodeQuery{IsWeight: true})
		if err != nil {
			util.Alarm(err.Error())
			time.Sleep(time.Second * 10)
			continue
		}
		nodesNum := len(nodes)
		if nodesNum == 0 {
			util.Alarm("no node for assigning task")
			time.Sleep(time.Second * 10)
			continue
		}

		// step3 assign task to node
		for _, v := range list {
			node := nodes[util.RandInt(nodesNum)] // simple rand,to do better,can consider cpu、mem、load...
			err = dao.ModifyData(&model.TaskMeta{ID: v.ID, ExecutorSerial: node.Host})
			if err != nil {
				util.Alarm("no node for assigning task")
				continue
			}
		}
	}
}
