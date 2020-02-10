package main

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"os"
)

type mediaplayer2 struct {
	conn     *dbus.Conn
	properties *prop.Properties
}

func (m mediaplayer2) Raise() *dbus.Error {
	println("fuck this")
	notify("raise", "body")
	enc := newEncoder(os.Stdout)

	out := &message {
		Type: "prev",
	}

	if err := enc.Encode(out); err != nil {
		notify("error", err.Error())
	}
	return nil
}

func (m mediaplayer2) Quit() *dbus.Error {
	println("quit")
	return nil
}


var MP2 *mediaplayer2

func ExportMP2(conn *dbus.Conn) introspect.Interface {

	propSpec := map[string]*prop.Prop{
		"CanQuit" : { true, false, 1, nil },
		"CanRaise" : { true, false, 1, nil },
		"HasTrackList" : { false, false, 1, nil },
		"Identity" : { "Nightly", false, 1, nil },
		"SupportedUriSchemes" : { []string{"http"}, false, 1, nil },
		"SupportedMimeTypes" : { []string{"application/octet-stream", "audio/mpeg"}, false, 1, nil },
	}

	props, _ := prop.Export(conn, s_path, map[string]map[string]*prop.Prop{
		s_mp2: propSpec,
	})

	MP2 = &mediaplayer2 {
		conn,
		props,
	}

	conn.Export(MP2, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2")

	return introspect.Interface{
		Name: s_mp2,
		Methods: introspect.Methods(MP2),
		Properties: MP2.properties.Introspection(s_mp2),
	}
}
