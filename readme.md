[![Go Reference](https://pkg.go.dev/badge/github.com/dlist-top/client-go.svg)](https://pkg.go.dev/github.com/dlist-top/client-go)

Official dlist.top gateway client written in go.

To add to the project, run:

`go get github.com/dlist-top/client-go`

Example code

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

You can have up to 2 connections (per token) at the same time.
