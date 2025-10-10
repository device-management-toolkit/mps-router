package db

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newSQLMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	return db, mock
}

func TestGetMPSInstance_Success(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	guid := "11111111-1111-1111-1111-111111111111"
	rows := sqlmock.NewRows([]string{"guid", "mpsinstance"}).AddRow(guid, "mps-host")
	mock.ExpectQuery(`SELECT guid, mpsinstance FROM devices WHERE guid = \$1;`).WithArgs(guid).WillReturnRows(rows)

	// Connect should return the injected connection
	_, err := pm.Connect()
	assert.NoError(t, err)

	got, err := pm.GetMPSInstance(pm.connection, guid)
	assert.NoError(t, err)
	assert.Equal(t, "mps-host", got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMPSInstance_NoRows(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	guid := "22222222-2222-2222-2222-222222222222"
	mock.ExpectQuery(`SELECT guid, mpsinstance FROM devices WHERE guid = \$1;`).WithArgs(guid).WillReturnError(sql.ErrNoRows)

	got, err := pm.GetMPSInstance(pm.connection, guid)
	assert.NoError(t, err)
	assert.Equal(t, "", got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMPSInstance_QueryError(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	guid := "33333333-3333-3333-3333-333333333333"
	mock.ExpectQuery(`SELECT guid, mpsinstance FROM devices WHERE guid = \$1;`).WithArgs(guid).WillReturnError(assert.AnError)

	got, err := pm.GetMPSInstance(pm.connection, guid)
	assert.Error(t, err)
	assert.Equal(t, "", got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealth_Success(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	ok := pm.Health()
	assert.True(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQuery_ErrorPath(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	guid := "44444444-4444-4444-4444-444444444444"
	mock.ExpectQuery(`SELECT guid, mpsinstance FROM devices WHERE guid = \$1;`).WithArgs(guid).WillReturnError(assert.AnError)

	got := pm.Query(guid)
	assert.Equal(t, "", got)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealth_QueryError(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	mock.ExpectQuery("SELECT 1").WillReturnError(assert.AnError)

	ok := pm.Health()
	assert.False(t, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQuery_Success(t *testing.T) {
	pm := &PostgresManager{}
	db, mock := newSQLMock(t)
	defer func() { _ = db.Close() }()
	pm.connection = db

	guid := "55555555-5555-5555-5555-555555555555"
	rows := sqlmock.NewRows([]string{"guid", "mpsinstance"}).AddRow(guid, "mps-instance-2")
	mock.ExpectQuery(`SELECT guid, mpsinstance FROM devices WHERE guid = \$1;`).WithArgs(guid).WillReturnRows(rows)

	got := pm.Query(guid)
	assert.Equal(t, "mps-instance-2", got)
	assert.NoError(t, mock.ExpectationsWereMet())
}
