package models

import (
	"math"

	"golang.org/x/net/websocket"
)

// Point is a location on a two-axis plane
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p Point) Distance(p2 Point) float64 {
	first := math.Pow(float64(p2.X-p.X), 2)
	second := math.Pow(float64(p2.Y-p.Y), 2)
	return math.Sqrt(first + second)
}

// Player is a user
type Player struct {
	Point        Point
	WsConnection *websocket.Conn
	GUID         string `json:"GUID"`
	Name         string `json:"Name"`
	Health       int
}

// Projectile is a moving object
type Projectile struct {
	Point  Point
	XSpeed float64
	YSpeed float64
}

// Shield is a static object
type Shield struct {
	Point  Point
	Health int
}
