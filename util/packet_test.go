package util_test

import (
	"testing"

	"github.com/hamburghammer/grcon"
	"github.com/hamburghammer/grcon/util"
)

func TestNewExecCommandPacket(t *testing.T) {
	got := util.NewExecCommandPacket(1, "foo")

	if got.Type != grcon.SERVERDATA_EXECCOMMAND {
		t.Errorf("type error:\nexpected: %v\ngot: %v\n", grcon.SERVERDATA_EXECCOMMAND, got.Type)
	}
	if got.Id != 1 {
		t.Errorf("type error:\nexpected: %d\ngot: %d\n", 1, got.Id)
	}
	if string(got.Body) != "foo" {
		t.Errorf("body error:\nexpected: %s\ngot: %s\n", "foo", string(got.Body))
	}
}

func TestNewEmptyResponseValuePacket(t *testing.T) {
	got := util.NewEmptyResponseValuePacket(1)

	if got.Type != grcon.SERVERDATA_RESPONSE_VALUE {
		t.Errorf("type error:\nexpected: %v\ngot: %v", grcon.SERVERDATA_RESPONSE_VALUE, got.Type)
	}
	if got.Id != 1 {
		t.Errorf("type error:\nexpected: %d\ngot: %d\n", 1, got.Id)
	}
	if string(got.Body) != "" {
		t.Errorf("body error:\nexpected: %s\ngot: %s", "", string(got.Body))
	}
}
