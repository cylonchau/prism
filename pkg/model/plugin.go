package models

type Plugin struct {
	Id   string `gorm:"type:char(32);primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);uniqueIndex;comment:名称;not null" json:"name"`
}

func (Plugin) TableName() string {
	return "plugin"
}
