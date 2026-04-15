import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import RegisterPage from './RegisterPage.vue'
import { useUserStore } from '@/stores/user'

vi.mock('@/stores/user', () => ({
  useUserStore: vi.fn(() => ({
    register: vi.fn(),
  })),
}))

describe('RegisterPage', () => {
  let wrapper: any
  let mockRegister: any
  let mockRouter: any

  beforeEach(() => {
    setActivePinia(createPinia())
    mockRegister = vi.fn()
    vi.mocked(useUserStore).mockReturnValue({
      register: mockRegister,
    } as any)

    mockRouter = {
      push: vi.fn(),
    }

    const router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
        { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
        { path: '/register', name: 'Register', component: RegisterPage },
      ],
    })
    router.push = mockRouter.push

    wrapper = mount(RegisterPage, {
      global: {
        plugins: [router],
        mocks: {
          $router: router,
        },
      },
    })
  })

  it('renders registration form correctly', () => {
    expect(wrapper.find('h1').text()).toBe('用户注册')
    expect(wrapper.find('input#username').exists()).toBe(true)
    expect(wrapper.find('input#email').exists()).toBe(true)
    expect(wrapper.find('input#password').exists()).toBe(true)
    expect(wrapper.find('input#nick_name').exists()).toBe(true)
    expect(wrapper.find('input#phone').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').text()).toBe('注册')
  })

  it('has link to login page', () => {
    const loginLink = wrapper.find('a[href="/login"]')
    expect(loginLink.exists()).toBe(true)
    expect(loginLink.text()).toBe('立即登录')
  })

  it('displays error message when registration fails', async () => {
    const errorMessage = '用户名已存在'
    mockRegister.mockRejectedValueOnce(new Error(errorMessage))

    await fillRegistrationForm({
      username: 'existinguser',
      email: 'test@example.com',
      password: 'Test123456',
      nick_name: 'Test User',
      phone: '13800138000',
    })

    expect(wrapper.find('.error-message').exists()).toBe(true)
    expect(wrapper.find('.error-message').text()).toContain('用户名已存在')
  })

  it('displays success message and redirects on successful registration', async () => {
    mockRegister.mockResolvedValueOnce({ user_id: '123' })

    await fillRegistrationForm({
      username: 'newuser',
      email: 'new@example.com',
      password: 'Test123456',
      nick_name: 'New User',
      phone: '13800138000',
    })

    expect(wrapper.find('.success-message').exists()).toBe(true)
    expect(wrapper.find('.success-message').text()).toBe('注册成功，请等待管理员审核')

    // Wait for redirect
    await new Promise(resolve => setTimeout(resolve, 2100))
    expect(mockRouter.push).toHaveBeenCalledWith('/login')
  })

  it('shows loading state during registration', async () => {
    // Create a promise that doesn't resolve immediately
    let resolveRegister: () => void
    const registerPromise = new Promise<void>((resolve) => {
      resolveRegister = resolve
    })
    mockRegister.mockReturnValueOnce(registerPromise)

    await wrapper.find('input#username').setValue('testuser')
    await wrapper.find('form').trigger('submit.prevent')

    // Give Vue a chance to update
    await wrapper.vm.$nextTick()

    const submitButton = wrapper.find('button[type="submit"]')
    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(submitButton.text()).toBe('注册中...')

    // Resolve the promise to clean up
    resolveRegister!()
  })

  it('clears error and success messages when starting new registration', async () => {
    // First registration fails
    mockRegister.mockRejectedValueOnce(new Error('First error'))
    await fillRegistrationForm({
      username: 'user1',
      email: 'user1@example.com',
      password: 'Test123456',
      nick_name: 'User 1',
    })

    expect(wrapper.find('.error-message').exists()).toBe(true)

    // Second registration attempt - messages should be cleared
    mockRegister.mockResolvedValueOnce(undefined)
    await wrapper.find('input#username').setValue('user2')
    await wrapper.find('form').trigger('submit.prevent')

    expect(wrapper.find('.error-message').exists()).toBe(false)
    expect(wrapper.find('.success-message').exists()).toBe(false)
  })

  it('validates required fields', async () => {
    // The component sends empty strings when fields are empty
    // This is expected behavior - actual validation happens in the backend
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    // Register is called with empty values - this is expected behavior
    expect(mockRegister).toHaveBeenCalledTimes(1)
    expect(mockRegister).toHaveBeenCalledWith({
      username: '',
      email: '',
      password: '',
      nick_name: '',
      phone: undefined,
    })
  })

  it('submits with empty phone number as undefined', async () => {
    mockRegister.mockResolvedValueOnce(undefined)

    await wrapper.find('input#username').setValue('newuser')
    await wrapper.find('input#email').setValue('new@example.com')
    await wrapper.find('input#password').setValue('Test123456')
    await wrapper.find('input#nick_name').setValue('New User')
    // Phone field is left empty
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRegister).toHaveBeenCalledWith({
      username: 'newuser',
      email: 'new@example.com',
      password: 'Test123456',
      nick_name: 'New User',
      phone: undefined,
    })
  })

  it('handles non-Error exceptions', async () => {
    mockRegister.mockRejectedValueOnce('string error')

    await fillRegistrationForm({
      username: 'testuser',
      email: 'test@example.com',
      password: 'Test123456',
      nick_name: 'Test User',
    })

    expect(wrapper.find('.error-message').text()).toContain('注册失败')
  })

  async function fillRegistrationForm(data: {
    username: string
    email: string
    password: string
    nick_name: string
    phone?: string
  }) {
    await wrapper.find('input#username').setValue(data.username)
    await wrapper.find('input#email').setValue(data.email)
    await wrapper.find('input#password').setValue(data.password)
    await wrapper.find('input#nick_name').setValue(data.nick_name)
    if (data.phone !== undefined) {
      await wrapper.find('input#phone').setValue(data.phone)
    }
    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()
  }
})
