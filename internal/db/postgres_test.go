/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package db

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectToDB(t *testing.T) {
	var db *sql.DB
	pm := PostgresManager{}
	result, err := pm.Connect()
	assert.Nil(t, err, "test failed to connect db")
	assert.Equal(t, reflect.TypeOf(result), reflect.TypeOf(db))
}

func TestConnectionPoolConfiguration(t *testing.T) {
	pm := PostgresManager{}
	t.Setenv("MPS_DB_MAX_OPEN_CONNS", "7")
	_, err := pm.Connect()
	assert.Nil(t, err, "test failed to connect db")
	assert.Equal(t, 7, pm.connection.Stats().MaxOpenConnections, "connection pool max open conns not configured")
}

func TestConnectionPoolConfigurationInvalid(t *testing.T) {
	pm := PostgresManager{}
	t.Setenv("MPS_DB_MAX_OPEN_CONNS", "not-a-number")
	_, err := pm.Connect()
	assert.Error(t, err)
}

func TestGetMPSInstancewithGUID(t *testing.T) {
	pm := PostgresManager{}

	db, err := pm.Connect()
	assert.Nil(t, err, "test failed to connect db")
	result := ""
	result, err = pm.GetMPSInstance(db, "d12428be-9fa1-4226-9784-54b2038beab6")
	if err != nil {
		log.Println("test failed to get the mps instance", err)
	}
	assert.Equal(t, "", result)
}

func TestGetMPSInstancewithNoDB(t *testing.T) {
	pm := PostgresManager{}

	var db *sql.DB
	_, err := pm.GetMPSInstance(db, "d12428be-9fa1-4226-9784-54b2038beab6")
	if err != nil {
		log.Println("test failed to get the mps instance", err)
	}
	assert.Equal(t, "invalid db connection", err.Error())
}

func TestQuery(t *testing.T) {
	pm := PostgresManager{}

	// Set an Environment Variable
	t.Setenv("MPS_CONNECTION_STRING", "postgresql://")
	result := pm.Query("d12428be-9fa1-4226-9784-54b2038beab6")
	assert.Equal(t, "", result)
}

func TestHealth(t *testing.T) {
	pm := PostgresManager{}
	result := pm.Health()
	assert.Equal(t, false, result)
}

func TestGetMPSInstance_InvalidDatabaseType(t *testing.T) {
	pm := PostgresManager{}
	// Pass something that's not a *sql.DB to trigger type assertion failure
	invalidDB := "not-a-database"
	result, err := pm.GetMPSInstance(invalidDB, "some-guid")
	assert.Error(t, err)
	assert.Equal(t, "invalid database type for PostgreSQL", err.Error())
	assert.Empty(t, result)
}

func TestConnect_ReuseExistingConnection(t *testing.T) {
	pm := PostgresManager{}
	// First connection
	db1, err1 := pm.Connect()
	assert.NoError(t, err1)
	assert.NotNil(t, db1)

	// Second call should reuse the same connection
	db2, err2 := pm.Connect()
	assert.NoError(t, err2)
	assert.NotNil(t, db2)
	assert.Same(t, db1, db2, "Connect should reuse existing connection")
}
