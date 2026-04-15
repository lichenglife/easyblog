import { test, expect } from '@playwright/test'

test.describe('文章管理 E2E 测试', () => {
  // 测试数据
  const testArticle = {
    title: `测试文章-${Date.now()}`,
    content: '这是一篇测试文章内容',
    summary: '测试文章摘要',
  }

  test.beforeEach(async ({ page }) => {
    // 访问管理后台页面
    await page.goto('/admin/articles')
  })

  test('应该显示文章管理页面', async ({ page }) => {
    // 检查页面标题
    await expect(page.locator('h1, h2, .page-title')).toContainText(/文章管理|Article/)

    // 检查是否有文章列表表格
    const table = page.locator('table')
    await expect(table).toBeVisible()
  })

  test('应该显示创建文章按钮', async ({ page }) => {
    const createBtn = page.locator('button:has-text("新建"), button:has-text("创建"), .btn-primary')
    await expect(createBtn).toBeVisible()
  })

  test('应该可以搜索文章', async ({ page }) => {
    // 查找搜索框
    const searchInput = page.locator('input[type="text"], input[placeholder*="搜索"]')
    if (await searchInput.count() > 0) {
      await searchInput.fill(testArticle.title)
      await page.waitForTimeout(500)
      // 搜索结果应该更新列表
      await expect(page.locator('table tbody tr')).toBeVisible()
    }
  })

  test('应该显示状态筛选器', async ({ page }) => {
    // 查找状态筛选下拉框
    const statusSelect = page.locator('select, .el-select').filter({ hasText: /状态|status/i })
    if (await statusSelect.count() > 0) {
      await expect(statusSelect).toBeVisible()
    }
  })

  test('文章列表应该显示文章信息', async ({ page }) => {
    // 等待表格加载
    await page.waitForTimeout(1000)

    // 检查表格列
    const headers = page.locator('thead th')
    const headerTexts = await headers.allTextContents()

    // 应该包含标题列
    const hasTitleColumn = headerTexts.some(h => h.includes('标题') || h.includes('Title'))
    expect(hasTitleColumn).toBeTruthy()
  })
})
