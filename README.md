## 分布式作业调度系统

### 特性
- 选举模式，减少分布式部署的机器任务争夺。（see task.electMaster）
- 异常节点监控，踢除异常节点，任务自动迁移。(see task.monitorNode)
- 负载监控，自动告警并降低任务执行量。
- 任务重试。

### 可迭代优化点
- 选择算法
- 作业分配，结束cpu、内存等机器负载。
- 告警完善


### 如何执行？
- 配置见conf/app.conf
- mysql 表数据见doc/
- 最后go run main.go

### 相关api
- 查看postman发布的链接：https://documenter.getpostman.com/view/1280333/UzBpMSH1
