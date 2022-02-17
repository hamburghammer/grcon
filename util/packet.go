package util

import "github.com/hamburghammer/grcon"

// NewExecCommandPacket creates a new grcon.Packet with
// the correct Type and the given Id and command string as Body.
func NewExecCommandPacket(id int32, cmd string) grcon.Packet {
	return grcon.Packet{
		Id:   grcon.PacketId(id),
		Type: grcon.SERVERDATA_EXECCOMMAND,
		Body: []byte(cmd),
	}
}

// NewEmptyResponseValuePacket returns an empty response value Type packet
// with an empty Body and the given Id.
func NewEmptyResponseValuePacket(id int32) grcon.Packet {
	return grcon.Packet{
		Id:   grcon.PacketId(id),
		Type: grcon.SERVERDATA_RESPONSE_VALUE,
		Body: []byte(""),
	}
}
