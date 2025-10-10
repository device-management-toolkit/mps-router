/*********************************************************************
 * Copyright (c) Intel Corporation 2021
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/device-management-toolkit/mps-router/internal/db"
	"github.com/device-management-toolkit/mps-router/internal/proxy"
)

func main() {
	// Keep main tiny and testable by delegating to run.
	code := run(
		os.Args[1:],
		os.Getenv,
		startServerReal,
		func(s string) db.Manager { return db.NewMongoManager(s) },
		func(s string) db.Manager { return db.NewPostgresManager(s) },
	)
	os.Exit(code)
}

func isMongoConnectionString(connectionString string) bool {
	return strings.HasPrefix(connectionString, "mongodb")
}

// run is the testable entry point for the application. It parses args/env,
// selects DB implementation, and starts the proxy server. It returns a process
// exit code (0=success, non-zero=failure) instead of exiting directly.
func run(
	args []string,
	getenv func(string) string,
	startServer func(db.Manager, string, string) error,
	newMongo func(string) db.Manager,
	newPostgres func(string) db.Manager,
	// return
) int {
	fs := flag.NewFlagSet("mps-router", flag.ContinueOnError)
	// Suppress default output in tests; errors will be handled via return code.
	fs.SetOutput(io.Discard)
	health := fs.Bool("health", false, "check health of service")
	if err := fs.Parse(args); err != nil {
		log.Println("failed to parse flags:", err)
		return 1
	}

	connectionString := getenv("MPS_CONNECTION_STRING")
	if connectionString == "" {
		// Preserve original message text to avoid surprising users/logs.
		log.Println("MPS_CONNECTION_STRING env is not set,default is mps")
		return 1
	}

	// Select DB implementation based on connection string.
	var dbImplementation db.Manager
	if isMongoConnectionString(connectionString) {
		dbImplementation = newMongo(connectionString)
	} else {
		dbImplementation = newPostgres(connectionString)
	}

	// Health check mode short-circuits server startup.
	if *health {
		if dbImplementation.Health() {
			return 0
		}
		return 1
	}

	// Resolve envs with defaults.
	routerPort := getenv("PORT")
	if routerPort == "" {
		log.Println("PORT env is not set, default is 8003")
		routerPort = "8003"
	}
	mpsPort := getenv("MPS_PORT")
	if mpsPort == "" {
		log.Println("MPS_PORT env is not set, default is 3000")
		mpsPort = "3000"
	}
	mpsHost := getenv("MPS_HOST")
	if mpsHost == "" {
		log.Println("MPS_HOST env is not set,default is mps")
		mpsHost = "mps"
	}

	addr := ":" + routerPort
	target := mpsHost + ":" + mpsPort
	if err := startServer(dbImplementation, addr, target); err != nil {
		log.Println("ListenAndServe:", err)
		return 1
	}
	return 0
}

// startServerReal constructs the proxy server and starts it. This is split out
// to allow tests to inject a fake to avoid binding a real port.
func startServerReal(m db.Manager, addr, target string) error {
	p := proxy.NewServer(m, addr, target)
	log.Println("Proxying from " + p.Addr + " to :" + p.Target)
	if err := p.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
