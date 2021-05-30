[![Go Reference](https://pkg.go.dev/badge/github.com/dlist-top/client-go.svg)](https://pkg.go.dev/github.com/dlist-top/client-go)

# DList.top Go client

Official [dlist.top](https://dlist.top) gateway client for Go.

## Installation

`go get github.com/dlist-top/client-go`

## Setup

To get your token please refer to [Gateway Docs](https://github.com/dlist-top/docs/wiki/Gateway)

## Example code

```go
package main

import (
	"context"
	"log"

	dlist "github.com/dlist-top/client-go"
)

func main() {
	c := dlist.NewClient("YOUR_API_TOKEN")
	if err := c.Connect(context.Background()); err != nil {
		panic(err)
	}

	c.OnVote(func(data dlist.VoteData) {
		log.Printf("%v voted for our bot / server. Total: %v", data.UserID, data.TotalVotes)
	})

	<-make(chan bool)
}

```

Notice: You can have up to 2 connections (per token) at the same time.
