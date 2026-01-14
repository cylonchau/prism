package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// TerraformResourceAttributeDAO provides resource attribute data access operations.
type TerraformResourceAttributeDAO struct {
	db *gorm.DB
}

// NewTerraformResourceAttributeDAO creates a new resource attribute DAO.
func NewTerraformResourceAttributeDAO(db *gorm.DB) *TerraformResourceAttributeDAO {
	db.AutoMigrate(&models.TerraformResourceAttribute{})
	return &TerraformResourceAttributeDAO{db: db}
}

// Create creates a new attribute.
func (d *TerraformResourceAttributeDAO) Create(attr *models.TerraformResourceAttribute) error {
	return d.db.Create(attr).Error
}

// CreateBatch creates multiple attributes.
func (d *TerraformResourceAttributeDAO) CreateBatch(attrs []models.TerraformResourceAttribute) error {
	return d.db.CreateInBatches(attrs, 100).Error
}

// Get retrieves attribute by ID.
func (d *TerraformResourceAttributeDAO) Get(id int64) (*models.TerraformResourceAttribute, error) {
	var attr models.TerraformResourceAttribute
	result := d.db.First(&attr, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &attr, nil
}

// ListByResourceID lists attributes by resource ID.
func (d *TerraformResourceAttributeDAO) ListByResourceID(resourceID int64) ([]models.TerraformResourceAttribute, error) {
	var attrs []models.TerraformResourceAttribute
	result := d.db.Where("resource_id = ?", resourceID).Find(&attrs)
	return attrs, result.Error
}

// ListByResourceAndIndex lists attributes by resource ID and index.
func (d *TerraformResourceAttributeDAO) ListByResourceAndIndex(resourceID int64, index int) ([]models.TerraformResourceAttribute, error) {
	var attrs []models.TerraformResourceAttribute
	result := d.db.Where("resource_id = ? AND resource_index = ?", resourceID, index).Find(&attrs)
	return attrs, result.Error
}

// GetByResourceAndName gets attribute by resource ID and name.
func (d *TerraformResourceAttributeDAO) GetByResourceAndName(resourceID int64, name string) (*models.TerraformResourceAttribute, error) {
	var attr models.TerraformResourceAttribute
	result := d.db.Where("resource_id = ? AND attribute_name = ?", resourceID, name).First(&attr)
	if result.Error != nil {
		return nil, result.Error
	}
	return &attr, nil
}

// ListByMappedName lists attributes by mapped name.
func (d *TerraformResourceAttributeDAO) ListByMappedName(mappedName string) ([]models.TerraformResourceAttribute, error) {
	var attrs []models.TerraformResourceAttribute
	result := d.db.Where("mapped_name = ?", mappedName).Find(&attrs)
	return attrs, result.Error
}

// Update updates attribute.
func (d *TerraformResourceAttributeDAO) Update(attr *models.TerraformResourceAttribute) error {
	return d.db.Save(attr).Error
}

// Delete deletes attribute.
func (d *TerraformResourceAttributeDAO) Delete(id int64) error {
	return d.db.Delete(&models.TerraformResourceAttribute{}, id).Error
}

// DeleteByResourceID deletes all attributes for a resource.
func (d *TerraformResourceAttributeDAO) DeleteByResourceID(resourceID int64) error {
	return d.db.Where("resource_id = ?", resourceID).Delete(&models.TerraformResourceAttribute{}).Error
}
