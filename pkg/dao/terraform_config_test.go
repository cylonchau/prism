package dao

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/cylonchau/prism/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestTerraformConfigDAO_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformConfigDAO(db)

	config := &models.TerraformConfig{
		ID:           1,
		Name:         "ec2-config",
		Provider:     "aws",
		BlockType:    "resource",
		ResourceType: "ec2",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `terraform_config`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := dao.Create(config)
	assert.NoError(t, err)
}

func TestTerraformConfigDAO_GetByName(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformConfigDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "provider"}).
		AddRow(1, "ec2-config", "aws")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `terraform_config` WHERE name = ?")).
		WithArgs("ec2-config").
		WillReturnRows(rows)

	config, err := dao.GetByName("ec2-config")
	assert.NoError(t, err)
	assert.Equal(t, "ec2-config", config.Name)
}

func TestTerraformConfigDAO_ListByProvider(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformConfigDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "provider"}).
		AddRow(1, "config1", "aws").
		AddRow(2, "config2", "aws")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `terraform_config` WHERE provider = ?")).
		WithArgs("aws").
		WillReturnRows(rows)

	configs, err := dao.ListByProvider("aws")
	assert.NoError(t, err)
	assert.Len(t, configs, 2)
}

func TestTerraformConfigDAO_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformConfigDAO(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `terraform_config` WHERE `terraform_config`.`id` = ?")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Delete(1)
	assert.NoError(t, err)
}
