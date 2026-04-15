import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'

// Mock user API
vi.mock('@/api/user', () => ({
  userApi: {
    login: vi.fn(),
    register: vi.fn(),
    getCurrentUser: vi.fn(),
    updateUserInfo: vi.fn(),
    changePassword: vi.fn(),
    logout: vi.fn(),
    getUserList: vi.fn(),
    getUserById: vi.fn(),
    deleteUser: vi.fn(),
    auditUser: vi.fn(),
    banUser: vi.fn(),
  },
}))

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}

  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      store = {}
    }),
  }
})()

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
})

describe('useUserStore', () => {
  let store: ReturnType<typeof useUserStore>

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useUserStore()
    vi.clearAllMocks()
  })

  afterEach(() => {
    localStorageMock.clear()
  })

  describe('initial state', () => {
    it('should have empty token and null userInfo', () => {
      expect(store.token).toBe('')
      expect(store.userInfo).toBeNull()
    })

    it('should have correct computed properties', () => {
      expect(store.isLoggedIn).toBe(false)
      expect(store.isAdmin).toBe(false)
    })
  })

  describe('initFromStorage', () => {
    it('should restore token and userInfo from localStorage', () => {
      const mockToken = 'test-token-123'
      const mockUserInfo = {
        userID: '1',
        username: 'testuser',
        nickname: 'Test User',
        email: 'test@example.com',
        role: 1,
        status: 1,
        createAt: '2024-01-01',
      }

      localStorageMock.getItem.mockImplementation((key: string) => {
        if (key === 'token') return mockToken
        if (key === 'userInfo') return JSON.stringify(mockUserInfo)
        return null
      })

      store.initFromStorage()

      expect(store.token).toBe(mockToken)
      expect(store.userInfo).toEqual(mockUserInfo)
      expect(localStorageMock.getItem).toHaveBeenCalledWith('token')
      expect(localStorageMock.getItem).toHaveBeenCalledWith('userInfo')
    })

    it('should handle invalid JSON in userInfo', () => {
      localStorageMock.getItem.mockImplementation((key: string) => {
        if (key === 'token') return 'test-token'
        if (key === 'userInfo') return 'invalid-json'
        return null
      })

      store.initFromStorage()

      expect(store.userInfo).toBeNull()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('userInfo')
    })
  })

  describe('setToken', () => {
    it('should set token and save to localStorage', () => {
      const newToken = 'new-token-456'

      store.setToken(newToken)

      expect(store.token).toBe(newToken)
      expect(localStorageMock.setItem).toHaveBeenCalledWith('token', newToken)
    })
  })

  describe('setUserInfo', () => {
    it('should set userInfo', () => {
      const userInfo = {
        userID: '1',
        username: 'testuser',
        nickname: 'Test User',
        email: 'test@example.com',
        role: 1,
        status: 1,
        createAt: '2024-01-01',
      }

      store.setUserInfo(userInfo)

      expect(store.userInfo).toEqual(userInfo)
    })
  })

  describe('login', () => {
    it('should login successfully and set token and userInfo', async () => {
      const mockResponse = {
        token: 'login-token',
        user: {
          userID: '1',
          username: 'testuser',
          nickname: 'Test User',
          email: 'test@example.com',
          role: 1,
          status: 1,
          createAt: '2024-01-01',
        },
      }

      vi.mocked(userApi.login).mockResolvedValue(mockResponse)

      await store.login('testuser', 'password123')

      expect(userApi.login).toHaveBeenCalledWith({ username: 'testuser', password: 'password123' })
      expect(store.token).toBe('login-token')
      expect(store.userInfo).toEqual(mockResponse.user)
      expect(localStorageMock.setItem).toHaveBeenCalledWith('token', 'login-token')
      expect(localStorageMock.setItem).toHaveBeenCalledWith('userInfo', JSON.stringify(mockResponse.user))
    })

    it('should throw error on login failure', async () => {
      const error = new Error('Invalid credentials')
      vi.mocked(userApi.login).mockRejectedValue(error)

      await expect(store.login('testuser', 'wrongpassword')).rejects.toThrow('Invalid credentials')
      expect(store.token).toBe('')
      expect(store.userInfo).toBeNull()
    })
  })

  describe('register', () => {
    it('should register successfully', async () => {
      const mockResponse = { user_id: '123' }
      vi.mocked(userApi.register).mockResolvedValue(mockResponse)

      const result = await store.register({
        username: 'newuser',
        password: 'NewPassword123',
        email: 'new@example.com',
        nick_name: 'New User',
        phone: '13800138000',
      })

      expect(userApi.register).toHaveBeenCalledWith({
        username: 'newuser',
        password: 'NewPassword123',
        email: 'new@example.com',
        nick_name: 'New User',
        phone: '13800138000',
      })
      expect(result).toEqual(mockResponse)
    })

    it('should throw error on registration failure', async () => {
      const error = new Error('Username already exists')
      vi.mocked(userApi.register).mockRejectedValue(error)

      await expect(
        store.register({
          username: 'existinguser',
          password: 'Password123',
          email: 'existing@example.com',
          nick_name: 'Existing User',
        })
      ).rejects.toThrow('Username already exists')
    })
  })

  describe('fetchCurrentUser', () => {
    it('should fetch current user successfully', async () => {
      const mockUser = {
        userID: '1',
        username: 'testuser',
        nickname: 'Test User',
        email: 'test@example.com',
        role: 1,
        status: 1,
        createAt: '2024-01-01',
      }

      vi.mocked(userApi.getCurrentUser).mockResolvedValue(mockUser)

      const result = await store.fetchCurrentUser()

      expect(userApi.getCurrentUser).toHaveBeenCalled()
      expect(result).toEqual(mockUser)
      expect(store.userInfo).toEqual(mockUser)
      expect(localStorageMock.setItem).toHaveBeenCalledWith('userInfo', JSON.stringify(mockUser))
    })

    it('should logout and throw error on fetch failure', async () => {
      const error = new Error('Token expired')
      vi.mocked(userApi.getCurrentUser).mockRejectedValue(error)

      await expect(store.fetchCurrentUser()).rejects.toThrow('Token expired')
      expect(store.token).toBe('')
      expect(store.userInfo).toBeNull()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('userInfo')
    })
  })

  describe('logout', () => {
    it('should clear token and userInfo', () => {
      store.token = 'existing-token'
      store.userInfo = {
        userID: '1',
        username: 'testuser',
        nickname: 'Test User',
        email: 'test@example.com',
        role: 1,
        status: 1,
        createAt: '2024-01-01',
      }

      store.logout()

      expect(store.token).toBe('')
      expect(store.userInfo).toBeNull()
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('userInfo')
    })
  })

  describe('computed properties', () => {
    it('should return isLoggedIn true when token exists', () => {
      store.token = 'some-token'
      expect(store.isLoggedIn).toBe(true)
    })

    it('should return isLoggedIn false when token is empty', () => {
      store.token = ''
      expect(store.isLoggedIn).toBe(false)
    })

    it('should return isAdmin true when role is 2', () => {
      store.userInfo = {
        userID: '1',
        username: 'admin',
        nickname: 'Admin User',
        email: 'admin@example.com',
        role: 2,
        status: 1,
        createAt: '2024-01-01',
      }
      expect(store.isAdmin).toBe(true)
    })

    it('should return isAdmin false when role is not 2', () => {
      store.userInfo = {
        userID: '1',
        username: 'user',
        nickname: 'Regular User',
        email: 'user@example.com',
        role: 1,
        status: 1,
        createAt: '2024-01-01',
      }
      expect(store.isAdmin).toBe(false)
    })

    it('should return isAdmin false when userInfo is null', () => {
      store.userInfo = null
      expect(store.isAdmin).toBe(false)
    })
  })
})
