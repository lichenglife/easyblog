package store

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewStore(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	// 初始化gorm
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm %v", err)
	}

	// 关闭数据库资源
	cleanUp := func() {
		defer db.Close()
	}
	defer cleanUp()

	//执行测试代码
	_ = NewStore(gdb)

}

func Test_dataStore_User(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   UserStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dataStore{
				db: tt.fields.db,
			}
			if got := ds.User(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataStore.User() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataStore_Post(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   PostStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dataStore{
				db: tt.fields.db,
			}
			if got := ds.Post(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataStore.Post() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataStore_Close(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &dataStore{
				db: tt.fields.db,
			}
			if err := ds.Close(); (err != nil) != tt.wantErr {
				t.Errorf("dataStore.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
