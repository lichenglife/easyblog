---
name: frontend-test-pattern
description: Vue Router composition API mocking pattern for Vitest tests using vi.mock
type: feedback
---

**Rule:** When testing Vue components that use `useRoute()` from vue-router, mock the `vue-router` module directly with `vi.mock()` at the module level, not individual router instances.

**Why:** The composition API `useRoute()` doesn't pick up mocked `$route` in global mocks. Each test needs to mock `useRoute` to return the expected route params.

**How to apply:**
```typescript
// Module level
let mockRouteParams = {}
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...actual,
    useRoute: vi.fn(() => ({ params: mockRouteParams })),
    useRouter: vi.fn(() => ({ push: vi.fn(), back: vi.fn() })),
  }
})

// In beforeEach
beforeEach(() => {
  mockRouteParams = { id: '1' } // Reset per test
  vi.mocked(useRoute).mockReturnValue({ params: mockRouteParams } as any)
})
```
