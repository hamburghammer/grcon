package client

import (
	"bytes"

	"github.com/hamburghammer/grcon"
	"github.com/hamburghammer/grcon/util"
)

// NewSimpleClient is a constructor for the SimpleClient struct.
// The util.GenerateNewId can be used as idGenFunc.
func NewSimpleClient(r util.RemoteConsole, idGenFunc func() grcon.PacketId) SimpleClient {
	return SimpleClient{RemoteConsole: r, IdGenFunc: idGenFunc}
}

// SimpleClient is a wrapper for a RemoteConsole that provides some utility functions.
// It simplifies the interaction with a remote console.
type SimpleClient struct {
	// RemoteConsole is the console to use for the interactions.
	util.RemoteConsole
	// IdGenFunc is the function to use to generate ids.
	IdGenFunc func() grcon.PacketId
}

// Auth should be used to authenticate the connection.
//
// It expects to receive an empty initial response value packet.
// https://developer.valvesoftware.com/wiki/Source_RCON_Protocol#SERVERDATA_AUTH_RESPONSE
//
// It can return following errors:
//	- InvalidResponseTypeError
//	- ResponseIdMismatchError
//	- ResponseBodyError
//	- AuthFailedError
func (sc SimpleClient) Auth(password string) error {
	reqID := sc.IdGenFunc()
	err := sc.Write(grcon.Packet{Id: reqID, Type: grcon.SERVERDATA_AUTH, Body: []byte(password)})
	if err != nil {
		return err
	}

	// read first empty SERVERDATA_RESPONSE_VALUE
	packet, err := sc.Read()
	if err != nil {
		return err
	}
	if packet.Type != grcon.SERVERDATA_RESPONSE_VALUE {
		return newInvalidResponseTypeError(grcon.SERVERDATA_RESPONSE_VALUE, packet.Type)
	}
	if packet.Id != reqID {
		return newResponseIdMismatchError(reqID, packet.Id)
	}
	// check if response is empty
	if !bytes.Equal(packet.Body, []byte{}) {
		return newResponseBodyError(string([]byte{}), string(packet.Body))
	}

	// read final SERVERDATA_AUTH_RESPONSE
	packet, err = sc.Read()
	if err != nil {
		return err
	}
	if packet.Type != grcon.SERVERDATA_AUTH_RESPONSE {
		return newInvalidResponseTypeError(grcon.SERVERDATA_AUTH_RESPONSE, packet.Type)
	}
	if packet.Id == -1 {
		return newAuthFailedError()
	}
	if packet.Id != reqID {
		return newResponseIdMismatchError(reqID, packet.Id)
	}

	return nil
}

// Exec executes the command on the given RemoteConsole implementation and
// waits till the response is read returns it.
// Supports multi-packet responses.
//
// The server has to response synchronously!
//
// Errors:
// Returns all errors returned from the Write and Read methode from the RemoteConsole implementation.
// Can also return an InvalidResponseTypeError if the response is not of the type
// grcon.SERVERDATA_RESPONSE_VALUE.
func (sc SimpleClient) Exec(cmd string) ([]byte, error) {
	cmdPacket := grcon.Packet{
		Id:   sc.IdGenFunc(),
		Type: grcon.SERVERDATA_EXECCOMMAND,
		Body: []byte(cmd),
	}
	err := sc.Write(cmdPacket)
	if err != nil {
		return []byte{}, err
	}
	// write delimiter packet
	delimiterPacket := grcon.Packet{
		Id:   sc.IdGenFunc(),
		Type: grcon.SERVERDATA_RESPONSE_VALUE,
		Body: []byte(""),
	}
	err = sc.Write(delimiterPacket)
	if err != nil {
		return []byte{}, err
	}

	// we assume that it won't be a multi packet response by giving the slice an initial capacity of 1.
	responsePackets := make([]grcon.Packet, 0, 1)

	// read until delimiterPacket is reached.
	for {
		packet, err := sc.Read()
		if err != nil {
			return []byte{}, err
		}
		if packet.Type != grcon.SERVERDATA_RESPONSE_VALUE {
			return []byte{}, newInvalidResponseTypeError(grcon.SERVERDATA_RESPONSE_VALUE, packet.Type)
		}
		// early break if delimiter packet is read.
		if packet.Id == delimiterPacket.Id {
			break
		}
		if packet.Id != cmdPacket.Id {
			return []byte{}, newResponseIdMismatchError(cmdPacket.Id, packet.Id)
		}
		responsePackets = append(responsePackets, packet)
	}

	// concatenate bodies
	response := make([]byte, 0, len(responsePackets))
	for _, packet := range responsePackets {
		response = append(response, packet.Body...)
	}

	return response, nil
}
