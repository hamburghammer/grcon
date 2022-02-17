package util_test

import (
	"testing"

	"github.com/hamburghammer/grcon"
	"github.com/hamburghammer/grcon/util"
)

func TestRemoteConsole(t *testing.T) {
	t.Run("*rcon.RemoteConsole should implement interface", func(t *testing.T) {
		var r interface{} = &grcon.RemoteConsole{}
		_, ok := r.(util.RemoteConsole)
		if !ok {
			t.Error("*rcon.RemoteConsole does not implement util.RemoteConsole interface")
		}
	})
}
