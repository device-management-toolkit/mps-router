/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type PostgresManager struct {
	ConnectionString string
	connection       *sql.DB
}

func NewPostgresManager(connectionString string) *PostgresManager {
	return &PostgresManager{
		ConnectionString: connectionString,
	}
}

func (pm *PostgresManager) Connect() (Database, error) {
	if pm.connection != nil {
		return pm.connection, nil
	}

	log.Println("Creating database connection pool")
	db, err := sql.Open("postgres", pm.ConnectionString)
	if err != nil {
		return nil, err
	}

	if maxOpenConnsStr, ok := os.LookupEnv("MPS_DB_MAX_OPEN_CONNS"); ok {
		maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
		if err != nil {
			return nil, err
		}

		db.SetMaxOpenConns(maxOpenConns)
	}

	pm.connection = db

	return db, nil
}

func (pm *PostgresManager) GetMPSInstance(db Database, guid string) (string, error) {
	client, ok := db.(*sql.DB)
	if !ok {
		return "", errors.New("invalid database type for PostgreSQL")
	}
	var device Device
	deviceSql := "SELECT guid, mpsinstance FROM devices WHERE guid = $1;"
	if client != nil {
		row := client.QueryRow(deviceSql, guid)
		switch err := row.Scan(&device.GUID, &device.MPSinstance); err {
		case sql.ErrNoRows:
			log.Println("no rows were returned!")
		case nil:
			{
				return device.MPSinstance, nil
			}
		default:
			{
				log.Println("failed to execute query: ", err)
				return "", err
			}
		}
		return "", nil
	}
	return "", errors.New("invalid db connection")
}

func (pm *PostgresManager) Health() bool {
	db, err := pm.Connect()
	if err != nil {
		log.Println("Failed to open a DB connection: ", err)
		return false
	}

	result := db.(*sql.DB).QueryRow("SELECT 1")
	if result.Err() != nil {
		log.Println(result.Err().Error())
		return false
	}
	return true
}

func (pm *PostgresManager) Query(guid string) string {
	db, err := pm.Connect()
	if err != nil {
		log.Println("Failed to open a DB connection: ", err)
		return ""
	}
	mpsInstance := ""
	mpsInstance, err = pm.GetMPSInstance(db, guid)
	if err != nil {
		log.Println("Failed to open a DB connection: ", err)
		return ""
	}
	return mpsInstance
}
