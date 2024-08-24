/*
 Navicat Premium Data Transfer

 Source Server         : MYSQL-3307
 Source Server Type    : MySQL
 Source Server Version : 50744 (5.7.44)
 Source Host           : 43.139.223.7:3307
 Source Schema         : goadmin

 Target Server Type    : MySQL
 Target Server Version : 50744 (5.7.44)
 File Encoding         : 65001

 Date: 21/08/2024 19:25:18
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_auth_assignment
-- ----------------------------
DROP TABLE IF EXISTS `sys_auth_assignment`;
CREATE TABLE `sys_auth_assignment` (
                                       `item_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
                                       `user_id` bigint(20) NOT NULL,
                                       `create_time` datetime DEFAULT NULL,
                                       PRIMARY KEY (`item_name`,`user_id`) USING BTREE,
                                       CONSTRAINT `sys_auth_assignment_ibfk_1` FOREIGN KEY (`item_name`) REFERENCES `sys_auth_item` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for sys_auth_item
-- ----------------------------
DROP TABLE IF EXISTS `sys_auth_item`;
CREATE TABLE `sys_auth_item` (
                                 `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
                                 `type` int(32) NOT NULL,
                                 `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `rule_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                                 `data` blob,
                                 `create_time` datetime DEFAULT NULL,
                                 `update_time` datetime DEFAULT NULL,
                                 PRIMARY KEY (`name`) USING BTREE,
                                 KEY `type` (`type`) USING BTREE,
                                 KEY `rule_name` (`rule_name`) USING BTREE,
                                 CONSTRAINT `sys_auth_item_ibfk_1` FOREIGN KEY (`rule_name`) REFERENCES `sys_auth_rule` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for sys_auth_item_child
-- ----------------------------
DROP TABLE IF EXISTS `sys_auth_item_child`;
CREATE TABLE `sys_auth_item_child` (
                                       `parent` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
                                       `child` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
                                       PRIMARY KEY (`parent`,`child`) USING BTREE,
                                       KEY `child` (`child`) USING BTREE,
                                       CONSTRAINT `sys_auth_item_child_ibfk_1` FOREIGN KEY (`parent`) REFERENCES `sys_auth_item` (`name`),
                                       CONSTRAINT `sys_auth_item_child_ibfk_2` FOREIGN KEY (`child`) REFERENCES `sys_auth_item` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for sys_auth_rule
-- ----------------------------
DROP TABLE IF EXISTS `sys_auth_rule`;
CREATE TABLE `sys_auth_rule` (
                                 `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
                                 `data` blob,
                                 `create_time` datetime DEFAULT NULL,
                                 `update_time` datetime DEFAULT NULL,
                                 PRIMARY KEY (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
