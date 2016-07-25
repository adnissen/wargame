// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
		if clients[k] == nil {
			continue
		}
		if clients[k].CompareWebSocketConn(c) {
			if clients[k].CurrentGame.String() != "00000000-0000-0000-0000-000000000000" {
				runningGames[clients[k].CurrentGame.String()].PlayerDisconnect(clients[k])
			}
			clients[k] = nil
			return
		}
	}
}

func findMatches() {
	var c1 *gameclient.GameClient
	for i := range clients {
		if clients[i] != nil && clients[i].CurrentGame.String() == "00000000-0000-0000-0000-000000000000" {
			if c1 == nil {
				c1 = clients[i]
			} else {
				fmt.Println("found match")
				gstat := gamestat.CreateGame(c1, clients[i])
				runningGames[gstat.Uid.String()] = gstat
				break
			}
		}
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	log.Printf("client connected %s", r.RemoteAddr)
	newClient := insertConnIntoClients(c)

	/*
		  typeof("aimee") string
			var a = "aimee"
			b = &a
			typeof(b) *string

			reverse(a)
	*/

	newA := army.Army{Squads: []units.Squad{units.CreateSquad(0), units.CreateSquad(1), units.CreateSquad(0)}}

	newClient.Army = newA

	//information about the game so that the client can download it if need be
	//(we really don't want lag in loading images and whatnot once they're in the game)

	newClient.SendMessage(units.UnitInformation())
	newClient.SendMessage(units.SquadInformation())
	newClient.SendMessage(newA.ArmyInformation())

	findMatches()
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func() {
		removeConnFromClients(c)
		c.Close()
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			break
		}
		//log.Printf("game_logic: %s", message)
		if string(message) == "game" {
			newClient.SendMessage([]byte(newClient.CurrentGame.String()))
		}

		if string(message) == "end_turn" {
			if newClient.CurrentGame.String() != "00000000-0000-0000-0000-000000000000" {
				runningGames[newClient.CurrentGame.String()].EndTurn(newClient)
			}
		}

		cm := ClientMessage{}
		json.Unmarshal(message, &cm)

		if cm.MessageType == "map_export_data" {
			gamemap.InsertMap(gamemap.ImportMap(cm.Message))
		}

		if cm.MessageType == "game_use_weapon" {
			if newClient.CurrentGame.String() == "00000000-0000-0000-0000-000000000000" {
				return
			}
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(cm.Message), &dat); err != nil {
				panic(err)
			}

			g := runningGames[newClient.CurrentGame.String()]
			u := g.GetUnit(dat["uid"].(string), g.GetPlayerIndex(newClient))
			t := g.GetUnitGlobal(dat["target"].(string))
			w := g.GetWeapon(dat["weapon"].(string), g.GetPlayerIndex(newClient))
			used, damage, roll := g.UseWeapon(u, t, w, g.GetPlayerIndex(newClient))
			if used == true {
				ret := map[string]interface{}{
					"uid":    u.Uid.String(),
					"weapon": w.Uid.String(),
					"target": t.Uid.String(),
					"roll":   strconv.Itoa(roll),
					"damage": strconv.Itoa(damage)}
				str, err := json.Marshal(ret)
				if err != nil {
					fmt.Println("Error encoding JSON")
					return
				}
				g.SendMessageToAllPlayers("game_use_weapon", str)
			}
		}

		if cm.MessageType == "game_move" {
			if newClient.CurrentGame.String() == "00000000-0000-0000-0000-000000000000" {
				return
			}
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(cm.Message), &dat); err != nil {
				panic(err)
			}

			// "[[1, 2], [1, 3], [1, 4]]"
			moves := dat["moves"].([]interface{})
			// moves = ["[1, 2]", "[1, 3]", "[1, 4]"]
			mvs := make([][]int, len(moves))
			for f := 0; f < len(mvs); f++ {
				mvs[f] = make([]int, 2)
			}
			// mvs = [[], [], []]

			for i := range moves {
				t := moves[i].([]interface{})
				for k := range t {
					mvs[i][k] = int(t[k].(float64))
				}
			}
			g := runningGames[newClient.CurrentGame.String()]
			moved := g.MoveUnit(g.GetUnit(dat["uid"].(string), g.GetPlayerIndex(newClient)), mvs)
			fmt.Println(moved)

			if moved == true {
				g.SendMessageToAllPlayers("game_move", []byte(cm.Message))
			}
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

	units.LoadWeapons()
	log.Println("Loaded Weapons")
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
