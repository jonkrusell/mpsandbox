package models

// Message structure
type Message struct {
	// the json tag means this will serialize as a lowercased field
	Type  string `json:"type"`
	Value string `json:"value"`
}

// PlayerShoot message structure
type PlayerShootMessage struct {
	FromPoint Point `json:"fromPoint"`
	ToPoint   Point `json:"toPoint"`
}
