import { test, expect } from '@playwright/test'

test.describe('用户登录流程', () => {
  test('应该显示登录页面', async ({ page }) => {
    await page.goto('/login')
    await expect(page).toHaveTitle(/EasyBlog/)
    await expect(page.locator('input[name="username"]')).toBeVisible()
    await expect(page.locator('input[name="password"]')).toBeVisible()
  })

  test('用户名为空应该显示错误', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[name="password"]', 'Test123456')
    await page.click('button[type="submit"]')
    // 等待错误消息
    await page.waitForSelector('.el-message--error', { timeout: 5000 })
    await expect(page.locator('.el-message--error')).toBeVisible()
  })

  test('密码为空应该显示错误', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[name="username"]', 'testuser')
    await page.click('button[type="submit"]')
    await page.waitForSelector('.el-message--error', { timeout: 5000 })
    await expect(page.locator('.el-message--error')).toBeVisible()
  })
})
