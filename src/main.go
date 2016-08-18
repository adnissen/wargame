package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/adnissen/wargame/src/packages/army"
	"github.com/adnissen/wargame/src/packages/gameclient"
	"github.com/adnissen/wargame/src/packages/gamemap"
	"github.com/adnissen/wargame/src/packages/gamestat"
	"github.com/adnissen/wargame/src/packages/units"
	"github.com/adnissen/wargame/src/packages/userpkg"

	"github.com/adnissen/websocket"

	"github.com/adnissen/gorm"
	_ "github.com/adnissen/gorm/dialects/postgres"
)

var db *gorm.DB

var addr = flag.String("addr", "localhost:8080", "http service address")

var clients = make([]*gameclient.GameClient, 100) // represents the max number of players this server will run
var mmQueue []*gameclient.GameClient
var runningGames = make(map[string]*gamestat.GameStat)

type ClientMessage struct {
	MessageType string
	Message     string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func addPlayerToMMQueue(c *gameclient.GameClient) {
	mmQueue = append(mmQueue, c)
}

func removePlayerFromMMQueue(c *gameclient.GameClient) {
	for i := range mmQueue {
		if mmQueue[i] == c {
			mmQueue = append(mmQueue[:i], mmQueue[i+1:]...)
			return
		}
	}
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

func findMatches(c *gameclient.GameClient) {
	if len(mmQueue) != 0 && mmQueue[0] != nil {
		g := gamestat.CreateGame(db, mmQueue[0], c)
		runningGames[g.Uid.String()] = g
		removePlayerFromMMQueue(mmQueue[0])
	} else {
		addPlayerToMMQueue(c)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	log.Printf("client connected %s", r.RemoteAddr)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
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

	newClient.SendMessageOfType("announce", []byte("Welcome to Elder Runes!"))

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func() {
		removePlayerFromMMQueue(newClient)
		removeConnFromClients(c)
		c.Close()
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			break
		}

		if string(message) == "find_game" {
			findMatches(newClient)
		}

		if string(message) == "exit_queue" {
			removePlayerFromMMQueue(newClient)
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

		fmt.Println(cm)

		if cm.MessageType == "login" {
			if newClient.LoggedIn() {
				//get the user back into the previous gameclient, for now just return
				newClient.SendMessageOfType("create_user_result", []byte("failure"))

				return
			}
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(cm.Message), &dat); err != nil {
				panic(err)
			}

			username := dat["username"].(string)
			pass := dat["password"].(string)

			record := userpkg.VerifyUser(db, username, pass)
			if record != nil {
				newClient.User = record
				newClient.SendMessageOfType("login_result", []byte("success"))
			} else {
				newClient.SendMessageOfType("login_result", []byte("failure"))
			}
		}

		if cm.MessageType == "create_user" {
			if newClient.LoggedIn() {
			}

			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(cm.Message), &dat); err != nil {
				panic(err)
			}

			username := dat["username"].(string)
			pass := dat["password"].(string)
			email := dat["email"].(string)

			record := userpkg.CreateUser(db, username, pass, email)
			if record != nil {
				newClient.User = record
				newClient.SendMessageOfType("create_user_result", []byte("success"))
			} else {
				newClient.SendMessageOfType("create_user_result", []byte("failure"))
			}
		}

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

			moves := dat["moves"].([]interface{})
			mvs := make([][]int, len(moves))
			for f := 0; f < len(mvs); f++ {
				mvs[f] = make([]int, 2)
			}

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

	http.HandleFunc("/", echo)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	db, _ = gorm.Open("postgres", "host=localhost user=gorm dbname=gorm sslmode=disable password=mypassword")

	//migrate the schema
	db.AutoMigrate(&userpkg.User{})
	db.AutoMigrate(&army.Army{})

	h := sha256.New()
	io.WriteString(h, "1234test32")
	s := h.Sum(nil)
	newu := userpkg.CreateUser(db, "adn", "a@a.com", hex.EncodeToString(s))

	if newu != nil {
		fmt.Println("created account!")
		newu.AddArmy(db, army.Army{SquadIds: []int{0, 2, 0}})
	}

	log.Fatal(http.ListenAndServe(*addr, nil))
}
