package models

type TerraformResourceOutput struct {
	ID           int64  `gorm:"type:bigint;primaryKey;autoIncrement:false" json:"id"`
	Provider     string `gorm:"type:varchar(100);not null;index:idx_provider;comment:terraform provider" json:"provider"`
	ResourceType string `gorm:"type:varchar(64);not null;index:idx_resource_type;comment:资源类型" json:"resource_type"`
	Field        string `gorm:"type:varchar(128);not null;comment:要提取的字段名" json:"field"`
	TfStatePath  string `gorm:"type:varchar(256);not null;comment:tfstate 中的路径" json:"tfstate_path"`
	DisplayName  string `gorm:"type:varchar(128);comment:显示名称" json:"display_name"`
	Description  string `gorm:"type:varchar(512);comment:字段说明" json:"description"`
}

func (TerraformResourceOutput) TableName() string {
	return "terraform_resource_output"
}

func (TerraformResourceOutput) Indexes() map[string]string {
	return map[string]string{
		"uk_provider_resource_field": "unique:provider,resource_type,field",
	}
}
