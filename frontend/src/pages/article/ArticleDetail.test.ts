import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { useRoute, useRouter } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import ArticleDetail from './ArticleDetail.vue'
import { useUserStore } from '@/stores/user'
import type { CommentListResponse, CommentInfo } from '@/api/comment'

// Mock route params
let mockRouteParams = { id: '1' }

// Mock vue-router's useRoute and useRouter
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...actual,
    useRoute: vi.fn(() => ({
      params: mockRouteParams,
      path: '/articles/1',
      name: 'ArticleDetail',
    })),
    useRouter: vi.fn(() => ({
      push: vi.fn(),
      back: vi.fn(),
      replace: vi.fn(),
    })),
  }
})

// Mock user store
vi.mock('@/stores/user', () => ({
  useUserStore: vi.fn(() => ({
    isLoggedIn: false,
    userInfo: null,
  })),
}))

// Mock APIs
vi.mock('@/api/post', () => ({
  postApi: {
    getById: vi.fn(),
    delete: vi.fn(),
  },
}))

vi.mock('@/api/comment', () => ({
  commentApi: {
    getByPostId: vi.fn(),
    create: vi.fn(),
    delete: vi.fn(),
    like: vi.fn(),
  },
}))

vi.mock('@/api/like', () => ({
  likeApi: {
    getPostStatus: vi.fn(),
    likePost: vi.fn(),
    unlikePost: vi.fn(),
  },
}))

import { postApi } from '@/api/post'
import { commentApi } from '@/api/comment'
import { likeApi } from '@/api/like'

