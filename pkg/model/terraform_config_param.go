package models

type TerraformConfigParam struct {
	ID                int64  `gorm:"type:bigint;primaryKey;autoIncrement:false;comment:雪花算法 ID" json:"id"`
	TerraformConfigID int64  `gorm:"type:bigint;not null;index:idx_terraform_config_id;index:idx_config_param;comment:配置 ID (FK to terraform_config.id)" json:"terraform_config_id"`
	ParamName         string `gorm:"type:varchar(128);not null;index:idx_resource_param;comment:参数名" json:"param_name"`
	ParamValue        string `gorm:"type:text;not null;default:'';comment:参数值" json:"param_value"`
	DefaultValue      string `gorm:"type:text;not null;default:'';comment:参数默认值" json:"default_value"`
	ValueType         string `gorm:"type:varchar(128);not null;default:'string';comment:参数值类型" json:"value_type"`
	Description       string `gorm:"type:varchar(128);not null;default:'';comment:参数描述" json:"description"`

	TerraformConfig *TerraformConfig `gorm:"foreignKey:TerraformConfigID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (TerraformConfigParam) TableName() string {
	return "terraform_param"
}

func (TerraformConfigParam) Indexes() map[string]string {
	return map[string]string{
		"uk_terraform_config_param": "unique:terraform_config_id,param_name",
	}
}
