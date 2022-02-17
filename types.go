package grcon

// PacketId is the id for a packet.
// It may be set to any positive integer.
// It need not be unique, but if a unique packet id is assigned,
// it can be used to match incoming responses to their corresponding requests.
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Packet_ID
type PacketId int32

// PacketType indicate the purpose of the packet.
//
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Packet_Type
type PacketType int32

// Package types
const (
	// SERVERDATA_AUTH Typically, the first packet sent by the client will be a SERVERDATA_AUTH packet,
	// which is used to authenticate the connection with the server.
	// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#SERVERDATA_AUTH
	SERVERDATA_AUTH PacketType = 3

	// SERVERDATA_EXECCOMMAND packet type represents a command issued to the server by a client.
	// This can be a ConCommand such as mp_switchteams or changelevel, a command to set a cvar such as sv_cheats 1,
	// or a command to fetch the value of a cvar, such as sv_cheats.
	// The response will vary depending on the command issued.
	// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#SERVERDATA_EXECCOMMAND
	SERVERDATA_EXECCOMMAND PacketType = 2

	// SERVERDATA_RESPONSE_VALUE packet is the response to a SERVERDATA_EXECCOMMAND request.
	// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#SERVERDATA_RESPONSE_VALUE
	SERVERDATA_RESPONSE_VALUE PacketType = 0

	// This packet is a notification of the connection's current auth status.
	// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#SERVERDATA_AUTH_RESPONSE
	SERVERDATA_AUTH_RESPONSE PacketType = 2
)

type Packet struct {
	Id   PacketId
	Type PacketType
	Body []byte
}
