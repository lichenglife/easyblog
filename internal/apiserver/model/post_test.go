package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lichenglife/easyblog/internal/apiserver/model"
)

func TestPost_TableName(t *testing.T) {
	expected := "post"
	post := model.Post{}
	// Ensure that the TableName method returns the expected value
	if post.TableName() != expected {
		t.Errorf("expected TableName %q, got %q", expected, post.TableName())
	}
}

func TestPost_JSONMarshalling(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	p := &model.Post{
		ID:        1,
		UserID:    "user-123",
		PostID:    "post-456",
		Content:   "Test content",
		Title:     "Test title",
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err = json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if result["id"] != float64(1) { // json package decodes numbers as float64
		t.Errorf("expected id %v, got %v", 1, result["id"])
	}
	if result["userID"] != p.UserID {
		t.Errorf("expected userID %q, got %q", p.UserID, result["userID"])
	}
	if result["postID"] != p.PostID {
		t.Errorf("expected postID %q, got %q", p.PostID, result["postID"])
	}
	if result["content"] != p.Content {
		t.Errorf("expected content %q, got %q", p.Content, result["content"])
	}
	if result["title"] != p.Title {
		t.Errorf("expected title %q, got %q", p.Title, result["title"])
	}
	if result["createAt"] == "" {
		t.Error("expected createAt to be present")
	}
	if result["updateAt"] == "" {
		t.Error("expected updateAt to be present")
	}
}

func TestCreatePostRequest_JSONUnmarshal(t *testing.T) {
	jsonData := `{"userID": "user-789", "content": "New post content", "title": "New post title"}`
	var req model.CreatePostRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if req.UserID != "user-789" {
		t.Errorf("expected userID %q, got %q", "user-789", req.UserID)
	}
	if req.Content != "New post content" {
		t.Errorf("expected content %q, got %q", "New post content", req.Content)
	}
	if req.Title != "New post title" {
		t.Errorf("expected title %q, got %q", "New post title", req.Title)
	}
}

func TestUpdatePostRequest_JSONUnmarshal(t *testing.T) {
	jsonData := `{"postID": "post-001", "content": "Updated content", "title": "Updated title"}`
	var req model.UpdatePostRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if req.PostID != "post-001" {
		t.Errorf("expected postID %q, got %q", "post-001", req.PostID)
	}
	if req.Content != "Updated content" {
		t.Errorf("expected content %q, got %q", "Updated content", req.Content)
	}
	if req.Title != "Updated title" {
		t.Errorf("expected title %q, got %q", "Updated title", req.Title)
	}
}

func TestGetPostRequest_JSONUnmarshal(t *testing.T) {
	jsonData := `{"id": 10}`
	var req model.GetPostRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if req.ID != 10 {
		t.Errorf("expected id %d, got %d", 10, req.ID)
	}
}

func TestPageListRequest_JSONUnmarshal(t *testing.T) {
	jsonData := `{"page": 2, "pageSize": 5}`
	var req model.PageListRequest
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if req.Page != 2 {
		t.Errorf("expected page %d, got %d", 2, req.Page)
	}
	if req.PageSize != 5 {
		t.Errorf("expected pageSize %d, got %d", 5, req.PageSize)
	}
}

func TestListPostResponse_JSONMarshalling(t *testing.T) {
	posts := []*model.Post{
		{
			ID:        1,
			UserID:    "user-1",
			PostID:    "post-1",
			Content:   "Content 1",
			Title:     "Title 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    "user-2",
			PostID:    "post-2",
			Content:   "Content 2",
			Title:     "Title 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	resp := model.ListPostResponse{
		TotalCount: 2,
		HasMore:    false,
		Posts:      posts,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err = json.Unmarshal(data, &result); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if int(result["totalCount"].(float64)) != 2 {
		t.Errorf("expected totalCount %d, got %v", 2, result["totalCount"])
	}

	if result["hasMore"] != false {
		t.Errorf("expected hasMore to be false, got %v", result["hasMore"])
	}

	postsArr, ok := result["posts"].([]interface{})
	if !ok || len(postsArr) != 2 {
		t.Errorf("expected posts array of length 2, got %v", result["posts"])
	}
}
