// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/adnissen/wargame/src/packages/army"
	"github.com/adnissen/wargame/src/packages/gameclient"
	"github.com/adnissen/wargame/src/packages/gamemap"
	"github.com/adnissen/wargame/src/packages/gamestat"
	"github.com/adnissen/wargame/src/packages/units"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}               // use default options
var clients = make([]*gameclient.GameClient, 100) // represents the max number of players this server will run
var runningGames = make(map[string]*gamestat.GameStat)

type ClientMessage struct {
	MessageType string
	Message     string
}

func sendMessageToAllClients(m []byte) {
	for k, _ := range clients {
		if clients[k] == nil {
			clients[k].SendMessage(m)
		}
	}
}

func insertConnIntoClients(c *websocket.Conn) *gameclient.GameClient {
	var ret *gameclient.GameClient
	for k, _ := range clients {
		if clients[k] == nil {
			clients[k] = gameclient.CreateGameClient(c)
			ret = clients[k]
			return ret
		}
	}
	return ret
}

func removeConnFromClients(c *websocket.Conn) {
	for k, _ := range clients {
		if clients[k].CompareWebSocketConn(c) {
			if clients[k].CurrentGame.String() != "00000000-0000-0000-0000-000000000000" {
				runningGames[clients[k].CurrentGame.String()].EndGame()
			}
			clients[k] = nil
			return
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	log.Printf("client connected %s", r.RemoteAddr)
	newClient := insertConnIntoClients(c)
	newA := army.Army{Squads: []units.Squad{units.GetSquad(0), units.GetSquad(0), units.GetSquad(0)}}
	newClient.Army = newA

	//information about the game so that the client can download it if need be
	//(we really don't want lag in loading images and whatnot once they're in the game)

	newClient.SendMessage(units.UnitInformation())
	newClient.SendMessage(units.SquadInformation())
	newClient.SendMessage(newA.ArmyInformation())

	//VERY poor mans matchmaking, only works for the first two people
	if clients[0] != nil && clients[1] != nil {
		gstat := gamestat.CreateGame(clients[0], clients[1])
		runningGames[gstat.Uid.String()] = gstat
	}
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			removeConnFromClients(c)
			break
		}
		//log.Printf("game_logic: %s", message)
		if string(message) == "game" {
			newClient.SendMessage([]byte(newClient.CurrentGame.String()))
		}

		cm := ClientMessage{}
		json.Unmarshal(message, &cm)

		if cm.MessageType == "map_export_data" {
			gamemap.InsertMap(gamemap.ImportMap(cm.Message))
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	dat, _ := ioutil.ReadFile("./public/client/index.html")
	var homeTemplate = template.Must(template.New("").Parse(string(dat)))
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	units.LoadUnits()
	log.Println("Loaded Units")
	units.LoadSquads()
	log.Println("Loaded Squads")
	gamemap.LoadMaps()
	log.Println("Loaded Maps")

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	//http.Handle("/", http.FileServer(http.Dir("./public/client")))

	/*ticker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for t := range ticker.C {
			log.Println(t)
			sendMessageToAllClients([]byte("Tick"))
		}
	}()*/

	log.Fatal(http.ListenAndServe(*addr, nil))
}
