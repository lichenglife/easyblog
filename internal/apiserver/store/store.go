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
	// Post() PostStore
	User() UserStore

	Post() PostStore

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

func (ds *dataStore) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
