package grcon

// size of a packet or a field.
type size int32

// Sizes of the individual fields.
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Basic_Packet_Structure
const (
	sizeField    size = 4
	idField      size = 4
	typeField    size = 4
	minBodyField size = 1
	endField     size = 1
)

// MinPacket contains all fields except the size field.
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Packet_Size
const MinPacket size = idField + typeField + minBodyField + endField

// MaxPacket of a request/response packet.
// This size does not include the size field.
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#Packet_Size
const MaxPacket size = 4096

// MaxBody is the maximal size of the packet body.
const MaxBody size = MaxPacket - MinPacket
