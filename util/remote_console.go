package util

import "github.com/hamburghammer/grcon"

// RemoteConsole is an interface that the grcon.RemoteConsole struct implements.
type RemoteConsole interface {
	// Read a packet
	Read() (grcon.Packet, error)
	// Write a packet
	Write(grcon.Packet) error
}
