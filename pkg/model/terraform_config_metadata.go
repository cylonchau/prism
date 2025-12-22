package models

type TerraformConfigMetadata struct {
	ID             int64  `gorm:"type:bigint;primaryKey;autoIncrement:false" json:"id"`
	Attribute      string `gorm:"type:varchar(128);comment:属性名;not null;default:''" json:"attribute" xorm:"attribute"`
	DisplayName    string `gorm:"type:varchar(128);comment:显示名称;not null;default:''" json:"display_name" xorm:"display_name"`
	ValueType      string `gorm:"type:char(32);comment:值类型;not null;default:''" json:"value_type" xorm:"value_type"`
	IsRequired     bool   `gorm:"not null;default:false;comment:是否必填" json:"is_required" xorm:"is_required"`
	DefaultValue   string `gorm:"type:text;not null;default:'';comment:默认值" json:"default_value" xorm:"default_value"`
	ValidationRule string `gorm:"type:text;not null;default:'';comment:验证规则 (JSON)" json:"validation_rule" xorm:"validation_rule"`
	Description    string `gorm:"type:varchar(256);comment:说明;not null;default:''" json:"description" xorm:"description"`
	Category       string `gorm:"type:varchar(64);comment:分类;not null;default:''" json:"category" xorm:"category"`
}

func (TerraformConfigMetadata) TableName() string {
	return "terraform_config_metadata"
}
