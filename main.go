package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"mpsandbox/models"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/labstack/echo"
)

var players []models.Player
var projectiles []models.Projectile
var shields []models.Shield

func main() {

	go initializeWebsockets()

	go initializeHost()

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
	for {
		message := parseMessage(ws)
		if message.Type == "playerUpdate" {
			player := parsePlayerMessage(message.Value)
			player.WsConnection = ws
			players = append(players, player)
			break
		}
	}
	for {
		message := parseMessage(ws)
		switch message.Type {
		case "playerUpdate":
			player := parsePlayerMessage(message.Value)
			updatePlayer(player)
		case "playerShoot":
			projectile := parsePlayerShootMessage(message.Value)
			projectiles = append(projectiles, projectile)
		case "playerShield":
			shield := parsePlayerShieldMessage(message.Value)
			shields = append(shields, shield)
		}
	}
}

func parseMessage(ws *websocket.Conn) models.Message {
	var messageString string
	if err := websocket.Message.Receive(ws, &messageString); err != nil {
		panic(err)
	}
	var message models.Message
	messageStringBytes := []byte(messageString)
	if err := json.Unmarshal(messageStringBytes, &message); err != nil {
		panic(err)
	}
	return message
}

func parsePlayerMessage(messageValue string) models.Player {
	var player models.Player
	messageBytes := []byte(messageValue)
	if err := json.Unmarshal(messageBytes, &player); err != nil {
		panic(err)
	}
	return player
}

func parsePlayerShootMessage(messageValue string) models.Projectile {
	var playerShootMessage models.PlayerShootMessage
	messageBytes := []byte(messageValue)
	if err := json.Unmarshal(messageBytes, &playerShootMessage); err != nil {
		panic(err)
	}
	var projectile models.Projectile
	projectile.Point = playerShootMessage.FromPoint
	var yChange = playerShootMessage.ToPoint.Y - playerShootMessage.FromPoint.Y
	var xChange = playerShootMessage.ToPoint.X - playerShootMessage.FromPoint.X
	var theta = math.Atan2(float64(yChange), float64(xChange))
	projectile.XSpeed = math.Cos(theta) * float64(20)
	projectile.YSpeed = math.Sin(theta) * float64(20)
	projectile.Point.X += projectile.XSpeed * float64(2)
	projectile.Point.Y += projectile.YSpeed * float64(2)
	return projectile
}

func parsePlayerShieldMessage(messageValue string) models.Shield {
	var playerShieldMessage models.Point
	messageBytes := []byte(messageValue)
	if err := json.Unmarshal(messageBytes, &playerShieldMessage); err != nil {
		panic(err)
	}
	var shield models.Shield
	shield.Point = playerShieldMessage
	shield.Health = 10
	return shield
}

func updatePlayer(updatedPlayer models.Player) {
	for i := 0; i < len(players); i++ {
		player := &players[i]
		if player.GUID == updatedPlayer.GUID {
			player.Point.X = updatedPlayer.Point.X
			player.Point.Y = updatedPlayer.Point.Y
			player.Name = updatedPlayer.Name
		}
	}
}

func initializeHost() {
	e := echo.New()
	e.Static("/", "mpsandbox-web/dist")
	e.Logger.Fatal(e.Start(":8080"))
}

func worldTick() {
	playersJSON, err := json.Marshal(players)
	if err != nil {
		fmt.Println(err)
		return
	}
	playersJSONString := string(playersJSON)

	var projectilesToRemove []int
	var shieldsToRemove []int
	for i := 0; i < len(projectiles); i++ {
		projectiles[i].Point.X += projectiles[i].XSpeed
		projectiles[i].Point.Y += projectiles[i].YSpeed
		for j := 0; j < len(players); j++ {
			if players[j].Point.Distance(projectiles[i].Point) < 30 {
				projectilesToRemove = AppendIfMissing(projectilesToRemove, i)
				players[j].Health -= 10
			}
		}
		for j := 0; j < len(shields); j++ {
			if shields[j].Point.Distance(projectiles[i].Point) < 30 {
				projectilesToRemove = AppendIfMissing(projectilesToRemove, i)
				shields[j].Health -= 5
				if shields[j].Health <= 0 {
					shieldsToRemove = AppendIfMissing(shieldsToRemove, j)
				}
			}
		}
	}
	cleanUpShields(shieldsToRemove)
	shieldsJSON, err := json.Marshal(shields)
	if err != nil {
		fmt.Println(err)
		return
	}
	shieldsJSONString := string(shieldsJSON)

	cleanUpProjectiles(projectilesToRemove)
	projectilesJSON, err := json.Marshal(projectiles)
	if err != nil {
		fmt.Println(err)
		return
	}
	projectilesJSONString := string(projectilesJSON)

	var playersToRemove []int
	for i, p := range players {
		m := models.Message{"playersUpdate", playersJSONString}
		if err := websocket.JSON.Send(p.WsConnection, m); err != nil {
			log.Println(err)
			playersToRemove = append(playersToRemove, i)
			break
		}
		m = models.Message{"projectilesUpdate", projectilesJSONString}
		if err := websocket.JSON.Send(p.WsConnection, m); err != nil {
			log.Println(err)
			playersToRemove = append(playersToRemove, i)
			break
		}
		m = models.Message{"shieldsUpdate", shieldsJSONString}
		if err := websocket.JSON.Send(p.WsConnection, m); err != nil {
			log.Println(err)
			playersToRemove = append(playersToRemove, i)
			break
		}
	}
	cleanUpWsConnections(playersToRemove)
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

func cleanUpProjectiles(toRemove []int) {
	reverse(toRemove)
	for _, removeIndex := range toRemove {
		projectiles = append(projectiles[:removeIndex], projectiles[removeIndex+1:]...)
	}
}

func cleanUpShields(toRemove []int) {
	reverse(toRemove)
	for _, removeIndex := range toRemove {
		shields = append(shields[:removeIndex], shields[removeIndex+1:]...)
	}
}

func AppendIfMissing(slice []int, i int) []int {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}
