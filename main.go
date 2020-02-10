package main

import (
	"encoding/json"
	"fmt"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"io"
	"os"
	"os/exec"
	"strconv"
)

const s_mp2 = "org.mpris.MediaPlayer2"
const s_player = "org.mpris.MediaPlayer2.Player"
const s_name = "org.mpris.MediaPlayer2.gotoaster"
const s_path = "/org/mpris/MediaPlayer2"

type message struct {
	Type string `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

func handleMessage(msg *message) {
	switch msg.Type {
	case "soundbadge": {
		song := &struct {
			Artist   string   `json:"artist"`
			Title    string   `json:"song"`
			ArtUrl   string   `json:"icon_url"`
			Time     string   `json:"time"`
			Length   string   `json:"length"`
			Location string   `json:"song_href"`
		}{}

		_ = json.Unmarshal(msg.Data, &song)

		length, _ := strconv.ParseInt(song.Length, 10, 64)

		metadata := map[string]dbus.Variant{
			"mpris:trackid": dbus.MakeVariant(song.Location),
			"xesam:url": dbus.MakeVariant("https://soundcloud.com" + song.Location),
			"xesam:title": dbus.MakeVariant(song.Title),
			"xesam:artist": dbus.MakeVariant([]string{song.Artist}),
			"mpris:artUrl": dbus.MakeVariant(song.ArtUrl),
			"mpris:length": dbus.MakeVariant(length),
		}
		Player.properties.SetMust(s_player, "Metadata", metadata)
	}
	case "playctrl": {
		playPause := &struct {
			Status string   `json:"status"`
		}{}

		_ = json.Unmarshal(msg.Data, &playPause)
		Player.properties.SetMust(s_player, "PlaybackStatus", playPause.Status)
	}
	default : {
		notify("Error", "unsupported message")
	}
	}
}

func readPump() {
	dec := newDecoder(os.Stdin)
	for {
		msg := &message{}
		if err := dec.Decode(msg); err != nil {
			if err == io.EOF {
				return
			}
			notify("err", err.Error())
			continue
		}

		handleMessage(msg)
	}
}

func notify(title string, body string) {
	cmd := exec.Command("notify-send", title, body)
	_ = cmd.Run()
}

func Export(conn *dbus.Conn) {
	mp2 := ExportMP2(conn)
	player := ExportPlayer(conn)

	node := &introspect.Node{
		Name: s_path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			mp2,
			player,
		},
	}
	_ = conn.Export(introspect.NewIntrospectable(node), s_path,
		"org.freedesktop.DBus.Introspectable")
}

func driveBus(conn *dbus.Conn) {
	Export(conn)

	reply, err := conn.RequestName(s_name, dbus.NameFlagReplaceExisting)
	if err != nil { panic(err) }

	if reply != dbus.RequestNameReplyPrimaryOwner {
		_, _ = fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	go driveBus(conn)
	readPump()

	select{}
}
