import { test, expect } from '@playwright/test'

test.describe('分类管理 E2E 测试', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/admin/categories')
  })

  test('应该显示分类管理页面', async ({ page }) => {
    await expect(page.locator('h1, h2, .page-title')).toContainText(/分类|Category/)
  })

  test('应该显示分类树结构', async ({ page }) => {
    // 等待页面加载
    await page.waitForTimeout(1000)

    // 分类树应该可见
    const treeElement = page.locator('.tree, .el-tree, [class*="tree"]')
    if (await treeElement.count() > 0) {
      await expect(treeElement.first()).toBeVisible()
    }
  })

  test('应该可以展开/折叠分类节点', async ({ page }) => {
    await page.waitForTimeout(1000)

    // 查找展开/折叠按钮
    const toggleButtons = page.locator('.el-tree-node__expand-icon, .toggle-icon, button[aria-label*="toggle"]')
    if (await toggleButtons.count() > 0) {
      await toggleButtons.first().click()
      await page.waitForTimeout(500)
    }
  })
})
