import { test, expect } from '@playwright/test'

// API 测试 - 无需启动前端服务器
test.describe('文章管理 API E2E 测试', () => {
  const API_BASE_URL = 'http://localhost:8080/v1'

  let authToken: string
  const testPostId: number[] = []

  // 测试用户凭据
  const testUser = {
    username: 'admin',
    password: 'admin123'
  }

  test.beforeAll('登录获取 Token', async ({ request }) => {
    const response = await request.post(`${API_BASE_URL}/user/login`, {
      data: testUser
    })
    const data = await response.json()
    expect(data.code).toBe(200)
    authToken = data.data.token
    console.log('登录成功，获取 Token')
  })

  test.afterAll('清理测试数据', async ({ request }) => {
    // 删除创建的测试文章
    for (const id of testPostId) {
      await request.delete(`${API_BASE_URL}/post/${id}`, {
        headers: { Authorization: `Bearer ${authToken}` }
      })
    }
  })

  test('创建文章', async ({ request }) => {
    const response = await request.post(`${API_BASE_URL}/post`, {
      headers: { Authorization: `Bearer ${authToken}` },
      data: {
        title: `E2E 测试文章-${Date.now()}`,
        content: '这是 E2E 测试内容',
        summary: 'E2E 测试摘要',
        status: 1,
        categoryId: 1,
        tags: ['test', 'e2e']
      }
    })
    const data = await response.json()
    console.log('创建文章响应:', data)
    expect(response.status()).toBe(200)
    expect(data.code).toBe(200)
    expect(data.data.id).toBeTruthy()
    testPostId.push(data.data.id)
  })

  test('获取文章列表', async ({ request }) => {
    const response = await request.get(`${API_BASE_URL}/post/list`, {
      headers: { Authorization: `Bearer ${authToken}` }
    })
    const data = await response.json()
    console.log('文章列表响应:', data)
    expect(response.status()).toBe(200)
    expect(data.code).toBe(200)
    expect(Array.isArray(data.data.list)).toBe(true)
  })

  test('搜索文章', async ({ request }) => {
    const response = await request.get(`${API_BASE_URL}/posts/search`, {
      headers: { Authorization: `Bearer ${authToken}` },
      params: { keyword: '测试' }
    })
    const data = await response.json()
    console.log('搜索文章响应:', data)
    expect(response.status()).toBe(200)
    expect(data.code).toBe(200)
  })

  test('获取我的文章', async ({ request }) => {
    const response = await request.get(`${API_BASE_URL}/posts/me`, {
      headers: { Authorization: `Bearer ${authToken}` }
    })
    const data = await response.json()
    console.log('我的文章响应:', data)
    expect(response.status()).toBe(200)
    expect(data.code).toBe(200)
  })
})
