package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"world/models"

	"golang.org/x/net/websocket"

	"github.com/labstack/echo"
)

var points []models.Point
var players []models.Player

func main() {

	go initializeWebsockets()

	go initializeHost()

	go initializeWorld()

	for {
		<-time.After(50 * time.Millisecond)
		go worldTick()
	}
}

func initializeWebsockets() {
	http.Handle("/ws", websocket.Handler(wsConnectionHandler))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("Websocket ListenAndServe: " + err.Error())
	}
}

func wsConnectionHandler(ws *websocket.Conn) {
	player := parsePlayerMessage(ws)
	player.WsConnection = ws
	players = append(players, player)
	for {
		player = parsePlayerMessage(ws)
		updatePlayer(player.Point.GUID, player.Point)
	}
}

func parsePlayerMessage(ws *websocket.Conn) models.Player {
	var messageString string
	if err := websocket.Message.Receive(ws, &messageString); err != nil {
		panic(err)
	}
	var message models.Message
	messageStringBytes := []byte(messageString)
	if err := json.Unmarshal(messageStringBytes, &message); err != nil {
		panic(err)
	}
	var point models.Point
	pointByte := []byte(message.Value)
	if err := json.Unmarshal(pointByte, &point); err != nil {
		panic(err)
	}
	var player models.Player
	player.Point = point
	return player
}

func updatePlayer(GUID string, point models.Point) {
	for i := 0; i < len(players); i++ {
		player := &players[i]
		if player.Point.GUID == GUID {
			player.Point.X = point.X
			player.Point.Y = point.Y
			player.Point.Name = point.Name
		}
	}
}

func initializeHost() {
	e := echo.New()
	e.Static("/", "mpsandbox-web/dist")
	e.Logger.Fatal(e.Start(":8080"))
}

func initializeWorld() {
	var point1 models.Point
	point1.X = 200
	point1.Y = 400
	var point2 models.Point
	point2.X = 750
	point2.Y = 500
	points = append(points, point1)
	points = append(points, point2)
}

func worldTick() {
	for i := 0; i < len(points); i++ {
		point := &points[i]
		point.X += (rand.Intn(10) - 5)
		point.Y += (rand.Intn(10) - 5)
	}
	var temp []models.Point
	temp = append(temp, points...)
	for i := 0; i < len(players); i++ {
		temp = append(temp, players[i].Point)
	}
	pointJSON, err := json.Marshal(temp)
	if err != nil {
		fmt.Println(err)
		return
	}
	pointJSONString := string(pointJSON)
	var toRemove []int
	for i, p := range players {
		m := models.Message{"pointsUpdate", pointJSONString}
		if err := websocket.JSON.Send(p.WsConnection, m); err != nil {
			log.Println(err)
			toRemove = append(toRemove, i)
			break
		}
	}
	cleanUpWsConnections(toRemove)
}

func reverse(j []int) {
	last := len(j) - 1
	for i := 0; i < len(j)/2; i++ {
		j[i], j[last-i] = j[last-i], j[i]
	}
}

func cleanUpWsConnections(toRemove []int) {
	reverse(toRemove)
	for _, removeIndex := range toRemove {
		players = append(players[:removeIndex], players[removeIndex+1:]...)
	}
}
