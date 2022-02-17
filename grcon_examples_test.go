package grcon_test

import (
	"log"
	"net"

	"github.com/hamburghammer/grcon"
)

func ExampleRemoteConsole_Read() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("establishing connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)

	packet, err := remoteConsole.Read()
	if err != nil {
		log.Fatalf("reading packet failed: %s", err.Error())
	}

	log.Printf("new packet read:\nid: %d\ntype: %d\nbody: %s\n", packet.Id, packet.Type, string(packet.Body))
}

func ExampleRemoteConsole_Write() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("establishing connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)
	packet := grcon.Packet{
		Id:   1,
		Type: grcon.SERVERDATA_EXECCOMMAND,
		Body: []byte("foo"),
	}

	err = remoteConsole.Write(packet)
	if err != nil {
		log.Fatalf("writing packet failed: %s", err.Error())
	}
}

func ExampleRemoteConsole() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("establishing connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)

	// authenticated connection.
	authPacket := grcon.Packet{
		Id:   1,
		Type: grcon.SERVERDATA_AUTH,
		Body: []byte("password"),
	}
	err = remoteConsole.Write(authPacket)
	authResponsePacket, err := remoteConsole.Read()
	if err != nil {
		log.Fatalf("auth failed: %s", err.Error())
	}
	if authResponsePacket.Id == -1 {
		log.Fatal("auth failed: auth rejected password invalid")
	}

	// packet that contais the command to execute
	cmdPacket := grcon.Packet{
		Id:   2,
		Type: grcon.SERVERDATA_EXECCOMMAND,
		Body: []byte("foo"),
	}
	err = remoteConsole.Write(cmdPacket)
	if err != nil {
		log.Fatalf("writing packet failed: %s", err.Error())
	}

	// empty packet that the server should mirror to indicate the end of the previous executed command.
	endPacket := grcon.Packet{
		Id:   3,
		Type: grcon.SERVERDATA_RESPONSE_VALUE,
		Body: []byte(""),
	}
	err = remoteConsole.Write(cmdPacket)
	if err != nil {
		log.Fatalf("writing packet failed: %s", err.Error())
	}

	responsePackets := make([]grcon.Packet, 0, 1)

	// read until endPacket is reached.
	for {
		packet, err := remoteConsole.Read()
		if err != nil {
			log.Fatalf("reading packet failed: %s", err.Error())
		}
		if packet.Id == endPacket.Id {
			break
		}
		responsePackets = append(responsePackets, packet)
	}

	var response string
	for _, packet := range responsePackets {
		response += string(packet.Body)
	}

	log.Println(response)
}
