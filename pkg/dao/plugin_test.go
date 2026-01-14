package dao

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/cylonchau/prism/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestPluginDAO_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewPluginDAO(db)

	plugin := &models.Plugin{
		Id:   "plugin-1",
		Name: "terraform-plugin",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `plugin`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := dao.Create(plugin)
	assert.NoError(t, err)
}

func TestPluginDAO_Get(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewPluginDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("plugin-1", "terraform-plugin")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `plugin` WHERE id = ?")).
		WithArgs("plugin-1").
		WillReturnRows(rows)

	plugin, err := dao.Get("plugin-1")
	assert.NoError(t, err)
	assert.Equal(t, "terraform-plugin", plugin.Name)
}

func TestPluginDAO_List(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewPluginDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("plugin-1", "plugin1").
		AddRow("plugin-2", "plugin2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `plugin`")).
		WillReturnRows(rows)

	plugins, err := dao.List()
	assert.NoError(t, err)
	assert.Len(t, plugins, 2)
}

func TestPluginDAO_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewPluginDAO(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `plugin` WHERE id = ?")).
		WithArgs("plugin-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Delete("plugin-1")
	assert.NoError(t, err)
}
