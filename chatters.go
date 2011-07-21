package main

import (
	markov "github.com/zeebo/markov"
	irc "github.com/zeebo/goirc"
	"time"
	"rand"
	"fmt"
	"math"
)

type Chatter struct {
	Messages chan string
	mark *markov.Markov
	average float64
	stddev float64
}

func (c* Chatter) SendChats() {
	for {
		duration := -1 * math.Log(rand.Float64()) * c.average * 1e9
		time.Sleep(int64(duration))
		c.Messages <- c.mark.Generate()
	}
}

func NewChatter(u* UserInfo) *Chatter {
	mark := markov.New()
	for _, line := range u.Lines {
		mark.Analyze(line.message)
	}
	u.Calculate()
	return &Chatter{
		mark: mark,
		average: u.MeanDelay,
		stddev: u.StdDelay,
		Messages: make(chan string),
	}
}

func main() {
	rand.Seed(time.Nanoseconds())

	infos, err := Analyze("log.txt")
	if err != nil {
		panic(err)
	}

	for name := range infos {
		c := NewChatter(infos[name])

		info := irc.Info{
			Channel: "#okcodev",
			Nick: fmt.Sprintf("%sbot", name),
			AltNick: fmt.Sprintf("_%sbot", name),
			Server: "zeeb.us.to:6667",
		}

		conn, err := irc.NewConnection(info)
		if err != nil {
			panic(err)
		}

		conn.SetupModcmd("zeebo")

		go c.SendChats()
		go func(con *irc.Connection) {
			for {
				fmt.Fprintln(con, <-c.Messages)
			}
		}(conn)
		go conn.Handle()
	}

	q := make(chan bool)
	<-q
}