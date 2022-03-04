package client

import (
	"github.com/hamburghammer/grcon"
	"github.com/hamburghammer/grcon/util"
)

// NewMinecraftClient is a constructor for the MinecraftClient struct.
// The util.GenerateNewId can be used as idGenFunc.
func NewMinecraftClient(r util.RemoteConsole, idGenFunc func() grcon.PacketId) MinecraftClient {
	return MinecraftClient{RemoteConsole: r, IdGenFunc: idGenFunc}
}

// MinecraftClient is a wrapper for a RemoteConsole that provides some utility functions.
// It simplifies the interaction with a remote console of a minecraft server.
type MinecraftClient struct {
	// RemoteConsole is the console to use for the interactions.
	util.RemoteConsole
	// IdGenFunc is the function to use to generate ids.
	IdGenFunc func() grcon.PacketId
}

// Auth should be used to authenticate the connection.
//
// It can return following errors:
//	- InvalidResponseTypeError
//	- ResponseIdMismatchError
//	- AuthFailedError
func (sc MinecraftClient) Auth(password string) error {
	reqID := sc.IdGenFunc()
	err := sc.Write(grcon.Packet{Id: reqID, Type: grcon.SERVERDATA_AUTH, Body: []byte(password)})
	if err != nil {
		return err
	}

	packet, err := sc.Read()
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
//
// Errors:
// Returns all errors returned from the Write and Read methode from the RemoteConsole implementation.
// Can also return an InvalidResponseTypeError if the response is not of the type
// grcon.SERVERDATA_RESPONSE_VALUE.
func (sc MinecraftClient) Exec(cmd string) ([]byte, error) {
	cmdPacket := grcon.Packet{
		Id:   sc.IdGenFunc(),
		Type: grcon.SERVERDATA_EXECCOMMAND,
		Body: []byte(cmd),
	}
	err := sc.Write(cmdPacket)
	if err != nil {
		return []byte{}, err
	}

	packet, err := sc.Read()
	if err != nil {
		return []byte{}, err
	}
	if packet.Type != grcon.SERVERDATA_RESPONSE_VALUE {
		return []byte{}, newInvalidResponseTypeError(grcon.SERVERDATA_RESPONSE_VALUE, packet.Type)
	}

	if packet.Id != cmdPacket.Id {
		return []byte{}, newResponseIdMismatchError(cmdPacket.Id, packet.Id)
	}

	return packet.Body, nil
}
