package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// TerraformResourceDAO provides terraform resource data access operations.
type TerraformResourceDAO struct {
	db *gorm.DB
}

// NewTerraformResourceDAO creates a new terraform resource DAO.
func NewTerraformResourceDAO(db *gorm.DB) *TerraformResourceDAO {
	db.AutoMigrate(&models.TerraformResource{})
	return &TerraformResourceDAO{db: db}
}

// Create creates a new terraform resource.
func (d *TerraformResourceDAO) Create(resource *models.TerraformResource) error {
	return d.db.Create(resource).Error
}

// Get retrieves resource by ID.
func (d *TerraformResourceDAO) Get(id int64) (*models.TerraformResource, error) {
	var resource models.TerraformResource
	result := d.db.First(&resource, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resource, nil
}

// ListByProvider lists resources by provider.
func (d *TerraformResourceDAO) ListByProvider(provider string) ([]models.TerraformResource, error) {
	var resources []models.TerraformResource
	result := d.db.Where("provider = ?", provider).Find(&resources)
	return resources, result.Error
}

// ListByType lists resources by type.
func (d *TerraformResourceDAO) ListByType(resourceType string) ([]models.TerraformResource, error) {
	var resources []models.TerraformResource
	result := d.db.Where("resource_type = ?", resourceType).Find(&resources)
	return resources, result.Error
}

// ListByRegion lists resources by region.
func (d *TerraformResourceDAO) ListByRegion(regionID string) ([]models.TerraformResource, error) {
	var resources []models.TerraformResource
	result := d.db.Where("region_id = ?", regionID).Find(&resources)
	return resources, result.Error
}

// ListByStatus lists resources by status.
func (d *TerraformResourceDAO) ListByStatus(status string) ([]models.TerraformResource, error) {
	var resources []models.TerraformResource
	result := d.db.Where("status = ?", status).Find(&resources)
	return resources, result.Error
}

// Update updates resource.
func (d *TerraformResourceDAO) Update(resource *models.TerraformResource) error {
	return d.db.Save(resource).Error
}

// UpdateStatus updates resource status.
func (d *TerraformResourceDAO) UpdateStatus(id int64, status string) error {
	return d.db.Model(&models.TerraformResource{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateTfState updates tfstate.
func (d *TerraformResourceDAO) UpdateTfState(id int64, tfState string) error {
	return d.db.Model(&models.TerraformResource{}).Where("id = ?", id).Update("tf_state", tfState).Error
}

// Delete deletes resource.
func (d *TerraformResourceDAO) Delete(id int64) error {
	return d.db.Delete(&models.TerraformResource{}, id).Error
}
