package test

import (
	"database/sql"

	"github.com/device-management-toolkit/mps-router/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockSQLDBManager struct {
	ConnectResult     *sql.DB
	ConnectError      error
	ConnectionStr     string
	MPSInstanceResult string
	MPSInstanceError  error
	HealthResult      bool
	QueryResult       string
}

func (mock *MockSQLDBManager) Connect() (db.Database, error) {
	return mock.ConnectResult, mock.ConnectError
}

func (mock *MockSQLDBManager) GetMPSInstance(db db.Database, guid string) (string, error) {
	if mock.MPSInstanceError != nil {
		return "", mock.MPSInstanceError
	}
	return mock.MPSInstanceResult, nil
}

func (mock *MockSQLDBManager) Health() bool {
	return mock.HealthResult
}

func (mock *MockSQLDBManager) Query(guid string) string {
	return mock.QueryResult
}

type MockNOSQLDBManager struct {
	ConnectResult     *mongo.Client
	ConnectError      error
	ConnectionStr     string
	MPSInstanceResult string
	MPSInstanceError  error
	HealthResult      bool
	QueryResult       string
}

func (mock *MockNOSQLDBManager) Connect() (db.Database, error) {
	return mock.ConnectResult, mock.ConnectError
}

func (mock *MockNOSQLDBManager) GetMPSInstance(db db.Database, guid string) (string, error) {
	if mock.MPSInstanceError != nil {
		return "", mock.MPSInstanceError
	}
	return mock.MPSInstanceResult, nil
}

func (mock *MockNOSQLDBManager) Health() bool {
	return mock.HealthResult
}

func (mock *MockNOSQLDBManager) Query(guid string) string {
	return mock.QueryResult
}
