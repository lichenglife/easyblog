import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import ArticleList from './ArticleList.vue'
import { useUserStore } from '@/stores/user'
import type { PostListResponse } from '@/api/post'

// Mock user store
vi.mock('@/stores/user', () => ({
  useUserStore: vi.fn(() => ({
    isLoggedIn: false,
  })),
}))

// Mock post API
vi.mock('@/api/post', () => ({
  postApi: {
    getList: vi.fn(),
  },
}))

import { postApi } from '@/api/post'

describe('ArticleList', () => {
  let wrapper: any
  let mockRouter: any
  let router: any

  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(useUserStore).mockReturnValue({
      isLoggedIn: false,
    } as any)
    vi.mocked(postApi.getList).mockReset()

    mockRouter = {
      push: vi.fn(),
    }

    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', name: 'Home', component: { template: '<div>Home</div>' } },
        { path: '/articles', name: 'Articles', component: ArticleList },
        { path: '/articles/:id', name: 'ArticleDetail', component: { template: '<div>Detail</div>' } },
        { path: '/editor', name: 'Editor', component: { template: '<div>Editor</div>' } },
      ],
    })
    router.push = mockRouter.push
  })

  it('renders loading state initially', async () => {
    // Mock empty response to avoid error in fetchPosts
    vi.mocked(postApi.getList).mockResolvedValue({ posts: [], totalCount: 0, hasMore: false } as PostListResponse)

    wrapper = mount(ArticleList, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter,
        },
      },
    })

    // Check loading state immediately after mount (before async fetchPosts completes)
    expect(wrapper.vm.loading).toBe(true)

    // Wait for onMounted to complete
    await flushPromises()
    expect(wrapper.vm.loading).toBe(false)
  })

  it('renders article list after data is loaded', async () => {
    const mockPosts = [
      {
        id: 1,
        postId: 'post-1',
        userId: 1,
        title: 'Test Post 1',
        content: 'Test content 1',
        summary: 'Test summary 1',
        coverImage: '/cover1.jpg',
        categoryId: 1,
        status: 1,
        isTop: 0,
        viewCount: 10,
        likeCount: 5,
        commentCount: 2,
        createAt: '2024-01-01T00:00:00Z',
        updateAt: '2024-01-01T00:00:00Z',
      },
    ]

    vi.mocked(postApi.getList).mockResolvedValue({
      posts: mockPosts,
      totalCount: 1,
      hasMore: false,
    } as PostListResponse)

    wrapper = mount(ArticleList, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter,
        },
      },
    })

    await flushPromises()

    expect(wrapper.vm.loading).toBe(false)
    expect(wrapper.vm.posts.length).toBe(1)
    expect(wrapper.find('.post-card').exists()).toBe(true)
  })

  it('renders empty state when no posts', async () => {
    vi.mocked(postApi.getList).mockResolvedValue({
      posts: [],
      totalCount: 0,
      hasMore: false,
    } as PostListResponse)

    wrapper = mount(ArticleList, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter,
        },
      },
    })

    await flushPromises()

    expect(wrapper.vm.loading).toBe(false)
    expect(wrapper.vm.posts.length).toBe(0)
    expect(wrapper.find('.empty').exists()).toBe(true)
  })

  it('navigates to article detail on card click', async () => {
    const mockPosts = [
      {
        id: 1,
        postId: 'post-1',
        userId: 1,
        title: 'Test Post',
        content: 'Test content',
        summary: 'Test summary',
        coverImage: '',
        categoryId: 1,
        status: 1,
        isTop: 0,
        viewCount: 0,
        likeCount: 0,
        commentCount: 0,
        createAt: '2024-01-01T00:00:00Z',
        updateAt: '2024-01-01T00:00:00Z',
      },
    ]

    vi.mocked(postApi.getList).mockResolvedValue({
      posts: mockPosts,
      totalCount: 1,
      hasMore: false,
    } as PostListResponse)

    wrapper = mount(ArticleList, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter,
        },
      },
    })

    await flushPromises()

    const postCard = wrapper.find('.post-card')
    await postCard.trigger('click')

    expect(mockRouter.push).toHaveBeenCalledWith('/articles/1')
  })

  it('filters posts by search query', async () => {
    const mockPosts = [
      {
        id: 1,
        postId: 'post-1',
        userId: 1,
        title: 'Vue Tutorial',
        content: 'Learn Vue',
        summary: '',
        coverImage: '',
        categoryId: 1,
        status: 1,
        isTop: 0,
        viewCount: 0,
        likeCount: 0,
        commentCount: 0,
        createAt: '2024-01-01T00:00:00Z',
        updateAt: '2024-01-01T00:00:00Z',
      },
      {
        id: 2,
        postId: 'post-2',
        userId: 1,
        title: 'React Tutorial',
        content: 'Learn React',
        summary: '',
        coverImage: '',
        categoryId: 1,
        status: 1,
        isTop: 0,
        viewCount: 0,
        likeCount: 0,
        commentCount: 0,
        createAt: '2024-01-01T00:00:00Z',
        updateAt: '2024-01-01T00:00:00Z',
      },
    ]

    vi.mocked(postApi.getList).mockResolvedValue({
      posts: mockPosts,
      totalCount: 2,
      hasMore: false,
    } as PostListResponse)

    wrapper = mount(ArticleList, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter,
        },
      },
    })

    await flushPromises()

    // Set search query
    wrapper.vm.searchQuery = 'Vue'
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.filteredPosts.length).toBe(1)
    expect(wrapper.vm.filteredPosts[0].title).toBe('Vue Tutorial')
  })

  it('handles page navigation', async () => {
    vi.mocked(postApi.getList).mockResolvedValue({
      posts: [],
      totalCount: 0,
      hasMore: false,
    } as PostListResponse)

    wrapper = mount(ArticleList, {
      global: {
        plugins: [router],
        mocks: {
          $router: mockRouter,
        },
      },
    })

    await flushPromises()

    // Initial page should be 1
    expect(wrapper.vm.currentPage).toBe(1)

    // Call handlePageChange
    wrapper.vm.handlePageChange(2)
    expect(wrapper.vm.currentPage).toBe(1) // Should stay at 1 since there's only one page
  })
})
