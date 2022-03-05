package client_test

import (
	"log"
	"net"

	"github.com/hamburghammer/grcon"
	"github.com/hamburghammer/grcon/client"
	"github.com/hamburghammer/grcon/util"
)

func ExampleNewMinecraftClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)

	// the returned MinecraftClient can now be used.
	// It will use the utility function to generate ids
	_ = client.NewMinecraftClient(remoteConsole, util.GenerateRequestId)
}

func ExampleMinecraftClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)

	minecraftClient := client.MinecraftClient{
		RemoteConsole: remoteConsole,
		// Use the utility function to generate ids
		IdGenFunc: util.GenerateRequestId,
	}

	err = minecraftClient.Auth("password")
	if err != nil {
		log.Fatalf("authentication failed: %s", err.Error())
	}

	result, err := minecraftClient.Exec("players")
	if err != nil {
		log.Fatalf("failed to retrive active players: %s", err.Error())
	}

	log.Println(string(result))
}

func ExampleMinecraftClient_Auth() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)

	minecraftClient := client.MinecraftClient{
		RemoteConsole: remoteConsole,
		IdGenFunc:     util.GenerateRequestId,
	}

	err = minecraftClient.Auth("password")
	if err != nil {
		log.Fatalf("authentication failed: %s", err.Error())
	}
}

func ExampleMinecraftClient_Exec() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatalf("connection failed: %s", err.Error())
	}
	defer conn.Close()

	remoteConsole := grcon.NewRemoteConsole(conn)

	minecraftClient := client.MinecraftClient{
		RemoteConsole: remoteConsole,
		IdGenFunc:     util.GenerateRequestId,
	}

	// before you execute something you might want to authenticated the connection.
	err = minecraftClient.Auth("password")
	if err != nil {
		log.Fatalf("authentication failed: %s", err.Error())
	}

	result, err := minecraftClient.Exec("players")
	if err != nil {
		log.Fatalf("failed to retrive active players: %s", err.Error())
	}

	log.Println(string(result))
}
