package main

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"os"
)

type player struct {
	conn     *dbus.Conn
	properties *prop.Properties
}

func (p player) Next() *dbus.Error {
	enc := newEncoder(os.Stdout)

	out := &message {
		Type: "next",
	}

	if err := enc.Encode(out); err != nil {
		notify("error", err.Error())
	}

	return nil
}

func (p player) Previous() *dbus.Error {
	enc := newEncoder(os.Stdout)

	out := &message {
		Type: "prev",
	}

	if err := enc.Encode(out); err != nil {
		notify("error", err.Error())
	}
	return nil
}

func (p player) Pause() *dbus.Error {
	println("Pause")
	notify("Pause", "noice")
	return nil
}

func (p player) PlayPause() *dbus.Error {
	enc := newEncoder(os.Stdout)

	out := &message {
		Type: "playpause",
	}

	if err := enc.Encode(out); err != nil {
		notify("error", err.Error())
	}

	return nil
}

func (p player) Stop() *dbus.Error {
	println("Stop")
	return nil
}

func (p player) Play() *dbus.Error {
	println("Play")
	notify("Play", "noice")
	return nil
}

func (p player) Seek(Offset int64) *dbus.Error {
	println("Seek x: ", Offset)
	return nil
}

func (p player) SetPosition(TrackId dbus.ObjectPath, Position int64) *dbus.Error {
	println("SetPosition o: ", TrackId, " x: ", Position)
	return nil
}

func (p player) OpenUri(Uri string) *dbus.Error {
	println("OpenUri s: ", Uri)
	return nil
}

var Player *player


func ExportPlayer(conn *dbus.Conn) introspect.Interface {
	propSpec := map[string]*prop.Prop{
		"PlaybackStatus": { "Stopped", false, 1, nil },
		"LoopStatus": { "None", true, 1, nil },
		"Rate": {1.0, true, 1, nil},
		"Shuffle": { false, true, 1, nil },
		"Metadata": { map[string]dbus.Variant {
			// metadata, not yet existent
		}, false, 1, nil},
		"Volume": {1.0, true, 1, nil},
		"Position": {int64(0), false, 0, nil},
		"MinimumRate": {1.0, false, 1, nil},
		"MaximumRate": {1.0, false, 1, nil},
		"CanGoNext": {true, false, 1, nil},
		"CanGoPrevious": {true, false, 1, nil},
		"CanPlay": {true, false, 1, nil},
		"CanPause": {true, false, 1, nil},
		"CanSeek": {false, false, 1, nil},
		"CanControl": {true, false, 0, nil},
	}

	props, _ := prop.Export(conn, s_path, map[string]map[string]*prop.Prop{
		s_player: propSpec,
	})

	Player = &player{
		conn,
		props,
	}

	_ = conn.Export(Player, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player")

	return introspect.Interface {
		Name: s_player,
		Signals: []introspect.Signal{{
			Name:        "Seeked",
			Args:        []introspect.Arg{{
				Name:      "Position",
				Type:      "x",
				Direction: "out",
			}},
			Annotations: nil,
		}},
		Methods: introspect.Methods(Player),
		Properties: Player.properties.Introspection(s_player),
	}
}
