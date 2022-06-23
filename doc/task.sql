/*
 Navicat MySQL Data Transfer

 Source Server         : localhsot
 Source Server Type    : MySQL
 Source Server Version : 80029
 Source Host           : localhost:3306
 Source Schema         : task

 Target Server Type    : MySQL
 Target Server Version : 80029
 File Encoding         : 65001

 Date: 23/06/2022 17:05:16
*/

CREATE DATABASE task;
USE task;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_node
-- ----------------------------
DROP TABLE IF EXISTS `t_node`;
CREATE TABLE `t_node` (
  `id` int NOT NULL AUTO_INCREMENT,
  `host` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `weight` int DEFAULT NULL COMMENT '选举权重值，0表示不可用',
  `master` enum('0','1') COLLATE utf8mb4_general_ci DEFAULT '0' COMMENT '是否master',
  `cpu` tinyint DEFAULT NULL COMMENT 'cpu使用率',
  `created_time` datetime DEFAULT NULL COMMENT '节点加入时间',
  `master_time` datetime DEFAULT NULL COMMENT '成为master的时间',
  `ping_time` datetime DEFAULT NULL COMMENT '最近一次心跳时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='机器节点表';

-- ----------------------------
-- Table structure for t_task_meta
-- ----------------------------
DROP TABLE IF EXISTS `t_task_meta`;
CREATE TABLE `t_task_meta` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务key',
  `task_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '任务名称',
  `task_param` json DEFAULT NULL COMMENT '任务参数',
  `task_type` tinyint DEFAULT NULL COMMENT '任务类型',
  `time_rule` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '时间规则',
  `entry_id` int DEFAULT '0' COMMENT '标识cron的id',
  `next_execute_time` datetime DEFAULT NULL COMMENT '下次执行时间',
  `task_state` tinyint DEFAULT NULL COMMENT '任务状态',
  `executed_times` int DEFAULT NULL COMMENT '执行次数',
  `task_desc` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '任务描述',
  `creator` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建者',
  `executor_serial` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '执行机器',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `task_key` (`task_key`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务表元数据';

-- ----------------------------
-- Table structure for t_task_record
-- ----------------------------
DROP TABLE IF EXISTS `t_task_record`;
CREATE TABLE `t_task_record` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务key',
  `task_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '任务名称',
  `task_param` json DEFAULT NULL COMMENT '任务参数',
  `task_progress` tinyint DEFAULT NULL COMMENT '任务进度',
  `task_result` json DEFAULT NULL COMMENT '任务结果',
  `executor_serial` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '执行机器',
  `created_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `task_key` (`task_key`)
) ENGINE=InnoDB AUTO_INCREMENT=152 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务执行记录';

SET FOREIGN_KEY_CHECKS = 1;
