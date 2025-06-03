// Package db provides abstractions for database operations.
package db

// Database represents a universal database object.
// For different databases, different underlying types can be used.
// For example, MongoDB would use a *mongo.Client, and SQL would use *sql.DB.
type Database interface{}

// Manager provides an interface for database management operations.
// It offers methods for connecting to the database, querying data,
// checking the health of the connection, and retrieving MPS instances by GUID.
type Manager interface {
	// Connect verifies the connection string to the database and returns the Database object. It does not guarantee that the connection is healthy.
	Connect() (Database, error)

	// GetMPSInstance fetches the MPS instance associated with the given GUID from the specified Database.
	GetMPSInstance(db Database, guid string) (string, error)

	// Health checks the status of the database connection.
	// It returns true if the connection is healthy; otherwise false.
	Health() bool

	// Query retrieves a result from the database based on the provided GUID.
	// The implementation details can vary depending on the underlying database.
	Query(guid string) string
}

// Device represents a database entity with information about a device.
// It includes a globally unique identifier (GUID) and an associated MPS instance, if available.
type Device struct {
	// GUID represents the Globally Unique Identifier for the device.
	GUID string `bson:"guid"`

	// MPSinstance holds the MPS instance associated with the device.
	MPSinstance string `bson:"mpsinstance"`
}
