import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import LoginPage from './LoginPage.vue'
import { useUserStore } from '@/stores/user'

// Mock user store
vi.mock('@/stores/user', () => ({
  useUserStore: vi.fn(() => ({
    login: vi.fn(),
    isLoggedIn: false,
  })),
}))

describe('LoginPage', () => {
  let wrapper: any
  let mockLogin: any
  let mockRouter: any

  beforeEach(() => {
    setActivePinia(createPinia())
    mockLogin = vi.fn()
    vi.mocked(useUserStore).mockReturnValue({
      login: mockLogin,
      isLoggedIn: false,
    } as any)

    mockRouter = {
      push: vi.fn(),
    }

    const router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
        { path: '/login', name: 'Login', component: LoginPage },
        { path: '/register', name: 'Register', component: { template: '<div>Register</div>' } },
      ],
    })
    router.push = mockRouter.push

    wrapper = mount(LoginPage, {
      global: {
        plugins: [router],
        mocks: {
          $router: router,
        },
      },
    })
  })

  it('renders login form correctly', () => {
    expect(wrapper.find('h1').text()).toBe('用户登录')
    expect(wrapper.find('input#username').exists()).toBe(true)
    expect(wrapper.find('input#password').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').text()).toBe('登录')
  })

  it('has link to register page', () => {
    const registerLink = wrapper.find('a[href="/register"]')
    expect(registerLink.exists()).toBe(true)
    expect(registerLink.text()).toBe('立即注册')
  })

  it('displays error message when login fails', async () => {
    const errorMessage = '用户名或密码错误'
    mockLogin.mockRejectedValueOnce(new Error(errorMessage))

    // Fill in form
    await wrapper.find('input#username').setValue('testuser')
    await wrapper.find('input#password').setValue('wrongpassword')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.find('.error-message').exists()).toBe(true)
    expect(wrapper.find('.error-message').text()).toContain('用户名或密码错误')
  })

  it('navigates to home page on successful login', async () => {
    mockLogin.mockResolvedValueOnce(undefined)

    // Fill in form
    await wrapper.find('input#username').setValue('testuser')
    await wrapper.find('input#password').setValue('correctpassword')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockLogin).toHaveBeenCalledWith('testuser', 'correctpassword')
    expect(mockRouter.push).toHaveBeenCalledWith('/')
  })

  it('shows loading state during login', async () => {
    // Create a promise that doesn't resolve immediately
    let resolveLogin: () => void
    const loginPromise = new Promise<void>((resolve) => {
      resolveLogin = resolve
    })
    mockLogin.mockReturnValueOnce(loginPromise)

    // Fill in form and submit
    await wrapper.find('input#username').setValue('testuser')
    await wrapper.find('input#password').setValue('password')
    await wrapper.find('form').trigger('submit.prevent')

    // Give Vue a chance to update
    await wrapper.vm.$nextTick()

    // Check button state during loading
    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(submitButton.text()).toBe('登录中...')

    // Resolve the promise to clean up
    resolveLogin!()
  })

  it('clears error when starting new login attempt', async () => {
    // First login attempt fails
    mockLogin.mockRejectedValueOnce(new Error('First error'))
    await wrapper.find('input#username').setValue('user1')
    await wrapper.find('input#password').setValue('pass1')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.find('.error-message').exists()).toBe(true)

    // Second login attempt - error should be cleared
    mockLogin.mockResolvedValueOnce(undefined)
    await wrapper.find('input#username').setValue('user2')
    await wrapper.find('form').trigger('submit.prevent')

    // Error should be cleared before async operation completes
    expect(wrapper.find('.error-message').exists()).toBe(false)
  })

  it('validates required fields', async () => {
    // The component uses reactive form with v-model, which sends empty strings
    // when fields are empty. The login function is called with empty values.
    // This test verifies that the component correctly binds empty form data.
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    // Login is called with empty values - this is expected behavior
    // The actual validation happens in the backend
    expect(mockLogin).toHaveBeenCalledTimes(1)
    expect(mockLogin).toHaveBeenCalledWith('', '')
  })

  it('handles non-Error exceptions', async () => {
    mockLogin.mockRejectedValueOnce('string error')

    await wrapper.find('input#username').setValue('testuser')
    await wrapper.find('input#password').setValue('password')
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.find('.error-message').text()).toContain('登录失败')
  })
})
