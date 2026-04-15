-- EasyBlog 数据库初始化脚本
-- 版本：v1.0
-- 创建时间：2026-04-01
-- 更新时间：2026-04-03

-- 设置客户端字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `easyblog` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `easyblog`;

-- ==============================
-- 数据表定义
-- ==============================

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户唯一 ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `nick_name` VARCHAR(50) DEFAULT '' COMMENT '昵称',
    `email` VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    `phone` VARCHAR(20) DEFAULT '' COMMENT '手机',
    `role` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '角色：1-普通用户 2-管理员',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态：0-待审核 1-正常 2-封禁',
    `github_openid` VARCHAR(50) DEFAULT '' COMMENT 'GitHub OpenID',
    `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像 URL',
    `bio` VARCHAR(500) DEFAULT '' COMMENT '个人简介',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_user_id` (`user_id`),
    UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 分类表
CREATE TABLE IF NOT EXISTS `categories` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `slug` VARCHAR(50) NOT NULL COMMENT '分类别名',
    `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父分类 ID',
    `level` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '层级（最多 3 级）',
    `sort` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序（越小越靠前）',
    `post_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章数',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_categories_name` (`name`),
    UNIQUE KEY `idx_categories_slug` (`slug`),
    KEY `idx_categories_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 标签表
CREATE TABLE IF NOT EXISTS `tags` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(30) NOT NULL COMMENT '标签名称',
    `slug` VARCHAR(30) NOT NULL COMMENT '标签别名',
    `post_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文章数',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-启用',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_tags_name` (`name`),
    UNIQUE KEY `idx_tags_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- 文章表
CREATE TABLE IF NOT EXISTS `posts` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '作者 ID',
    `post_id` VARCHAR(36) NOT NULL COMMENT '文章唯一 ID',
    `title` VARCHAR(255) NOT NULL COMMENT '标题',
    `content` LONGTEXT NOT NULL COMMENT '内容（Markdown）',
    `summary` VARCHAR(500) DEFAULT '' COMMENT '摘要',
    `cover_image` VARCHAR(255) DEFAULT '' COMMENT '封面图 URL',
    `category_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '分类 ID',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-草稿 1-已发布 2-已下架',
    `is_top` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否置顶：0-否 1-是',
    `view_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '阅读量',
    `like_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    `comment_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_posts_post_id` (`post_id`),
    KEY `idx_posts_user_id` (`user_id`),
    KEY `idx_posts_category_id` (`category_id`),
    KEY `idx_posts_status` (`status`),
    KEY `idx_posts_is_top` (`is_top`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- 文章标签关联表
CREATE TABLE IF NOT EXISTS `post_tags` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `post_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
    `tag_id` BIGINT UNSIGNED NOT NULL COMMENT '标签 ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_tags_unique` (`post_id`, `tag_id`),
    KEY `idx_post_tags_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章标签关联表';

-- 评论表
CREATE TABLE IF NOT EXISTS `comments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `post_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '评论者 ID',
    `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论 ID（楼中楼）',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `is_edited` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否已编辑：0-否 1-是',
    `like_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
    `reply_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '回复数',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-待审核 1-正常 2-已删除',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_comments_post_id` (`post_id`),
    KEY `idx_comments_user_id` (`user_id`),
    KEY `idx_comments_parent_id` (`parent_id`),
    KEY `idx_comments_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 点赞表
CREATE TABLE IF NOT EXISTS `likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标 ID（文章或评论）',
    `target_type` TINYINT UNSIGNED NOT NULL COMMENT '目标类型：1-文章 2-评论',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_like_unique` (`user_id`, `target_id`, `target_type`),
    KEY `idx_likes_target` (`target_id`, `target_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';

-- 图片表
CREATE TABLE IF NOT EXISTS `images` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '上传者 ID',
    `file_path` VARCHAR(255) NOT NULL COMMENT '文件路径',
    `file_name` VARCHAR(100) NOT NULL COMMENT '文件名',
    `file_size` INT UNSIGNED NOT NULL COMMENT '文件大小（字节）',
    `file_type` VARCHAR(50) NOT NULL COMMENT '文件类型（MIME）',
    `width` INT UNSIGNED DEFAULT NULL COMMENT '图片宽度',
    `height` INT UNSIGNED DEFAULT NULL COMMENT '图片高度',
    `storage` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '存储方式：1-本地 2-OSS',
    `url` VARCHAR(500) NOT NULL COMMENT '访问 URL',
    `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-禁用 1-正常',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_images_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图片表';

-- 系统配置表
CREATE TABLE IF NOT EXISTS `system_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `config_key` VARCHAR(50) NOT NULL COMMENT '配置键',
    `config_value` TEXT NOT NULL COMMENT '配置值',
    `config_type` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '配置类型：1-字符串 2-数字 3-布尔 4-JSON',
    `group_name` VARCHAR(50) DEFAULT '' COMMENT '分组名称',
    `description` VARCHAR(255) DEFAULT '' COMMENT '配置描述',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_at_by` BIGINT UNSIGNED DEFAULT NULL COMMENT '最后更新人 ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 初始化管理员账号（密码：Admin123456）
INSERT INTO `users` (`user_id`, `username`, `password`, `nick_name`, `email`, `role`, `status`) VALUES
(UUID(), 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '管理员', 'admin@easyblog.com', 2, 1)
ON DUPLICATE KEY UPDATE `username` = `username`;

-- 初始化系统配置
INSERT INTO `system_configs` (`config_key`, `config_value`, `config_type`, `group_name`, `description`) VALUES
('site_name', 'EasyBlog', 1, 'basic', '网站名称'),
('site_description', 'EasyBlog - 一个简洁的博客系统', 1, 'basic', '网站描述'),
('site_keywords', '博客，Go，Vue', 1, 'basic', '网站关键词'),
('allow_register', 'true', 3, 'user', '允许注册'),
('allow_comment', 'true', 3, 'post', '允许评论'),
('post_per_page', '10', 2, 'post', '每页文章数'),
('max_upload_size', '5242880', 2, 'upload', '最大上传大小（5MB）')
ON DUPLICATE KEY UPDATE `config_key` = `config_key`;

-- 初始化默认分类
INSERT INTO `categories` (`name`, `slug`, `level`, `sort`, `status`) VALUES
('技术', 'tech', 1, 1, 1),
('生活', 'life', 1, 2, 1),
('随笔', 'thoughts', 1, 3, 1)
ON DUPLICATE KEY UPDATE `name` = `name`;

-- 初始化默认标签
INSERT INTO `tags` (`name`, `slug`, `status`) VALUES
('Go', 'go', 1),
('Vue', 'vue', 1),
('教程', 'tutorial', 1)
ON DUPLICATE KEY UPDATE `name` = `name`;

-- ==============================
-- 初始化完成
-- ==============================
SET FOREIGN_KEY_CHECKS = 1;
COMMIT;
