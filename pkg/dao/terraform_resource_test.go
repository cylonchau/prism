package dao

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/cylonchau/prism/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestTerraformResourceDAO_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformResourceDAO(db)

	resource := &models.TerraformResource{
		ID:           1,
		Provider:     "aws",
		ResourceType: "ec2",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `terraform_resource`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := dao.Create(resource)
	assert.NoError(t, err)
}

func TestTerraformResourceDAO_Get(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformResourceDAO(db)

	rows := sqlmock.NewRows([]string{"id", "provider", "resource_type"}).
		AddRow(1, "aws", "ec2")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `terraform_resource` WHERE `terraform_resource`.`id` = ?")).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	resource, err := dao.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, "aws", resource.Provider)
}

func TestTerraformResourceDAO_UpdateStatus(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformResourceDAO(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `terraform_resource` SET `status`=? WHERE id = ?")).
		WithArgs("running", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.UpdateStatus(1, "running")
	assert.NoError(t, err)
}

func TestTerraformResourceDAO_ListByProvider(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformResourceDAO(db)

	rows := sqlmock.NewRows([]string{"id", "provider", "resource_type"}).
		AddRow(1, "aws", "ec2").
		AddRow(2, "aws", "vpc")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `terraform_resource` WHERE provider = ?")).
		WithArgs("aws").
		WillReturnRows(rows)

	resources, err := dao.ListByProvider("aws")
	assert.NoError(t, err)
	assert.Len(t, resources, 2)
}

func TestTerraformResourceDAO_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewTerraformResourceDAO(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `terraform_resource` WHERE `terraform_resource`.`id` = ?")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Delete(1)
	assert.NoError(t, err)
}
