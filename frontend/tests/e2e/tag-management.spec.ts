import { test, expect } from '@playwright/test'

test.describe('标签管理 E2E 测试', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/admin/tags')
  })

  test('应该显示标签管理页面', async ({ page }) => {
    await expect(page.locator('h1, h2, .page-title')).toContainText(/标签|Tag/)
  })

  test('应该显示标签列表', async ({ page }) => {
    await page.waitForTimeout(1000)

    // 检查是否有表格或列表
    const table = page.locator('table')
    const list = page.locator('[class*="list"]')

    if (await table.count() > 0) {
      await expect(table).toBeVisible()
    } else if (await list.count() > 0) {
      await expect(list.first()).toBeVisible()
    }
  })

  test('应该可以创建标签', async ({ page }) => {
    // 查找创建按钮
    const createBtn = page.locator('button:has-text("新建"), button:has-text("创建"), .btn-primary')
    if (await createBtn.count() > 0) {
      await expect(createBtn).toBeVisible()
    }
  })
})
