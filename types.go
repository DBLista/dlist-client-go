package dlist

import (
	"encoding/json"
	"strconv"
	"time"
)

type GatewayOP uint64

const (
	HelloOP      GatewayOP = iota + 1 // client <- gateway
	IdentifyOP                        // client -> gateway
	ReadyOP                           // client <- gateway
	DisconnectOP                      // client <- gateway
	EventOP                           // client <- gateway
)

func (o GatewayOP) String() string {
	return []string{"-", "HELLO", "IDENTIFY", "READY", "DISCONNECT", "EVENT"}[o]
}

type GatewayEvent string

const (
	VoteEvent GatewayEvent = "VOTE"
	RateEvent GatewayEvent = "RATE"
)

type gatewayPayload struct {
	Op    GatewayOP       `json:"op"`
	Data  json.RawMessage `json:"data"`
	Event GatewayEvent    `json:"event,omitempty"`
}

type sendData struct {
	Op   GatewayOP   `json:"op"`
	Data interface{} `json:"data"`
}

type identify struct {
	Token string `json:"token"`
}

type EntityType string

const (
	BotType    EntityType = "bots"
	ServerType EntityType = "servers"
)

type VoteData struct {
	UserID     string     `json:"authorID"`
	EntityType EntityType `json:"entityType"`
	EntityID   string     `json:"entityID"`
	Date       Timestamp  `json:"date"`
	// Total votes of entity NOT only this user
	TotalVotes uint64 `json:"totalVotes"`
}

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(bytes []byte) error {
	date, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return err
	}
	*t = Timestamp{
		Time: time.Unix(date/1000, 0),
	}
	return nil
}

type RateData struct {
	// number of "stars"
	Rating uint8 `json:"rating"`
	// message that the user added to the review
	Details    string     `json:"details"`
	UserID     string     `json:"authorID"`
	EntityType EntityType `json:"entityType"`
	EntityID   string     `json:"entityID"`
	Date       Timestamp  `json:"date"`
}