describe('ArticleDetail', () => {
  let wrapper: any
  let mockRouter: any

  const mockPost = {
    id: 1,
    postId: 'post-1',
    userId: 1,
    title: 'Test Article',
    content: 'Test article content',
    summary: 'Test summary',
    coverImage: '/cover.jpg',
    categoryId: 1,
    status: 1,
    isTop: 0,
    viewCount: 100,
    likeCount: 10,
    commentCount: 5,
    createAt: '2024-01-01T00:00:00Z',
    updateAt: '2024-01-01T00:00:00Z',
  }

  const mockComments: CommentInfo[] = [
    {
      id: 1,
      postId: 1,
      userId: 2,
      parentId: 0,
      content: 'Great article!',
      likeCount: 3,
      isEdited: 0,
      replyCount: 0,
      status: 1,
      createAt: '2024-01-02T00:00:00Z',
      updateAt: '2024-01-02T00:00:00Z',
    },
  ]

  const createWrapper = async () => {
    const testMockRouter = { push: vi.fn(), back: vi.fn(), replace: vi.fn() }
    vi.mocked(useRouter).mockReturnValue(testMockRouter as any)

    const w = mount(ArticleDetail, {
      global: {
        plugins: [createPinia()],
      },
    })

    await w.vm.$nextTick()

    return { wrapper: w, router: testMockRouter }
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    mockRouteParams = { id: '1' }
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: false,
      userInfo: null,
    } as any)
    vi.mocked(useRoute).mockReturnValue({
      params: mockRouteParams,
      path: '/articles/1',
      name: 'ArticleDetail',
    } as any)
    vi.mocked(postApi.getById).mockReset()
    vi.mocked(commentApi.getByPostId).mockReset()
    vi.mocked(likeApi.getPostStatus).mockReset()
  })

  it('renders loading state initially', async () => {
    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 0 })

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    expect(wrapper.vm.loading).toBe(true)
    await flushPromises()
    expect(wrapper.vm.loading).toBe(false)
  })

  it('renders article detail after data is loaded', async () => {
    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: mockComments, totalCount: 1, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    expect(wrapper.vm.post).toEqual(mockPost)
    expect(wrapper.find('.article-title').text()).toBe('Test Article')
    expect(wrapper.find('.article-content').exists()).toBe(true)
  })

  it('renders error state when article not found', async () => {
    vi.mocked(postApi.getById).mockRejectedValue(new Error('Article not found'))
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    expect(wrapper.vm.post).toBeNull()
    expect(wrapper.find('.error').exists()).toBe(true)
    expect(wrapper.text()).toContain('文章不存在或已被删除')
  })

  it('renders author actions when user is the author', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: true,
      userInfo: { userID: 1 },
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    expect(wrapper.vm.isAuthor).toBe(true)
    expect(wrapper.find('.article-actions').exists()).toBe(true)
    expect(wrapper.find('.btn-edit').exists()).toBe(true)
    expect(wrapper.find('.btn-delete').exists()).toBe(true)
  })

  it('hides author actions when user is not the author', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: true,
      userInfo: { userID: 2 },
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    expect(wrapper.vm.isAuthor).toBe(false)
    expect(wrapper.find('.article-actions').exists()).toBe(false)
  })

  it('navigates to editor page on edit button click', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: true,
      userInfo: { userID: 1 },
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    await wrapper.find('.btn-edit').trigger('click')
    expect(mockRouter.push).toHaveBeenCalledWith('/editor/1')
  })

  it('deletes article on delete button click', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: true,
      userInfo: { userID: 1 },
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })
    vi.mocked(postApi.delete).mockResolvedValue(undefined)

    // Mock confirm to return true
    const originalConfirm = window.confirm
    window.confirm = vi.fn(() => true)

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    await wrapper.find('.btn-delete').trigger('click')
    await flushPromises()

    expect(vi.mocked(postApi.delete)).toHaveBeenCalledWith(1)
    expect(mockRouter.push).toHaveBeenCalledWith('/articles')

    window.confirm = originalConfirm
  })

  it('handles like button click when logged in', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: true,
      userInfo: { userID: 1 },
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })
    vi.mocked(likeApi.likePost).mockResolvedValue(undefined)

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    const initialLikeCount = wrapper.vm.likeCount
    await wrapper.find('.meta-item-clickable').trigger('click')
    await flushPromises()

    expect(vi.mocked(likeApi.likePost)).toHaveBeenCalledWith(1)
    expect(wrapper.vm.likeCount).toBe(initialLikeCount + 1)
  })

  it('navigates to login on like when not logged in', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: false,
      userInfo: null,
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    await wrapper.find('.meta-item-clickable').trigger('click')
    expect(mockRouter.push).toHaveBeenCalledWith('/login')
  })

  it('submits comment when logged in', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: true,
      userInfo: { userID: 1 },
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })
    vi.mocked(commentApi.create).mockResolvedValue(undefined as any)

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    await wrapper.find('.comment-input').setValue('Great article!')
    await wrapper.find('.btn-submit').trigger('click')
    await flushPromises()

    expect(vi.mocked(commentApi.create)).toHaveBeenCalledWith({
      postId: 1,
      parentId: 0,
      content: 'Great article!',
    })
  })

  it('shows login hint for comments when not logged in', async () => {
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: false,
      userInfo: null,
    } as any)

    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: [], totalCount: 0, hasMore: false } as CommentListResponse)

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    expect(wrapper.find('.login-hint').exists()).toBe(true)
    expect(wrapper.text()).toContain('登录')
  })

  it('renders comments list correctly', async () => {
    vi.mocked(postApi.getById).mockResolvedValue(mockPost)
    vi.mocked(commentApi.getByPostId).mockResolvedValue({ comments: mockComments, totalCount: 1, hasMore: false } as CommentListResponse)
    vi.mocked(likeApi.getPostStatus).mockResolvedValue({ isLiked: false, count: 10 })

    const { wrapper: testWrapper, router: testRouter } = await createWrapper()
    wrapper = testWrapper
    mockRouter = testRouter

    await flushPromises()

    expect(wrapper.vm.comments.length).toBe(1)
    expect(wrapper.find('.comment-item').exists()).toBe(true)
    expect(wrapper.text()).toContain('Great article!')
  })
})
