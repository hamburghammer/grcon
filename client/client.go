// A package with client implementations to interact over RCON with different servers.
package client

// Client is the interface that should be implemented by ever client type.
type Client interface {
	// Auth authenticates the client against the server.
	Auth(password string) error
	// Exec executes a command on the remote console.
	Exec(cmd string) ([]byte, error)
}
