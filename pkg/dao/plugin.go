package dao

import (
	models "github.com/cylonchau/prism/pkg/model"
	"gorm.io/gorm"
)

// PluginDAO provides plugin data access operations.
type PluginDAO struct {
	db *gorm.DB
}

// NewPluginDAO creates a new plugin DAO.
func NewPluginDAO(db *gorm.DB) *PluginDAO {
	db.AutoMigrate(&models.Plugin{})
	return &PluginDAO{db: db}
}

// Create creates a new plugin.
func (d *PluginDAO) Create(plugin *models.Plugin) error {
	return d.db.Create(plugin).Error
}

// Get retrieves plugin by ID.
func (d *PluginDAO) Get(id string) (*models.Plugin, error) {
	var plugin models.Plugin
	result := d.db.First(&plugin, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &plugin, nil
}

// GetByName retrieves plugin by name.
func (d *PluginDAO) GetByName(name string) (*models.Plugin, error) {
	var plugin models.Plugin
	result := d.db.Where("name = ?", name).First(&plugin)
	if result.Error != nil {
		return nil, result.Error
	}
	return &plugin, nil
}

// List lists all plugins.
func (d *PluginDAO) List() ([]models.Plugin, error) {
	var plugins []models.Plugin
	result := d.db.Find(&plugins)
	return plugins, result.Error
}

// Update updates plugin.
func (d *PluginDAO) Update(plugin *models.Plugin) error {
	return d.db.Save(plugin).Error
}

// Delete deletes plugin.
func (d *PluginDAO) Delete(id string) error {
	return d.db.Delete(&models.Plugin{}, "id = ?", id).Error
}
