package models

type TerraformResource struct {
	ID           int64  `gorm:"type:bigint;primaryKey;autoIncrement:false;comment:雪花算法" json:"id"`
	Provider     string `gorm:"type:varchar(64);not null;index:idx_provider" json:"provider"`
	ResourceType string `gorm:"type:varchar(64);not null;index:idx_resource_type" json:"resource_type"`
	RegionId     string `gorm:"type:varchar(128);index:idx_region" json:"region_id"`
	TfConfig     string `gorm:"type:text;comment:Terraform 配置文件" json:"tf_config"`
	TfState      string `gorm:"type:text;comment:Terraform 状态文件 (完整 tfstate)" json:"tf_state"`
	Action       string `gorm:"type:varchar(32);comment:操作类型 (apply, destroy, import)" json:"action"`
	Status       string `gorm:"type:varchar(32);default:'pending';comment:资源状态" json:"status"`
}

func (TerraformResource) TableName() string {
	return "terraform_resource"
}
