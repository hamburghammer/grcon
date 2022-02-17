package grcon_test

import (
	"bytes"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/hamburghammer/grcon"
)

func TestRemoteConsole_Write(t *testing.T) {
	t.Run("compare written bytes", func(t *testing.T) {
		mockConn := &MockConn{}
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		expect := []byte{
			// size
			13, 0, 0, 0,
			// id
			1, 0, 0, 0,
			// type
			2, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,
		}

		// under test
		err := remoteConsole.Write(grcon.Packet{
			Id:   grcon.PacketId(1),
			Type: grcon.SERVERDATA_EXECCOMMAND,
			Body: []byte("foo"),
		})
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}
		got := mockConn.Send[0]

		if !bytes.Equal(expect, got) {
			t.Errorf("written bytes does not match:\nexpected:\n%b\ngot:\n%b\n", expect, got)
		}
	})

	t.Run("too long body", func(t *testing.T) {
		mockConn := &MockConn{}
		remoteConsole := grcon.NewRemoteConsole(mockConn)
		// under test
		err := remoteConsole.Write(grcon.Packet{
			Id:   grcon.PacketId(1),
			Type: grcon.SERVERDATA_EXECCOMMAND,
			Body: make([]byte, grcon.MaxPacket),
		})

		if _, ok := err.(grcon.RequestTooLongError); !ok {
			t.Errorf("error did not match:\nexpected:\n%T\ngot:\n%T", grcon.RequestTooLongError{}, err)
			t.FailNow()
		}
	})
}

func TestRemoteConsole_Read(t *testing.T) {
	t.Run("normal packet", func(t *testing.T) {
		mockConn := &MockConn{}
		mockConn.Receive = make([][]byte, 0, 1)
		mockConn.Receive = append(mockConn.Receive, []byte{
			// size
			13, 0, 0, 0,
			// id
			1, 0, 0, 0,
			// type
			0, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,
		})
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		expect := grcon.Packet{
			Id:   1,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: []byte("foo"),
		}
		// under test
		got, err := remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if !EqualPacket(expect, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect, got)
		}
	})

	t.Run("max packet", func(t *testing.T) {
		mockConn := &MockConn{}
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		// we use the write method to make it easier for us
		// for this we have to make sure that the method works as expected!
		remoteConsole.Write(grcon.Packet{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: make([]byte, grcon.MaxBody)})
		mockConn.Receive = make([][]byte, 0, 1)
		mockConn.Receive = append(mockConn.Receive, mockConn.Send...)

		expect := grcon.Packet{
			Id:   1,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: make([]byte, grcon.MaxBody),
		}
		// under test
		got, err := remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if len(got.Body) != len(expect.Body) {
			t.Errorf("packet body has not the expected size:\nexpected:%d\ngot: %d", len(expect.Body), len(got.Body))
		}

		if !EqualPacket(expect, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect, got)
		}
	})

	t.Run("min packet", func(t *testing.T) {
		mockConn := &MockConn{}
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		// we use the write method to make it easier for us
		// for this we have to make sure that the method works as expected!
		remoteConsole.Write(grcon.Packet{Id: 1, Type: grcon.SERVERDATA_RESPONSE_VALUE, Body: []byte{}})
		mockConn.Receive = make([][]byte, 0, 1)
		mockConn.Receive = append(mockConn.Receive, mockConn.Send...)

		expect := grcon.Packet{
			Id:   1,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: []byte{},
		}
		// under test
		got, err := remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if len(got.Body) != len(expect.Body) {
			t.Errorf("packet body has not the expected size:\nexpected:%d\ngot: %d", len(expect.Body), len(got.Body))
		}

		if !EqualPacket(expect, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect, got)
		}
	})

	t.Run("tow packets using the queue buffer.", func(t *testing.T) {
		mockConn := &MockConn{}
		mockConn.Receive = make([][]byte, 0, 1)
		mockConn.Receive = append(mockConn.Receive, []byte{
			// size
			13, 0, 0, 0,
			// id
			1, 0, 0, 0,
			// type
			0, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,

			// size
			13, 0, 0, 0,
			// id
			2, 0, 0, 0,
			// type
			0, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,
		})
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		expect1 := grcon.Packet{
			Id:   1,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: []byte("foo"),
		}
		expect2 := grcon.Packet{
			Id:   2,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: []byte("foo"),
		}

		// under test
		got, err := remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if !EqualPacket(expect1, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect1, got)
		}

		// second read
		got, err = remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if !EqualPacket(expect2, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect2, got)
		}

	})

	t.Run("receive size field over slow connection", func(t *testing.T) {
		mockConn := &MockConn{}
		mockConn.Receive = make([][]byte, 0, 5)
		// size
		mockConn.Receive = append(mockConn.Receive, []byte{13})
		mockConn.Receive = append(mockConn.Receive, []byte{0})
		mockConn.Receive = append(mockConn.Receive, []byte{0})
		mockConn.Receive = append(mockConn.Receive, []byte{0})

		mockConn.Receive = append(mockConn.Receive, []byte{
			// id
			1, 0, 0, 0,
			// type
			0, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,
		})

		remoteConsole := grcon.NewRemoteConsole(mockConn)

		expect := grcon.Packet{
			Id:   1,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: []byte("foo"),
		}
		// under test
		got, err := remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if !EqualPacket(expect, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect, got)
		}
	})

	t.Run("receive packet over slow connection execpt size field", func(t *testing.T) {
		mockConn := &MockConn{}
		mockConn.Receive = make([][]byte, 0, 5)
		// size
		mockConn.Receive = append(mockConn.Receive, []byte{13, 0, 0, 0})
		// id
		mockConn.Receive = append(mockConn.Receive, []byte{1, 0, 0, 0})
		// type
		mockConn.Receive = append(mockConn.Receive, []byte{0, 0, 0, 0})
		// body with null termination
		mockConn.Receive = append(mockConn.Receive, []byte{102, 111, 111, 0})
		// termination
		mockConn.Receive = append(mockConn.Receive, []byte{0})

		remoteConsole := grcon.NewRemoteConsole(mockConn)

		expect := grcon.Packet{
			Id:   1,
			Type: grcon.SERVERDATA_RESPONSE_VALUE,
			Body: []byte("foo"),
		}
		// under test
		got, err := remoteConsole.Read()
		if err != nil {
			t.Errorf("an error occurred that was not expected: %s", err.Error())
			t.FailNow()
		}

		if !EqualPacket(expect, got) {
			t.Errorf("packet are not equal:\nexpected:\n%+v\ngot:\n%+v", expect, got)
		}
	})

	t.Run("too large packet", func(t *testing.T) {
		mockConn := &MockConn{}
		mockConn.Receive = make([][]byte, 0, 1)
		mockConn.Receive = append(mockConn.Receive, []byte{
			// size
			// 4100
			4, 16, 0, 0,
			// id
			1, 0, 0, 0,
			// type
			0, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,
		})
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		// under test
		_, err := remoteConsole.Read()
		if _, ok := err.(grcon.ResponseTooLongError); !ok {
			t.Errorf("error did not match:\nexpected:\n%T\ngot:\n%T", grcon.ResponseTooLongError{}, err)
			t.FailNow()
		}
	})

	t.Run("too small packet", func(t *testing.T) {
		mockConn := &MockConn{}
		mockConn.Receive = make([][]byte, 0, 1)
		mockConn.Receive = append(mockConn.Receive, []byte{
			// size
			// 1
			1, 0, 0, 0,
			// id
			1, 0, 0, 0,
			// type
			0, 0, 0, 0,
			// body with null termination
			102, 111, 111, 0,
			// termination
			0,
		})
		remoteConsole := grcon.NewRemoteConsole(mockConn)

		// under test
		_, err := remoteConsole.Read()
		if _, ok := err.(grcon.UnexpectedFormatError); !ok {
			t.Errorf("error did not match:\nexpected:\n%T\ngot:\n%T", grcon.UnexpectedFormatError{}, err)
			t.FailNow()
		}
	})
}

