package shared

import "time"

type BaseSQLModel struct {
	Id        int       `json:"-" gorm:"primaryKey;autoIncrement"`
	FakeId    *UID      `json:"id" gorm:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;autoUpdateTime"`
}

func (m *BaseSQLModel) Mask(dbType DbType) {
	uid := NewUID(uint32(m.Id), int(dbType), 1)
	m.FakeId = &uid
}
