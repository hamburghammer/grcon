/*
Package grcon library for the RCON Protocol from Valve.

Information to the protocol can be found under:
https://developer.valvesoftware.com/wiki/Source_RCON_Protocol
*/
package grcon

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
)

// NewRemoteConsole creates a new RemoteConsole with the given connection and default values.
func NewRemoteConsole(conn net.Conn) *RemoteConsole {
	remoteConsole := &RemoteConsole{
		Conn: conn,
		// Initialize a buffer to hold at least a hole packet:
		// Since the only one of the packet values that can change in length is the body,
		// an easy way to calculate the size of a packet is to find the byte-length of the packet body,
		// then add 10 to it.
		ReadBuff: make([]byte, MaxPacket+sizeField),
	}

	return remoteConsole
}

// RemoteConsole holds the information to communicate withe remote console (server).
// To optain a preconfigured RemoteConsole use the NewRemoteConsole() function.
// The NewRemoteConsole() function is also the recommended way to get a *RemoteConsole.
//
// This struct can be used concurrently.
// All exported fields are not allowed to be nil!
type RemoteConsole struct {
	// Conn is the connection to read and write to.
	Conn net.Conn

	// ReadBuff should at least have the capacity for a hole packet.
	// Capacity >= MaxPacket + 4.
	ReadBuff []byte

	readMutex  sync.Mutex
	queuedBuff []byte
}

// Write writes a packet with a given id, type and body.
// The body should be a ASCII string.
// Returns an RequestTooLongError if the body is greater than the max
// length which is MaxPacket - MinPacket.
func (r *RemoteConsole) Write(packet Packet) error {
	bodySize := size(len(packet.Body))
	if bodySize > MaxBody {
		return newRequestTooLongError()
	}

	buffer := bytes.NewBuffer(make([]byte, 0, MinPacket+sizeField+bodySize))

	// size
	err := binary.Write(buffer, binary.LittleEndian, bodySize+MinPacket)
	if err != nil {
		return err
	}

	// id
	err = binary.Write(buffer, binary.LittleEndian, packet.Id)
	if err != nil {
		return err
	}

	// type
	err = binary.Write(buffer, binary.LittleEndian, packet.Type)
	if err != nil {
		return err
	}

	// body
	_, err = buffer.Write(packet.Body)
	if err != nil {
		return err
	}

	// double null termination
	err = binary.Write(buffer, binary.LittleEndian, byte(0))
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.LittleEndian, byte(0))
	if err != nil {
		return err
	}

	// writing to the connection
	_, err = r.Conn.Write(buffer.Bytes())

	return err
}

// Read returns all the parts of the read packet.
// Returns an ResponseTooLongError if the size of the packet is bigger
// than the MaxPacket siz. It can also return an UnexpectedForamatError
// if the packet size is smaller than the MinPacket size.
func (r *RemoteConsole) Read() (Packet, error) {
	r.readMutex.Lock()
	defer r.readMutex.Unlock()

	var readBytes int
	var err error
	if r.queuedBuff != nil {
		copy(r.ReadBuff, r.queuedBuff)
		readBytes = len(r.queuedBuff)
		r.queuedBuff = nil
	} else {
		readBytes, err = r.Conn.Read(r.ReadBuff)
		if err != nil {
			return Packet{}, err
		}
	}

	dataSize, readBytes, err := r.readPacketSize(readBytes)
	if err != nil {
		return Packet{}, err
	}

	if dataSize > MaxPacket {
		return Packet{}, newResponseTooLongError()
	}

	totalPacketSize := dataSize + sizeField
	readBytes, err = r.readPacket(totalPacketSize, readBytes)
	if err != nil {
		return Packet{}, err
	}

	// The data has to be explicitly selected to prevent copying empty bytes.
	data := r.ReadBuff[sizeField:totalPacketSize]

	// Save not packet related bytes for the next read.
	if readBytes > int(totalPacketSize) {
		// start of the next buffer was at the end of this packet.
		// save it for the next read.
		// The data has to be explicitly selected to prevent copying empty bytes.
		r.queuedBuff = r.ReadBuff[totalPacketSize:readBytes]
	}

	return r.parsePacket(data)
}

// readPacketSize wait until first 4 bytes are read to get the packet size.
// Takes as param how many bytes are already read. The returned size does not include the size field.
func (r *RemoteConsole) readPacketSize(readBytes int) (size, int, error) {
	for readBytes < int(sizeField) {
		// need the 4 byte packet size...
		b, err := r.Conn.Read(r.ReadBuff[readBytes:])
		if err != nil {
			return 0, 0, err
		}
		readBytes += b
	}

	// Does not include the packetSize field.
	var totalPacketSize size
	b := bytes.NewBuffer(r.ReadBuff[:sizeField])
	err := binary.Read(b, binary.LittleEndian, &totalPacketSize)
	if err != nil {
		return 0, 0, err
	}

	if totalPacketSize < MinPacket {
		return 0, 0, newUnexpectedFormatError()
	}

	return totalPacketSize, readBytes, nil
}

// readPacket waits until the whole packet is read including the size field.
func (r *RemoteConsole) readPacket(totalPacketSize size, readBytes int) (int, error) {
	for int(totalPacketSize) > readBytes {
		b, err := r.Conn.Read(r.ReadBuff[readBytes:])
		if err != nil {
			return readBytes, err
		}

		readBytes += b
	}

	return readBytes, nil
}

// parsePacket reads the a packet from an array byte.
// The array has only to contain the 'id', 'type' and 'body' data.
func (r *RemoteConsole) parsePacket(data []byte) (Packet, error) {
	var requestID, responseType int32
	buffer := bytes.NewBuffer(data)

	err := binary.Read(buffer, binary.LittleEndian, &requestID)
	if err != nil {
		return Packet{}, err
	}

	binary.Read(buffer, binary.LittleEndian, &responseType)
	if err != nil {
		return Packet{}, err
	}

	// the rest of the buffer is the body.
	body := buffer.Bytes()
	// remove the to null terminations
	body = body[:len(body)-2]

	parsedPacket := Packet{
		Id:   PacketId(requestID),
		Type: PacketType(responseType),
		Body: body,
	}

	return parsedPacket, nil
}
