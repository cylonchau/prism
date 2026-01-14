package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// TerraformResourceOutputDAO provides resource output mapping data access operations.
type TerraformResourceOutputDAO struct {
	db *gorm.DB
}

// NewTerraformResourceOutputDAO creates a new resource output DAO.
func NewTerraformResourceOutputDAO(db *gorm.DB) *TerraformResourceOutputDAO {
	db.AutoMigrate(&models.TerraformResourceOutput{})
	return &TerraformResourceOutputDAO{db: db}
}

// Create creates a new output mapping.
func (d *TerraformResourceOutputDAO) Create(output *models.TerraformResourceOutput) error {
	return d.db.Create(output).Error
}

// Get retrieves output by ID.
func (d *TerraformResourceOutputDAO) Get(id int64) (*models.TerraformResourceOutput, error) {
	var output models.TerraformResourceOutput
	result := d.db.First(&output, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &output, nil
}

// ListByProvider lists outputs by provider.
func (d *TerraformResourceOutputDAO) ListByProvider(provider string) ([]models.TerraformResourceOutput, error) {
	var outputs []models.TerraformResourceOutput
	result := d.db.Where("provider = ?", provider).Find(&outputs)
	return outputs, result.Error
}

// ListByResourceType lists outputs by resource type.
func (d *TerraformResourceOutputDAO) ListByResourceType(resourceType string) ([]models.TerraformResourceOutput, error) {
	var outputs []models.TerraformResourceOutput
	result := d.db.Where("resource_type = ?", resourceType).Find(&outputs)
	return outputs, result.Error
}

// ListByProviderAndType lists outputs by provider and resource type.
func (d *TerraformResourceOutputDAO) ListByProviderAndType(provider, resourceType string) ([]models.TerraformResourceOutput, error) {
	var outputs []models.TerraformResourceOutput
	result := d.db.Where("provider = ? AND resource_type = ?", provider, resourceType).Find(&outputs)
	return outputs, result.Error
}

// GetByField gets output by provider, resource type and field.
func (d *TerraformResourceOutputDAO) GetByField(provider, resourceType, field string) (*models.TerraformResourceOutput, error) {
	var output models.TerraformResourceOutput
	result := d.db.Where("provider = ? AND resource_type = ? AND field = ?", provider, resourceType, field).First(&output)
	if result.Error != nil {
		return nil, result.Error
	}
	return &output, nil
}

// Update updates output.
func (d *TerraformResourceOutputDAO) Update(output *models.TerraformResourceOutput) error {
	return d.db.Save(output).Error
}

// Delete deletes output.
func (d *TerraformResourceOutputDAO) Delete(id int64) error {
	return d.db.Delete(&models.TerraformResourceOutput{}, id).Error
}
