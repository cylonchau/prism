package models

type TerraformResourceAttribute struct {
	ID             int64  `gorm:"type:bigint;primaryKey;autoIncrement:false" json:"id"`
	ResourceId     int64  `gorm:"type:bigint;not null;index:idx_resource_id;index:idx_resource_index" json:"resource_id"`
	ResourceIndex  int    `gorm:"not null;index:idx_resource_index;comment:这条资源在源数据中第几个索引" json:"resource_index"`
	AttributeName  string `gorm:"type:varchar(128);not null;index:idx_attribute_name;comment:源资源的参数名字" json:"attribute_name"`
	AttributeValue string `gorm:"type:text;not null;comment:源资源的值" json:"attribute_value"`
	ValueType      string `gorm:"type:varchar(32);not null;default:'string';comment:源资源的数据类型" json:"value_type"`
	MappedName     string `gorm:"type:varchar(128);not null;default:'';index:idx_mapped_name;comment:翻译后统一的名字用于整合系统" json:"mapped_name"`

	Resource *TerraformResource `gorm:"foreignKey:ResourceId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (TerraformResourceAttribute) TableName() string {
	return "terraform_resource_attribute"
}

func (TerraformResourceAttribute) Indexes() map[string]string {
	return map[string]string{
		"uk_resource_index_attr": "unique:resource_id,resource_index,attribute_name",
	}
}
