package models

// message structure
type Message struct {
	// the json tag means this will serialize as a lowercased field
	Type  string `json:"type"`
	Value string `json:"value"`
}
