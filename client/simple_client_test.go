package client_test

import (
	"testing"

	"github.com/hamburghammer/grcon"
	"github.com/hamburghammer/grcon/client"
)

func TestSimpleClient_Auth(t *testing.T) {
	t.Run("successfull auth", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
			{Id: 1, Type: grcon.SERVERDATA_AUTH_RESPONSE, Body: []byte("")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}
		err := simpleClient.Auth("foo")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("pre empty response", func(t *testing.T) {
		t.Run("not match type", func(t *testing.T) {
			mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
			mock := &MockRemoteConsole{In: []grcon.Packet{
				{Id: 1, Type: grcon.SERVERDATA_AUTH_RESPONSE, Body: []byte("")},
			}}
			simpleClient := client.SimpleClient{
				RemoteConsole: mock,
				IdGenFunc:     mockIdGen.GetNextId,
			}
			err := simpleClient.Auth("foo")

			_, ok := err.(client.InvalidResponseTypeError)
			if err != nil && !ok {
				t.Errorf("expected: InvalidResponseTypeError\ngot: %T\n", err)
				t.Error(err)
				t.FailNow()
			}
		})

		t.Run("not matching id", func(t *testing.T) {
			mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
			mock := &MockRemoteConsole{In: []grcon.Packet{
				{Id: 2, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
			}}
			simpleClient := client.SimpleClient{
				RemoteConsole: mock,
				IdGenFunc:     mockIdGen.GetNextId,
			}
			err := simpleClient.Auth("foo")

			_, ok := err.(client.ResponseIdMismatchError)
			if err != nil && !ok {
				t.Errorf("expected: ResponseIdMismatchError\ngot: %T\n", err)
				t.Error(err)
				t.FailNow()
			}
		})

		t.Run("not empty body", func(t *testing.T) {
			mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
			mock := &MockRemoteConsole{In: []grcon.Packet{
				{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("ooops")},
			}}
			simpleClient := client.SimpleClient{
				RemoteConsole: mock,
				IdGenFunc:     mockIdGen.GetNextId,
			}
			err := simpleClient.Auth("foo")

			_, ok := err.(client.ResponseBodyError)
			if err != nil && !ok {
				t.Errorf("expected: ResponseBodyError\ngot: %T\n", err)
				t.Error(err)
				t.FailNow()
			}
		})
	})

	t.Run("auth response", func(t *testing.T) {
		t.Run("not match type", func(t *testing.T) {
			mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
			mock := &MockRemoteConsole{In: []grcon.Packet{
				{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
				{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
			}}
			simpleClient := client.SimpleClient{
				RemoteConsole: mock,
				IdGenFunc:     mockIdGen.GetNextId,
			}
			err := simpleClient.Auth("foo")

			_, ok := err.(client.InvalidResponseTypeError)
			if err != nil && !ok {
				t.Errorf("expected: InvalidResponseTypeError\ngot: %T\n", err)
				t.Error(err)
				t.FailNow()
			}
		})

		t.Run("auth failed", func(t *testing.T) {
			mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
			mock := &MockRemoteConsole{In: []grcon.Packet{
				{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
				{Id: -1, Type: grcon.SERVERDATA_AUTH_RESPONSE, Body: []byte("")},
			}}
			simpleClient := client.SimpleClient{
				RemoteConsole: mock,
				IdGenFunc:     mockIdGen.GetNextId,
			}
			err := simpleClient.Auth("foo")

			_, ok := err.(client.AuthFailedError)
			if err != nil && !ok {
				t.Errorf("expected: AuthFailedError\ngot: %T\n", err)
				t.Error(err)
				t.FailNow()
			}
		})

		t.Run("not matching ids", func(t *testing.T) {
			mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
			mock := &MockRemoteConsole{In: []grcon.Packet{
				{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
				{Id: 2, Type: grcon.SERVERDATA_AUTH_RESPONSE, Body: []byte("")},
			}}
			simpleClient := client.SimpleClient{
				RemoteConsole: mock,
				IdGenFunc:     mockIdGen.GetNextId,
			}
			err := simpleClient.Auth("foo")

			_, ok := err.(client.ResponseIdMismatchError)
			if err != nil && !ok {
				t.Errorf("expected: ResponseIdMismatchError\ngot: %T\n", err)
				t.Error(err)
				t.FailNow()
			}
		})
	})

	t.Run("written packet", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
			{Id: 1, Type: grcon.SERVERDATA_AUTH_RESPONSE, Body: []byte("")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}
		err := simpleClient.Auth("foo")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		got := mock.Out
		if len(got) < 1 {
			t.Errorf("expected at leased 1 packet but got %d\n", len(got))
		}
		gotPacket := got[0]
		if gotPacket.Id != 1 {
			t.Error("expected id 1 but got something different")
		}
		if gotPacket.Type != grcon.SERVERDATA_AUTH {
			t.Error("packet was not of type auth")
		}
		if string(gotPacket.Body) != "foo" {
			t.Error("body/password did not match")
		}
	})
}

func TestSimpleClient_Exec(t *testing.T) {
	t.Run("successfull execution", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1, 2}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("bar")},
			{Id: 2, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}
		got, err := simpleClient.Exec("foo")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if string(got) != "bar" {
			t.Errorf("response did not match:\nexpected: %s\ngot: %s\n", "bar", string(got))
		}
	})

	t.Run("successfull multi packet execution", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1, 2}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("foo")},
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("bar")},
			{Id: 2, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}
		got, err := simpleClient.Exec("foo")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if string(got) != "foobar" {
			t.Errorf("response did not match:\nexpected: %s\ngot: %s\n", "bar", string(got))
		}
	})

	t.Run("write cmd packet", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1, 2}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
			{Id: 2, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}
		_, err := simpleClient.Exec("foo")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		got := mock.Out

		if len(got) < 1 {
			t.Error("written less than 1 packet")
			t.FailNow()
		}

		gotPacket := got[0]

		if gotPacket.Id != 1 {
			t.Error("id missmatch")
		}
		if gotPacket.Type != grcon.SERVERDATA_EXECCOMMAND {
			t.Error("type missmatch")
		}
		if string(gotPacket.Body) != "foo" {
			t.Error("body missmatch")
		}
	})

	t.Run("write delimiter packet", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1, 2}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
			{Id: 2, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}
		_, err := simpleClient.Exec("foo")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		got := mock.Out

		if len(got) < 2 {
			t.Error("written less than 2 packet")
			t.FailNow()
		}

		gotPacket := got[1]

		if gotPacket.Id != 2 {
			t.Error("id missmatch")
		}
		if gotPacket.Type != grcon.SERVERDATA_RESPONSE_VALUE {
			t.Error("type missmatch")
		}
		if string(gotPacket.Body) != "" {
			t.Error("body missmatch")
		}
	})

	t.Run("invalid response type error", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1, 2}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 1, Type: grcon.SERVERDATA_AUTH_RESPONSE, Body: []byte("bar")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}

		_, err := simpleClient.Exec("foo")
		_, ok := err.(client.InvalidResponseTypeError)
		if err != nil && !ok {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("id response missmatch type error", func(t *testing.T) {
		mockIdGen := &MockIdGenerator{Ids: []grcon.PacketId{1, 2}}
		mock := &MockRemoteConsole{In: []grcon.Packet{
			{Id: 3, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte("bar")},
		}}
		simpleClient := client.SimpleClient{
			RemoteConsole: mock,
			IdGenFunc:     mockIdGen.GetNextId,
		}

		_, err := simpleClient.Exec("foo")
		_, ok := err.(client.ResponseIdMismatchError)
		if err != nil && !ok {
			t.Error(err)
			t.FailNow()
		}
	})
}

// Mock

type MockIdGenerator struct {
	Ids    []grcon.PacketId
	nextId int
}

func (i *MockIdGenerator) GetNextId() grcon.PacketId {
	id := i.Ids[i.nextId]
	i.nextId++
	return id
}

type MockRemoteConsole struct {
	Out    []grcon.Packet
	In     []grcon.Packet
	nextIn int
}

func (m *MockRemoteConsole) Read() (grcon.Packet, error) {
	packet := m.In[m.nextIn]
	m.nextIn++

	return packet, nil
}

func (m *MockRemoteConsole) Write(packet grcon.Packet) error {
	if m.Out == nil {
		m.Out = make([]grcon.Packet, 0, 2)
	}
	m.Out = append(m.Out, packet)

	return nil
}
