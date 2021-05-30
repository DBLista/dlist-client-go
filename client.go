// With this library you can receive events from dlist.top gateway
// At the moment, the available events are: VOTE, RATE
// This library is maintained by dlist.top developers, so it's official
// Libraries for other languages can be found in the readme file

package dlist

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/BOOMfinity-Developers/wshelper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type VoteHandler func(data VoteData)
type RateHandler func(data RateData)

type Client struct {
	token  string
	ws     *wshelper.Connection
	logger zerolog.Logger
	err    chan error
	conn   chan bool

	onVoteHandler VoteHandler
	onRateHandler RateHandler
}

// Connect creates a connection to the dlist.top gateway.
func (c *Client) Connect(ctx context.Context) (err error) {
	c.err = make(chan error)
	c.conn = make(chan bool)
	c.ws, err = wshelper.Dial(ctx, GatewayURL, &websocket.DialOptions{})
	if err != nil {
		return
	}
	c.ws.OnMessage(c.onMessage)
	c.ws.OnError(c.onError)
	c.ws.OnClose(c.onClose)
	select {
	case err = <-c.err:
		return
	case <-c.conn:
		break
	case <-ctx.Done():
		return ctx.Err()
	}
	c.logger.Info().Msg("Successfully connected to the dlist.top gateway")
	return
}

func (c *Client) onError(_ *wshelper.Connection, err error) {
	c.logger.Error().Err(err).Send()
}

func (c *Client) onClose(_ *wshelper.Connection, code websocket.StatusCode, reason string) {
	c.logger.Warn().Int("code", int(code)).Str("reason", reason).Msg("connection has been closed")
}

func (c *Client) onMessage(conn *wshelper.Connection, _ websocket.MessageType, data wshelper.Payload) {
	var payload gatewayPayload
	err := data.Into(&payload)
	if err != nil {
		c.logger.WithLevel(zerolog.ErrorLevel).Err(err).Msg("error unmarshalling gateway op")
		return
	}
	c.logger.Debug().Str("OP", payload.Op.String()).Msg(string(payload.Data))
	switch payload.Op {
	case HelloOP:
		if err := conn.WriteJSON(context.Background(), sendData{
			Op: IdentifyOP,
			Data: identify{
				Token: c.token,
			},
		}); err != nil {
			c.logger.WithLevel(zerolog.ErrorLevel).Err(err).Send()
			return
		}
	case ReadyOP:
		c.conn <- true
	case EventOP:
		switch payload.Event {
		case VoteEvent:
			var voteData VoteData
			if err := json.Unmarshal(payload.Data, &voteData); err != nil {
				c.logger.Error().Err(err).Send()
				return
			}
			if c.onVoteHandler != nil {
				c.onVoteHandler(voteData)
			}
		case RateEvent:
			var rateData RateData
			if err := json.Unmarshal(payload.Data, &rateData); err != nil {
				c.logger.Error().Err(err).Send()
				return
			}
			if c.onRateHandler != nil {
				c.onRateHandler(rateData)
			}
		}
	case DisconnectOP:
		c.err <- errors.New(string(payload.Data))
	}
}

// OnVote will be executed if someone votes for the entity assigned to this token.
func (c *Client) OnVote(handler VoteHandler) {
	c.onVoteHandler = handler
}

// OnRate will be executed if someone adds review about the entity assigned to this token.
func (c *Client) OnRate(handler RateHandler) {
	c.onRateHandler = handler
}

// NewClientWithLogger works the same as NewClient, BUT additionally sets the logger with the specified level.
func NewClientWithLogger(token string, level zerolog.Level) *Client {
	logger := log.Level(level)
	c := NewClient(token)
	c.logger = logger
	return c
}

// NewClient returns a new instance of Client type.
func NewClient(token string) *Client {
	c := new(Client)
	c.token = token
	c.logger = log.Level(zerolog.InfoLevel)
	return c
}
