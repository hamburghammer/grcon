package util

import (
	"time"

	"github.com/hamburghammer/grcon"
)

// GenerateRequestId is a convenience function to generate an id using the current time.
func GenerateRequestId() grcon.PacketId {
	return grcon.PacketId((time.Now().UnixNano() / 100000) % 100000)
}
