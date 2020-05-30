package model

import (
	"github.com/sony/sonyflake"
	"time"
)

type HasId interface {
	GetId() uint64
}

type Model struct {
	Id        uint64    `gorm:"type:bigint(20) UNSIGNED COMMENT '用户ID';PRIMARY_KEY;NOT NULL;" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime COMMENT '创建时间';NOT NULL;DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime COMMENT '最后更新时间';NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *Model) BeforeCreate() error {
	if m.Id == 0 {
		s := sonyflake.NewSonyflake(sonyflake.Settings{})
		id, err := s.NextID()
		if err != nil {
			panic(err)
		}
		m.Id = id
	}
	return nil
}

func (m *Model) GetId() uint64 {
	return m.Id
}
