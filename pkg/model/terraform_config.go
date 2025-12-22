package models

type TerraformConfig struct {
	ID           int64  `gorm:"type:bigint;primaryKey;autoIncrement:false" json:"id"`
	Name         string `gorm:"type:varchar(100);uniqueIndex;comment:名称;not null" json:"name"`
	Provider     string `gorm:"type:varchar(32);comment:Entity: 云厂商 (tencentcloud, alicloud, aws);not null;default:''" json:"provider" xorm:"provider"`
	BlockType    string `gorm:"type:varchar(32);not null;index:idx_block_type;comment:配置块类型" json:"block_type"`
	ResourceType string `gorm:"type:varchar(64);comment:资源类型 (cvm, vpc, subnet) ;not null;default:''" json:"resource_type" xorm:"resource_type"`
	Attribute    string `gorm:"type:varchar(128);comment:配置项名称;not null;default:''" json:"attribute" xorm:"attribute"`
	Value        string `gorm:"type:text;comment:配置项值 (支持 JSON) ;not null;default:''" json:"value" xorm:"value"`
	ValueType    string `gorm:"type:varchar(32);comment:值类型: string, json, int, bool;not null;default:'string'" json:"value_type" xorm:"value_type"`
	Description  string `gorm:"type:varchar(256);comment:说明;not null;default:''" json:"description" xorm:"description"`
}

func (TerraformConfig) TableName() string {
	return "terraform_config"
}

func (TerraformConfig) Indexes() map[string]string {
	return map[string]string{
		"uk_provider_block_resource_attr": "unique:provider,block_type,resource_type,attribute",
	}
}
