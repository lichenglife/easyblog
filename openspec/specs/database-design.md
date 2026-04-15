# 博客系统数据库设计

**状态**: 草稿  
**创建时间**: 2026-04-01

---

## 设计原则

- 使用 InnoDB 引擎
- 所有表使用 `id` 作为主键（自增 BIGINT）
- 使用 `created_at` / `updated_at` 记录时间
- 使用软删除的表添加 `deleted_at` 字段
- 外键逻辑关联，不强制数据库外键约束

---

## ER 图

```
┌─────────────────────────────────────────────────────────────────────┐
│                         数据库 ER 图                                  │
└─────────────────────────────────────────────────────────────────────┘

                    ┌──────────────────┐
                    │      user        │
                    ├──────────────────┤
                    │ id               │
                    │ user_id (UUID)   │──┬────────────────────────┐
                    │ username         │  │                        │
                    │ password         │  │                        │
                    │ nick_name        │  │                        │
                    │ email            │  │                        │
                    │ phone            │  │                        │
                    │ role             │  │                        │
                    │ status           │  │                        │
                    │ github_openid    │  │                        │
                    │ created_at       │  │                        │
                    │ updated_at       │  │                        │
                    └────────┬─────────┘  │                        │
                             │           │                        │
              ┌──────────────┼───────────┼────────────────────────┤
              │              │           │                        │
              ▼              ▼           ▼                        ▼
    ┌──────────────────┐ ┌──────────────────┐            ┌──────────────────┐
    │      post        │ │     comment      │            │      like        │
    ├──────────────────┤ ├──────────────────┤            ├──────────────────┤
    │ id               │ │ id               │            │ id               │
    │ post_id (UUID)   │ │ comment_id(UUID) │            │ user_id          │
    │ user_id ─────────┼─│ user_id ─────────┼────────────│ target_type      │
    │ title            │ │ post_id          │            │ target_id        │
    │ content          │ │ parent_id        │            │ created_at       │
    │ summary          │ │ content          │            └──────────────────┘
    │ cover_image      │ │ status           │
    │ category_id ─────┼─│ created_at       │
    │ status           │ │ updated_at       │
    │ view_count       │ └────────┬─────────┘
    │ like_count       │          │
    │ created_at       │          │
    │ updated_at       │          ▼
    └────────┬─────────┘  ┌──────────────────┐
             │            │   comment_like   │
             │            └──────────────────┘
             │
             │            ┌──────────────────┐
             ├───────────▶│   post_category  │
             │            ├──────────────────┤
             │            │ id               │
             │            │ name             │
             │            │ parent_id        │
             │            │ level            │
             │            │ sort             │
             │            │ created_at       │
             │            │ updated_at       │
             │            └──────────────────┘
             │
             │            ┌──────────────────┐
             └───────────▶│    post_tag      │
                          ├──────────────────┤
                          │ id               │
                          │ name             │
                          │ created_at       │
                          │ updated_at       │
                          └────────┬─────────┘
                                   │
                          ┌────────┴─────────┐
                          ▼                  ▼
                ┌──────────────────┐  ┌──────────────────┐
                │  post_tag_map    │  │   post_category  │
                ├──────────────────┤  │       map        │
                │ post_id          │  ├──────────────────┤
                │ tag_id           │  │ post_id          │
                │ created_at       │  │ category_id      │
                └──────────────────┘  └──────────────────┘

                ┌──────────────────┐
                │   audit_log      │
                ├──────────────────┤
                │ id               │
                │ user_id          │
                │ action           │
                │ target_type      │
                │ target_id        │
                │ result           │
                │ remark           │
                │ created_at       │
                └──────────────────┘

                ┌──────────────────┐
                │  system_config   │
                ├──────────────────┤
                │ id               │
                │ config_key       │
                │ config_value     │
                │ remark           │
                │ created_at       │
                │ updated_at       │
                └──────────────────┘

                ┌──────────────────┐
                │      image       │
                ├──────────────────┤
                │ id               │
                │ image_id (UUID)  │
                │ user_id          │
                │ file_name        │
                │ file_path        │
                │ file_size        │
                │ oss_path         │
                │ oss_url          │
                │ created_at       │
                └──────────────────┘
```

