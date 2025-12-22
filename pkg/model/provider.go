package models

type Provider struct {
	ID          int64  `gorm:"type:bigint;primaryKey;autoIncrement:false" json:"id"`
	Name        string `gorm:"type:varchar(100);uniqueIndex;comment:名称;not null;index:idx_name_version" json:"name"`
	Registry    string `gorm:"type:varchar(255);comment:provider 供应商;not null" json:"provider"`
	Version     string `gorm:"type:varchar(64);comment:provider 版本号;not null;index:idx_name_version" json:"version"`
	Namespace   string `gorm:"type:varchar(32);not null;default:'';comment:terraform namespace" json:"namespace" xorm:"namespace"`
	Initialized bool   `gorm:"not null;default:false;comment:是否初始化" json:"initialized" xorm:"initialized"`
	Enabled     bool   `gorm:"not null;default:true;comment:是否启用" json:"enabled"`
	IsDefault   bool   `gorm:"not null;default:false;comment:是否为默认版本" json:"is_default"`
	Description string `gorm:"type:varchar(512);default:''" json:"description"`
}

func (Provider) TableName() string {
	return "provider"
}

func (Provider) Indexes() map[string]string {
	return map[string]string{
		"uk_name_version": "unique:name,version",
	}
}
