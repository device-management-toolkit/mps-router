/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/device-management-toolkit/mps-router/internal/db"
	"github.com/device-management-toolkit/mps-router/internal/proxy"
)

func main() {

	result := flag.Bool("health", false, "check health of service")
	flag.Parse()
	connectionString := os.Getenv("MPS_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("MPS_CONNECTION_STRING env is not set,default is mps")
	}
	var dbImplementation db.Manager
	if isMongoConnectionString(connectionString) {
		// Handle MongoDB-related operations.
		dbImplementation = db.NewMongoManager(connectionString)
	} else {
		// Handle sql database-related operations.
		dbImplementation = db.NewPostgresManager(connectionString)
	}

	if *result {
		dbHealth := dbImplementation.Health()
		if dbHealth {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	routerPort := os.Getenv("PORT")
	if routerPort == "" {
		log.Println("PORT env is not set, default is 8003")
		routerPort = "8003"
	}
	mpsPort := os.Getenv("MPS_PORT")
	if mpsPort == "" {
		log.Println("MPS_PORT env is not set, default is 3000")
		mpsPort = "3000"
	}
	mpsHost := os.Getenv("MPS_HOST")
	if mpsHost == "" {
		log.Println("MPS_HOST env is not set,default is mps")
		mpsHost = "mps"
	}

	p := proxy.NewServer(dbImplementation, ":"+routerPort, mpsHost+":"+mpsPort)
	log.Println("Proxying from " + p.Addr + " to :" + p.Target)
	err := p.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func isMongoConnectionString(connectionString string) bool {
	return strings.HasPrefix(connectionString, "mongodb")
}
