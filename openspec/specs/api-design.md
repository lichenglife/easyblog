# 博客系统 API 接口设计

**状态**: 草稿  
**创建时间**: 2026-04-01  
**版本**: v1

---

## 接口规范

### 基础 URL

```
/api/v1
```

### 认证方式

```
Authorization: Bearer <jwt_token>
```

### 响应格式

**成功响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**错误响应**
```json
{
  "code": 10001,
  "message": "错误描述",
  "data": null
}
```

### 分页格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "has_more": true
  }
}
```

---

## 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10000 | 系统错误 |
| 10001 | 参数校验失败 |
| 10002 | 未授权 |
| 10003 | 权限不足 |
| 10004 | 资源不存在 |
| 10005 | 资源已存在 |

---

## 接口清单

### 用户模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 注册 | POST | /user/register | ❌ | 用户注册 |
| 登录 | POST | /user/login | ❌ | 用户名密码登录 |
| GitHub 登录 | GET | /user/github | ❌ | GitHub OAuth 入口 |
| GitHub 回调 | GET | /user/github/callback | ❌ | GitHub OAuth 回调 |
| 绑定 GitHub | POST | /user/bind-github | ✅ | 绑定 GitHub 账号 |
| 登出 | POST | /user/logout | ✅ | 用户登出 |
| 获取当前用户 | GET | /user/me | ✅ | 获取当前用户信息 |
| 更新用户信息 | PUT | /user/me | ✅ | 更新个人信息 |
| 修改密码 | PUT | /user/password | ✅ | 修改密码 |
| 获取用户列表 | GET | /users | ✅(管理员) | 用户管理列表 |
| 审核用户 | PUT | /users/:id/audit | ✅(管理员) | 用户审核 |
| 封禁用户 | PUT | /users/:id/ban | ✅(管理员) | 封禁用户 |
| 删除用户 | DELETE | /users/:id | ✅(管理员) | 删除用户 |

### 文章模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 创建文章 | POST | /posts | ✅ | 创建文章 |
| 文章列表 | GET | /posts | ❌ | 获取文章列表 |
| 文章详情 | GET | /posts/:id | ❌ | 获取文章详情 |
| 更新文章 | PUT | /posts/:id | ✅ | 更新文章 |
| 删除文章 | DELETE | /posts/:id | ✅ | 删除文章 |
| 置顶文章 | PUT | /posts/:id/top | ✅(管理员) | 置顶/取消置顶 |
| 我的文章 | GET | /posts/mine | ✅ | 获取我的文章 |

### 分类模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 分类列表 | GET | /categories | ❌ | 获取分类列表 |
| 创建分类 | POST | /categories | ✅(管理员) | 创建分类 |
| 更新分类 | PUT | /categories/:id | ✅(管理员) | 更新分类 |
| 删除分类 | DELETE | /categories/:id | ✅(管理员) | 删除分类 |
| 分类树 | GET | /categories/tree | ❌ | 获取分类树 |

### 标签模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 标签列表 | GET | /tags | ❌ | 获取标签列表 |
| 创建标签 | POST | /tags | ✅ | 创建标签 |
| 删除标签 | DELETE | /tags/:id | ✅(管理员) | 删除标签 |

### 评论模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 评论列表 | GET | /comments | ❌ | 获取评论列表 |
| 创建评论 | POST | /comments | ✅ | 发表评论 |
| 更新评论 | PUT | /comments/:id | ✅ | 更新评论 |
| 删除评论 | DELETE | /comments/:id | ✅ | 删除评论 |
| 点赞评论 | POST | /comments/:id/like | ✅ | 点赞评论 |

### 点赞模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 点赞文章 | POST | /posts/:id/like | ✅ | 点赞文章 |
| 取消点赞 | DELETE | /posts/:id/like | ✅ | 取消点赞 |

### 图片模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 上传图片 | POST | /images | ✅ | 上传图片 |
| 图片列表 | GET | /images | ✅ | 我的图片列表 |
| 删除图片 | DELETE | /images/:id | ✅ | 删除图片 |

### 系统模块

| 接口 | 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|------|
| 获取配置 | GET | /config | ❌ | 获取公开配置 |
| 更新配置 | PUT | /config | ✅(管理员) | 更新系统配置 |

---

## 接口详细设计

### 用户模块

#### 1. 用户注册

```
POST /api/v1/user/register
```

**请求**
```json
{
  "username": "zhangsan",
  "password": "123456",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "nick_name": "张三"
}
```

**响应**
```json
{
  "code": 0,
  "message": "注册成功，请等待审核",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440001"
  }
}
```

---

#### 2. 用户登录

```
POST /api/v1/user/login
```

**请求**
```json
{
  "username": "zhangsan",
  "password": "123456"
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expire_in": 7200,
    "user": {
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "username": "zhangsan",
      "nick_name": "张三",
      "email": "zhangsan@example.com",
      "role": 1
    }
  }
}
```

---

#### 3. 获取当前用户信息

```
GET /api/v1/user/me
Authorization: Bearer <token>
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "username": "zhangsan",
    "nick_name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "role": 1,
    "status": 1,
    "avatar": "https://...",
    "created_at": "2026-04-01T10:00:00Z"
  }
}
```

---

#### 4. 更新用户信息

```
PUT /api/v1/user/me
Authorization: Bearer <token>
```

**请求**
```json
{
  "nick_name": "新昵称",
  "email": "newemail@example.com",
  "phone": "13900139000",
  "avatar": "https://...",
  "bio": "个人简介"
}
```

---

#### 5. 修改密码

```
PUT /api/v1/user/password
Authorization: Bearer <token>
```

**请求**
```json
{
  "old_password": "123456",
  "new_password": "newpassword123"
}
```

---

#### 6. 用户列表（管理员）

```
GET /api/v1/users?page=1&page_size=10&status=1&keyword=xxx
Authorization: Bearer <token>
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "user_id": "550e8400-e29b-41d4-a716-446655440001",
        "username": "zhangsan",
        "nick_name": "张三",
        "email": "zhangsan@example.com",
        "role": 1,
        "status": 1,
        "created_at": "2026-04-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "has_more": true
  }
}
```

---

#### 7. 审核用户（管理员）

```
PUT /api/v1/users/:id/audit
Authorization: Bearer <token>
```

**请求**
```json
{
  "approved": true,
  "reason": ""
}
```

---

### 文章模块

#### 1. 创建文章

```
POST /api/v1/posts
Authorization: Bearer <token>
```

**请求**
```json
{
  "title": "我的第一篇博客",
  "content": "# Hello World\n\n这是正文内容...",
  "category_id": 1,
  "tag_ids": [1, 2, 3],
  "cover_image": "https://..."
}
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "post_id": "660e8400-e29b-41d4-a716-446655440001",
    "title": "我的第一篇博客",
    "created_at": "2026-04-01T10:00:00Z"
  }
}
```

---

#### 2. 文章列表

```
GET /api/v1/posts?page=1&page_size=10&category_id=1&tag_id=2&keyword=xxx
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "post_id": "660e8400-e29b-41d4-a716-446655440001",
        "title": "我的第一篇博客",
        "summary": "这是文章摘要...",
        "cover_image": "https://...",
        "author": {
          "user_id": "550e8400-e29b-41d4-a716-446655440001",
          "nick_name": "张三"
        },
        "category": {
          "id": 1,
          "name": "技术"
        },
        "tags": [
          {"id": 1, "name": "Go"},
          {"id": 2, "name": "后端"}
        ],
        "view_count": 100,
        "like_count": 10,
        "is_top": 1,
        "created_at": "2026-04-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10,
    "has_more": true
  }
}
```

---

#### 3. 文章详情

```
GET /api/v1/posts/:id
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "post_id": "660e8400-e29b-41d4-a716-446655440001",
    "title": "我的第一篇博客",
    "content": "# Hello World\n\n这是正文内容...",
    "cover_image": "https://...",
    "author": {
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "nick_name": "张三",
      "avatar": "https://..."
    },
    "category": {
      "id": 1,
      "name": "技术"
    },
    "tags": [
      {"id": 1, "name": "Go"},
      {"id": 2, "name": "后端"}
    ],
    "view_count": 101,
    "like_count": 10,
    "is_liked": false,
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-01T11:00:00Z"
  }
}
```

---

#### 4. 更新文章

```
PUT /api/v1/posts/:id
Authorization: Bearer <token>
```

**请求**
```json
{
  "title": "更新后的标题",
  "content": "更新后的内容...",
  "category_id": 2,
  "tag_ids": [3, 4],
  "cover_image": "https://..."
}
```

---

#### 5. 删除文章

```
DELETE /api/v1/posts/:id
Authorization: Bearer <token>
```

---

#### 6. 置顶文章（管理员）

```
PUT /api/v1/posts/:id/top
Authorization: Bearer <token>
```

**请求**
```json
{
  "is_top": true
}
```

---

### 分类模块

#### 1. 分类树

```
GET /api/v1/categories/tree
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "技术",
      "parent_id": null,
      "level": 0,
      "children": [
        {
          "id": 3,
          "name": "后端",
          "parent_id": 1,
          "level": 1,
          "children": [
            {"id": 5, "name": "Go", "parent_id": 3, "level": 2, "children": []}
          ]
        },
        {
          "id": 4,
          "name": "前端",
          "parent_id": 1,
          "level": 1,
          "children": []
        }
      ]
    },
    {
      "id": 2,
      "name": "生活",
      "parent_id": null,
      "level": 0,
      "children": []
    }
  ]
}
```

---

### 评论模块

#### 1. 评论列表

```
GET /api/v1/comments?post_id=1&page=1&page_size=20
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "comment_id": "770e8400-e29b-41d4-a716-446655440001",
        "user": {
          "user_id": "550e8400-e29b-41d4-a716-446655440001",
          "nick_name": "张三",
          "avatar": "https://..."
        },
        "content": "写得真好！",
        "like_count": 5,
        "is_liked": false,
        "level": 0,
        "reply_count": 2,
        "created_at": "2026-04-01T10:00:00Z",
        "replies": [
          {
            "comment_id": "770e8400-e29b-41d4-a716-446655440002",
            "user": {
              "user_id": "550e8400-e29b-41d4-a716-446655440002",
              "nick_name": "李四"
            },
            "content": "同意！",
            "created_at": "2026-04-01T11:00:00Z"
          }
        ]
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20,
    "has_more": true
  }
}
```

---

#### 2. 创建评论

```
POST /api/v1/comments
Authorization: Bearer <token>
```

**请求**
```json
{
  "post_id": 1,
  "content": "写得真好！",
  "parent_id": null
}
```

---

#### 3. 删除评论

```
DELETE /api/v1/comments/:id
Authorization: Bearer <token>
```

---

### 图片模块

#### 1. 上传图片

```
POST /api/v1/images
Content-Type: multipart/form-data
Authorization: Bearer <token>
```

**请求**
```
file: <image_file>
```

**响应**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "image_id": "880e8400-e29b-41d4-a716-446655440001",
    "url": "https://easyblog.com/images/xxx.jpg",
    "local_path": "/data/images/2026/04/01/xxx.jpg",
    "file_size": 102400
  }
}
```

---
