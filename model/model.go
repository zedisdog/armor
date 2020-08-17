package model

import (
	"github.com/jinzhu/gorm"
	"github.com/sony/sonyflake"
	"sync"
	"time"
)

var once sync.Once
var snow *sonyflake.Sonyflake

func GetSnow() *sonyflake.Sonyflake {
	once.Do(func() {
		snow = sonyflake.NewSonyflake(sonyflake.Settings{})
	})
	return snow
}

type HasId interface {
	GetId() uint64
}

type Model struct {
	ID        uint64    `gorm:"type:bigint(20) UNSIGNED COMMENT '用户ID';PRIMARY_KEY;NOT NULL;" json:"id,string"`
	CreatedAt time.Time `gorm:"type:datetime COMMENT '创建时间' DEFAULT CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime COMMENT '最后更新时间' DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;NOT NULL;" json:"updated_at"`
}

func (m *Model) GetId() uint64 {
	return m.ID
}

func (m *Model) BeforeSave(scope *gorm.Scope) error {
	if m.ID == 0 {
		id, err := GetSnow().NextID()
		if err != nil {
			panic(err)
		}
		scope.SetColumn("ID", id)
	}
	return nil
}
