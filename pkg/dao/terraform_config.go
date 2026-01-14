package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// TerraformConfigDAO provides terraform config data access operations.
type TerraformConfigDAO struct {
	db *gorm.DB
}

// NewTerraformConfigDAO creates a new terraform config DAO.
func NewTerraformConfigDAO(db *gorm.DB) *TerraformConfigDAO {
	db.AutoMigrate(&models.TerraformConfig{})
	return &TerraformConfigDAO{db: db}
}

// Create creates a new config.
func (d *TerraformConfigDAO) Create(config *models.TerraformConfig) error {
	return d.db.Create(config).Error
}

// Get retrieves config by ID.
func (d *TerraformConfigDAO) Get(id int64) (*models.TerraformConfig, error) {
	var config models.TerraformConfig
	result := d.db.First(&config, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// GetByName retrieves config by name.
func (d *TerraformConfigDAO) GetByName(name string) (*models.TerraformConfig, error) {
	var config models.TerraformConfig
	result := d.db.Where("name = ?", name).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// ListByProvider lists configs by provider.
func (d *TerraformConfigDAO) ListByProvider(provider string) ([]models.TerraformConfig, error) {
	var configs []models.TerraformConfig
	result := d.db.Where("provider = ?", provider).Find(&configs)
	return configs, result.Error
}

// ListByBlockType lists configs by block type.
func (d *TerraformConfigDAO) ListByBlockType(blockType string) ([]models.TerraformConfig, error) {
	var configs []models.TerraformConfig
	result := d.db.Where("block_type = ?", blockType).Find(&configs)
	return configs, result.Error
}

// ListByResourceType lists configs by resource type.
func (d *TerraformConfigDAO) ListByResourceType(resourceType string) ([]models.TerraformConfig, error) {
	var configs []models.TerraformConfig
	result := d.db.Where("resource_type = ?", resourceType).Find(&configs)
	return configs, result.Error
}

// GetAttribute gets specific attribute config.
func (d *TerraformConfigDAO) GetAttribute(provider, blockType, resourceType, attribute string) (*models.TerraformConfig, error) {
	var config models.TerraformConfig
	result := d.db.Where("provider = ? AND block_type = ? AND resource_type = ? AND attribute = ?",
		provider, blockType, resourceType, attribute).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// Update updates config.
func (d *TerraformConfigDAO) Update(config *models.TerraformConfig) error {
	return d.db.Save(config).Error
}

// Delete deletes config.
func (d *TerraformConfigDAO) Delete(id int64) error {
	return d.db.Delete(&models.TerraformConfig{}, id).Error
}
