/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 80018
Source Host           : localhost:3306
Source Database       : miniblog

Target Server Type    : MYSQL
Target Server Version : 80018
File Encoding         : 65001

Date: 2026-04-01 14:33:39
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=latin1;

-- ----------------------------
-- Table structure for post
-- ----------------------------
DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `userID` varchar(36) NOT NULL DEFAULT '' COMMENT '用户唯一 ID',
  `postID` varchar(35) NOT NULL DEFAULT '' COMMENT '博文唯一 ID',
  `title` varchar(256) NOT NULL DEFAULT '' COMMENT '博文标题',
  `content` longtext NOT NULL COMMENT '博文内容',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '博文创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '博文最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `post.postID` (`postID`),
  KEY `idx.post.userID` (`userID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='博文表';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `userID` varchar(36) NOT NULL DEFAULT '' COMMENT '用户唯一 ID',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名（唯一）',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码（加密后）',
  `nickname` varchar(30) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `email` varchar(256) NOT NULL DEFAULT '' COMMENT '用户电子邮箱地址',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '用户创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user.userID` (`userID`),
  UNIQUE KEY `user.username` (`username`),
  UNIQUE KEY `user.phone` (`phone`)
) ENGINE=MyISAM AUTO_INCREMENT=116 DEFAULT CHARSET=utf8 COMMENT='用户表';
