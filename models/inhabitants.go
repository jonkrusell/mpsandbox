package models

import "golang.org/x/net/websocket"

// Point is a location on a two-axis plane
type Point struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	GUID string `json:"GUID"`
	Name string `json:"Name"`
}

type Player struct {
	Point        Point
	WsConnection *websocket.Conn
}
