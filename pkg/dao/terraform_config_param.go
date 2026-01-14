package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// TerraformConfigParamDAO provides config param data access operations.
type TerraformConfigParamDAO struct {
	db *gorm.DB
}

// NewTerraformConfigParamDAO creates a new config param DAO.
func NewTerraformConfigParamDAO(db *gorm.DB) *TerraformConfigParamDAO {
	db.AutoMigrate(&models.TerraformConfigParam{})
	return &TerraformConfigParamDAO{db: db}
}

// Create creates a new param.
func (d *TerraformConfigParamDAO) Create(param *models.TerraformConfigParam) error {
	return d.db.Create(param).Error
}

// Get retrieves param by ID.
func (d *TerraformConfigParamDAO) Get(id int64) (*models.TerraformConfigParam, error) {
	var param models.TerraformConfigParam
	result := d.db.First(&param, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &param, nil
}

// ListByConfigID lists params by config ID.
func (d *TerraformConfigParamDAO) ListByConfigID(configID int64) ([]models.TerraformConfigParam, error) {
	var params []models.TerraformConfigParam
	result := d.db.Where("terraform_config_id = ?", configID).Find(&params)
	return params, result.Error
}

// GetByConfigAndName gets param by config ID and param name.
func (d *TerraformConfigParamDAO) GetByConfigAndName(configID int64, paramName string) (*models.TerraformConfigParam, error) {
	var param models.TerraformConfigParam
	result := d.db.Where("terraform_config_id = ? AND param_name = ?", configID, paramName).First(&param)
	if result.Error != nil {
		return nil, result.Error
	}
	return &param, nil
}

// Update updates param.
func (d *TerraformConfigParamDAO) Update(param *models.TerraformConfigParam) error {
	return d.db.Save(param).Error
}

// UpdateValue updates param value.
func (d *TerraformConfigParamDAO) UpdateValue(id int64, value string) error {
	return d.db.Model(&models.TerraformConfigParam{}).Where("id = ?", id).Update("param_value", value).Error
}

// Delete deletes param.
func (d *TerraformConfigParamDAO) Delete(id int64) error {
	return d.db.Delete(&models.TerraformConfigParam{}, id).Error
}

// DeleteByConfigID deletes all params for a config.
func (d *TerraformConfigParamDAO) DeleteByConfigID(configID int64) error {
	return d.db.Where("terraform_config_id = ?", configID).Delete(&models.TerraformConfigParam{}).Error
}