// Helper functions

func EqualPacket(expected, got grcon.Packet) bool {
	if expected.Id != got.Id {
		return false
	}
	if expected.Type != got.Type {
		return false
	}
	if !bytes.Equal(expected.Body, got.Body) {
		return false
	}

	return true
}

// Mock implementations for tests

var ErrNotImplemented = errors.New("not implemented")

type MockConn struct {
	Send        [][]byte
	Receive     [][]byte
	nextSend    int
	nextReceive int
	IsClosed    bool
}

func (c *MockConn) Read(b []byte) (n int, err error) {
	toReceive := c.Receive[c.nextReceive]
	n = copy(b, toReceive)
	c.nextReceive++
	return
}

func (c *MockConn) Write(b []byte) (n int, err error) {
	if c.Send == nil {
		c.Send = make([][]byte, 0, 1)
	}
	toSend := make([]byte, len(b))
	n = copy(toSend, b)
	c.Send = append(c.Send, toSend)
	c.nextSend++
	return
}

func (c *MockConn) Close() error {
	c.IsClosed = true
	return nil
}

func (c *MockConn) LocalAddr() net.Addr {
	return mockAddr{}
}

func (c *MockConn) RemoteAddr() net.Addr {
	return mockAddr{}
}

func (c *MockConn) SetDeadline(t time.Time) error {
	return ErrNotImplemented
}

func (c *MockConn) SetReadDeadline(t time.Time) error {
	return ErrNotImplemented
}

func (c *MockConn) SetWriteDeadline(t time.Time) error {
	return ErrNotImplemented
}

type mockAddr struct{}

func (mockAddr) Network() string { return "mock addr" }
func (mockAddr) String() string  { return "mock addr" }
