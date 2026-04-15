package store

import (
	"sync"

	"gorm.io/gorm"
)

var (

	// 确保store 是实例化一次
	once sync.Once
	// 全局变量
	S *dataStore
)

// IStore 存储层工厂接口
type IStore interface {
	User() UserStore
	Post() PostStore
	Comment() CommentStore
	Like() LikeStore
	Category() CategoryStore
	Tag() TagStore
	PostTag() PostTagStore
	Close() error
}

// dataStore 实现 IStore 接口
type dataStore struct {
	db *gorm.DB
}

// dataStore 实现 IStore 接口
var _ IStore = (*dataStore)(nil)

// NewStore 创建存储层工厂
func NewStore(db *gorm.DB) IStore {
	// 确保store 是实例化一次
	once.Do(func() {
		S = &dataStore{db}
	})
	return S
}

// User() UserStore
func (ds *dataStore) User() UserStore {
	return NewUsers(ds.db)
}

// Post() PostStore

func (ds *dataStore) Post() PostStore {
	return NewPosts(ds.db)
}

// Comment() CommentStore
func (ds *dataStore) Comment() CommentStore {
	return NewComments(ds.db)
}

// Like() LikeStore
func (ds *dataStore) Like() LikeStore {
	return NewLikes(ds.db)
}

// Category() CategoryStore
func (ds *dataStore) Category() CategoryStore {
	return NewCategories(ds.db)
}

// Tag() TagStore
func (ds *dataStore) Tag() TagStore {
	return NewTags(ds.db)
}

// PostTag() PostTagStore
func (ds *dataStore) PostTag() PostTagStore {
	return NewPostTagStore(ds.db)
}

func (ds *dataStore) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// NewTestStore 创建存储层工厂实例（用于测试，不使用 once）
func NewTestStore(db *gorm.DB) IStore {
	return &dataStore{db: db}
}
