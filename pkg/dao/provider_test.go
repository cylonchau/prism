package dao

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	models "github.com/cylonchau/prism/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestProviderDAO_Create(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewProviderDAO(db)

	provider := &models.Provider{
		ID:      1,
		Name:    "aws",
		Version: "5.0.0",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `provider`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := dao.Create(provider)
	assert.NoError(t, err)
}

func TestProviderDAO_Get(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewProviderDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "version"}).
		AddRow(1, "aws", "5.0.0")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provider` WHERE `provider`.`id` = ?")).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	provider, err := dao.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, "aws", provider.Name)
}

func TestProviderDAO_GetByName(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewProviderDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "version"}).
		AddRow(1, "aws", "5.0.0")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provider` WHERE name = ?")).
		WithArgs("aws").
		WillReturnRows(rows)

	provider, err := dao.GetByName("aws")
	assert.NoError(t, err)
	assert.Equal(t, "aws", provider.Name)
}

func TestProviderDAO_List(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewProviderDAO(db)

	rows := sqlmock.NewRows([]string{"id", "name", "version"}).
		AddRow(1, "aws", "5.0.0").
		AddRow(2, "azure", "3.0.0")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `provider`")).
		WillReturnRows(rows)

	providers, err := dao.List()
	assert.NoError(t, err)
	assert.Len(t, providers, 2)
}

func TestProviderDAO_Delete(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	dao := NewProviderDAO(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `provider` WHERE `provider`.`id` = ?")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dao.Delete(1)
	assert.NoError(t, err)
}
