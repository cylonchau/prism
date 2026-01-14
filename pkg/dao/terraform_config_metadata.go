package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// TerraformConfigMetadataDAO provides config metadata data access operations.
type TerraformConfigMetadataDAO struct {
	db *gorm.DB
}

// NewTerraformConfigMetadataDAO creates a new config metadata DAO.
func NewTerraformConfigMetadataDAO(db *gorm.DB) *TerraformConfigMetadataDAO {
	db.AutoMigrate(&models.TerraformConfigMetadata{})
	return &TerraformConfigMetadataDAO{db: db}
}

// Create creates a new metadata.
func (d *TerraformConfigMetadataDAO) Create(meta *models.TerraformConfigMetadata) error {
	return d.db.Create(meta).Error
}

// Get retrieves metadata by ID.
func (d *TerraformConfigMetadataDAO) Get(id int64) (*models.TerraformConfigMetadata, error) {
	var meta models.TerraformConfigMetadata
	result := d.db.First(&meta, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meta, nil
}

// GetByAttribute retrieves metadata by attribute name.
func (d *TerraformConfigMetadataDAO) GetByAttribute(attribute string) (*models.TerraformConfigMetadata, error) {
	var meta models.TerraformConfigMetadata
	result := d.db.Where("attribute = ?", attribute).First(&meta)
	if result.Error != nil {
		return nil, result.Error
	}
	return &meta, nil
}

// ListByCategory lists metadata by category.
func (d *TerraformConfigMetadataDAO) ListByCategory(category string) ([]models.TerraformConfigMetadata, error) {
	var metas []models.TerraformConfigMetadata
	result := d.db.Where("category = ?", category).Find(&metas)
	return metas, result.Error
}

// ListRequired lists required metadata.
func (d *TerraformConfigMetadataDAO) ListRequired() ([]models.TerraformConfigMetadata, error) {
	var metas []models.TerraformConfigMetadata
	result := d.db.Where("is_required = ?", true).Find(&metas)
	return metas, result.Error
}

// List lists all metadata.
func (d *TerraformConfigMetadataDAO) List() ([]models.TerraformConfigMetadata, error) {
	var metas []models.TerraformConfigMetadata
	result := d.db.Find(&metas)
	return metas, result.Error
}

// Update updates metadata.
func (d *TerraformConfigMetadataDAO) Update(meta *models.TerraformConfigMetadata) error {
	return d.db.Save(meta).Error
}

// Delete deletes metadata.
func (d *TerraformConfigMetadataDAO) Delete(id int64) error {
	return d.db.Delete(&models.TerraformConfigMetadata{}, id).Error
}
