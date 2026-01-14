package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

type ProviderDAO struct {
	db *gorm.DB
}

// NewProviderDAO creates a new provider DAO.
func NewProviderDAO(db *gorm.DB) *ProviderDAO {
	db.AutoMigrate(&models.Provider{})
	return &ProviderDAO{db: db}
}

// Create creates a new provider.
func (d *ProviderDAO) Create(provider *models.Provider) error {
	return d.db.Create(provider).Error
}

// Get retrieves provider by ID.
func (d *ProviderDAO) Get(id int64) (*models.Provider, error) {
	var provider models.Provider
	result := d.db.First(&provider, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &provider, nil
}

// GetByName retrieves provider by name.
func (d *ProviderDAO) GetByName(name string) (*models.Provider, error) {
	var provider models.Provider
	result := d.db.Where("name = ?", name).First(&provider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &provider, nil
}

// GetByNameVersion retrieves provider by name and version.
func (d *ProviderDAO) GetByNameVersion(name, version string) (*models.Provider, error) {
	var provider models.Provider
	result := d.db.Where("name = ? AND version = ?", name, version).First(&provider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &provider, nil
}

// GetDefault retrieves the default provider by name.
func (d *ProviderDAO) GetDefault(name string) (*models.Provider, error) {
	var provider models.Provider
	result := d.db.Where("name = ? AND is_default = ?", name, true).First(&provider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &provider, nil
}

// List lists all providers.
func (d *ProviderDAO) List() ([]models.Provider, error) {
	var providers []models.Provider
	result := d.db.Find(&providers)
	return providers, result.Error
}

// ListEnabled lists enabled providers.
func (d *ProviderDAO) ListEnabled() ([]models.Provider, error) {
	var providers []models.Provider
	result := d.db.Where("enabled = ?", true).Find(&providers)
	return providers, result.Error
}

// Update updates provider.
func (d *ProviderDAO) Update(provider *models.Provider) error {
	return d.db.Save(provider).Error
}

// SetDefault sets provider as default.
func (d *ProviderDAO) SetDefault(id int64) error {
	// Clear existing default for same name
	var provider models.Provider
	if err := d.db.First(&provider, id).Error; err != nil {
		return err
	}
	d.db.Model(&models.Provider{}).Where("name = ?", provider.Name).Update("is_default", false)
	return d.db.Model(&models.Provider{}).Where("id = ?", id).Update("is_default", true).Error
}

// Delete deletes provider.
func (d *ProviderDAO) Delete(id int64) error {
	return d.db.Delete(&models.Provider{}, id).Error
}
