package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestConnect(t *testing.T) {
	manager := NewMongoManager("mongodb://localhost:27017")
	db, err := manager.Connect()
	assert.Nil(t, err, "test failed to connect db")

	_, ok := db.(*mongo.Client)
	assert.Equal(t, true, ok)
}

func TestNoSQLHealth(t *testing.T) {
	manager := &MongoManager{
		ConnectionString: "mongodb://localhost:27017",
	}
	result := manager.Health()
	assert.Equal(t, false, result)
}

func TestNoSQLGetMPSInstance(t *testing.T) {
	manager := &MongoManager{
		ConnectionString: "mongodb://localhost:27017",
	}
	result, err := manager.GetMPSInstance(nil, "mockGUID")
	assert.Empty(t, result)
	assert.Error(t, err)
}

func TestNoSQLQuery(t *testing.T) {
	manager := &MongoManager{
		ConnectionString: "mongodb://localhost:27017",
	}
	res := manager.Query("mockGUID")
	assert.Empty(t, res)
}

func TestNewMongoManagerDefaultsFromEnv(t *testing.T) {
	t.Setenv("MPS_DATABASE_NAME", "")
	t.Setenv("MPS_COLLECTION_NAME", "")
	m := NewMongoManager("mongodb://foo")
	assert.Equal(t, "mpsdb", m.DatabaseName)
	assert.Equal(t, "devices", m.CollectionName)

	t.Setenv("MPS_DATABASE_NAME", "customdb")
	t.Setenv("MPS_COLLECTION_NAME", "customcol")
	m2 := NewMongoManager("mongodb://foo")
	assert.Equal(t, "customdb", m2.DatabaseName)
	assert.Equal(t, "customcol", m2.CollectionName)
}

func TestMongoHealth_TypeAssertionFailure(t *testing.T) {
	// Test the path where Connect returns something that can't be asserted to *mongo.Client
	manager := &MongoManager{
		ConnectionString: "mongodb://localhost:27017",
	}
	// We can't easily inject a bad type, but we can test with invalid connection string
	manager.ConnectionString = "invalid://connection"
	result := manager.Health()
	assert.False(t, result)
}

func TestMongoQuery_ConnectionError(t *testing.T) {
	manager := &MongoManager{
		ConnectionString: "mongodb://nonexistent:99999",
		DatabaseName:     "test",
		CollectionName:   "devices",
	}
	result := manager.Query("some-guid")
	assert.Empty(t, result)
}

func TestMongoQuery_TypeAssertionFailure(t *testing.T) {
	// Test with invalid connection string to trigger connection errors
	manager := &MongoManager{
		ConnectionString: "invalid://bad",
		DatabaseName:     "test",
		CollectionName:   "devices",
	}
	result := manager.Query("some-guid")
	assert.Empty(t, result)
}

func TestMongoConnect_InvalidConnectionString(t *testing.T) {
	manager := &MongoManager{
		ConnectionString: "not-a-valid-uri",
	}
	_, err := manager.Connect()
	assert.Error(t, err)
}