---

## 表结构详细设计

### 1. user - 用户表

```sql
CREATE TABLE `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` VARCHAR(36) NOT NULL COMMENT '用户唯一 ID',
  `username` VARCHAR(50) NOT NULL COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt 加密）',
  `nick_name` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
  `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号',
  `role` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '角色：1-普通用户 2-管理员',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '状态：0-待审核 1-正常 2-封禁',
  `github_openid` VARCHAR(50) DEFAULT NULL COMMENT 'GitHub OpenID',
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '头像 URL',
  `bio` VARCHAR(500) DEFAULT NULL COMMENT '个人简介',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_github_openid` (`github_openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 2. post - 文章表

```sql
CREATE TABLE `post` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `post_id` VARCHAR(36) NOT NULL COMMENT '文章唯一 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '作者 ID',
  `title` VARCHAR(255) NOT NULL COMMENT '标题',
  `content` LONGTEXT NOT NULL COMMENT '内容（Markdown）',
  `summary` VARCHAR(500) DEFAULT NULL COMMENT '摘要（自动生成或手动设置）',
  `cover_image` VARCHAR(255) DEFAULT NULL COMMENT '封面图 URL',
  `category_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '分类 ID',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-草稿 1-已发布 2-已下架',
  `is_top` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否置顶：0-否 1-是',
  `view_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '阅读量',
  `like_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `comment_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_id` (`post_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_top` (`is_top`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';
```

### 3. comment - 评论表

```sql
CREATE TABLE `comment` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `comment_id` VARCHAR(36) NOT NULL COMMENT '评论唯一 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '评论者 ID',
  `post_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
  `parent_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '父评论 ID（楼中楼）',
  `root_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '根评论 ID（便于查询）',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '状态：0-已删除 1-正常',
  `like_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数',
  `floor` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '楼层（用于排序）',
  `level` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '层级（0-一级评论，1-二级，2-三级）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_comment_id` (`comment_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_post_id` (`post_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_root_id` (`root_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评论表';
```

### 4. post_category - 文章分类表

```sql
CREATE TABLE `post_category` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
  `parent_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '父分类 ID',
  `level` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '层级（0-一级，1-二级，2-三级）',
  `sort` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序（小到大）',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '分类描述',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_level` (`level`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类表';
```

### 5. post_tag - 文章标签表

```sql
CREATE TABLE `post_tag` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `name` VARCHAR(50) NOT NULL COMMENT '标签名称',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签表';
```

### 6. post_tag_map - 文章标签关联表

```sql
CREATE TABLE `post_tag_map` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `post_id` BIGINT UNSIGNED NOT NULL COMMENT '文章 ID',
  `tag_id` BIGINT UNSIGNED NOT NULL COMMENT '标签 ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_post_tag` (`post_id`, `tag_id`),
  KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签关联表';
```

### 7. like - 点赞表

```sql
CREATE TABLE `like` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '点赞者 ID',
  `target_type` TINYINT UNSIGNED NOT NULL COMMENT '目标类型：1-文章 2-评论',
  `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标 ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_target` (`user_id`, `target_type`, `target_id`),
  KEY `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';
```

### 8. image - 图片表

```sql
CREATE TABLE `image` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `image_id` VARCHAR(36) NOT NULL COMMENT '图片唯一 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '上传者 ID',
  `file_name` VARCHAR(255) NOT NULL COMMENT '原始文件名',
  `file_path` VARCHAR(255) NOT NULL COMMENT '本地存储路径',
  `file_size` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小（字节）',
  `mime_type` VARCHAR(50) DEFAULT NULL COMMENT '文件类型',
  `oss_path` VARCHAR(255) DEFAULT NULL COMMENT '对象存储路径',
  `oss_url` VARCHAR(500) DEFAULT NULL COMMENT '对象存储 URL',
  `sync_status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '同步状态：0-待同步 1-已同步 2-同步失败',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_image_id` (`image_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_sync_status` (`sync_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='图片表';
```

### 9. system_config - 系统配置表

```sql
CREATE TABLE `system_config` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `config_key` VARCHAR(100) NOT NULL COMMENT '配置键',
  `config_value` TEXT NOT NULL COMMENT '配置值',
  `config_type` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '配置类型：1-字符串 2-数字 3-JSON',
  `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';
```

### 10. audit_log - 审核日志表

```sql
CREATE TABLE `audit_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '操作用户 ID',
  `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
  `target_type` VARCHAR(20) DEFAULT NULL COMMENT '目标类型',
  `target_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '目标 ID',
  `result` TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '结果：1-成功 2-失败',
  `remark` VARCHAR(500) DEFAULT NULL COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审核日志表';
```

---

## 初始化数据

### 默认管理员账号

```sql
-- 初始管理员：admin / admin123
-- 密码使用 bcrypt 加密，以下为示例 hash（实际需动态生成）
INSERT INTO `user` (`user_id`, `username`, `password`, `nick_name`, `email`, `role`, `status`) 
VALUES (
  '550e8400-e29b-41d4-a716-446655440000',
  'admin',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  '系统管理员',
  'admin@easyblog.com',
  2,
  1
);
```

### 默认分类

```sql
INSERT INTO `post_category` (`name`, `parent_id`, `level`, `sort`) VALUES
('技术', NULL, 0, 1),
('生活', NULL, 0, 2),
('后端', 1, 1, 1),
('前端', 1, 1, 2),
('Go', 3, 2, 1),
('Java', 3, 2, 2),
('JavaScript', 4, 2, 1);
```

### 默认系统配置

```sql
INSERT INTO `system_config` (`config_key`, `config_value`, `config_type`, `remark`) VALUES
('site_name', 'EasyBlog', 1, '站点名称'),
('site_logo', '', 1, '站点 Logo URL'),
('site_announcement', '欢迎来到 EasyBlog！', 1, '首页公告'),
('github_client_id', '', 1, 'GitHub OAuth Client ID'),
('github_client_secret', '', 1, 'GitHub OAuth Client Secret'),
('oss_endpoint', '', 1, '对象存储 Endpoint'),
('oss_bucket', '', 1, '对象存储 Bucket'),
('oss_access_key_id', '', 1, '对象存储 AccessKey ID'),
('oss_access_key_secret', '', 1, '对象存储 AccessKey Secret');
```

---

## 索引设计说明

| 表名 | 索引名 | 索引字段 | 说明 |
|------|--------|----------|------|
| user | uk_user_id | user_id | 用户唯一标识查询 |
| user | uk_username | username | 登录查询 |
| post | idx_user_id | user_id | 用户文章列表 |
| post | idx_category_id | category_id | 分类文章列表 |
| post | idx_is_top | is_top, created_at | 置顶文章排序 |
| comment | idx_post_id | post_id | 文章评论列表 |
| comment | idx_root_id | root_id, floor | 楼中楼查询 |
| like | uk_user_target | user_id, target_type, target_id | 防重复点赞 |

---

## 数据字典

### user.role

| 值 | 说明 |
|----|------|
| 1 | 普通用户 |
| 2 | 管理员 |

### user.status

| 值 | 说明 |
|----|------|
| 0 | 待审核 |
| 1 | 正常 |
| 2 | 封禁 |

### post.status

| 值 | 说明 |
|----|------|
| 0 | 草稿 |
| 1 | 已发布 |
| 2 | 已下架 |

### like.target_type

| 值 | 说明 |
|----|------|
| 1 | 文章 |
| 2 | 评论 |
