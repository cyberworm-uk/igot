package twitch

import (
	"fmt"
	"igot/runes"
	"log"
	"strconv"
	"strings"

	"github.com/adeithe/go-twitch/irc"
)

type Bot struct {
	conn     *irc.Conn
	channels []string
	username string
	oauth    string
	stopping bool
}

func (b *Bot) handleJoin(channel, username string) {
	log.Printf("JOIN: %q joined %q\n", username, channel)
}

func (b *Bot) handleJoin(channel, username string) {
	log.Printf("LEAVE: %q left %q\n", username, channel)
}

func (b *Bot) handleMsg(m irc.ChatMessage) {
	log.Printf("MESSAGE: %q said %q\n", m.IRCMessage.Sender.Nickname, m.IRCMessage.Text)
	if strings.HasPrefix(m.IRCMessage.Text, "!") {
		var args = strings.SplitN(m.IRCMessage.Text, " ", 2)
		switch args[0] {
		case "!igot":
			if len(args) > 1 {
				if r, e := strconv.Atoi(args[1]); e == nil {
					log.Printf("Calculating runes %d\n", r)
					var rl = runes.Calc(r)
					log.Printf("Rune Level %d\n", rl)
					if rl != 0 {
						if e := b.conn.Say(m.Channel, fmt.Sprintf("Rune Level: %d", rl)); e != nil {
							log.Println(e)
						}
					}
				} else {
					log.Println(e)
				}
			}
		}
	}
}

func (b *Bot) handleNotice(notice irc.ServerNotice) {
	log.Printf("NOTICE: %q %q\n", notice.Channel, notice.Message)
}

func (b *Bot) handleUserNotice(notice irc.UserNotice) {
	log.Printf("USERNOTICE: %q %q\n", notice.Sender.Username, notice.Message)
}

func (b *Bot) handleDisconnect() {
	if b.stopping {
		return
	}
	log.Printf("DISCONNECT!\n")
	b.refreshAuth()
	b.Start()
}

func (b *Bot) setHandlers() {
	b.conn.OnChannelJoin(b.handleJoin)
	b.conn.OnChannelLeave(b.handleLeave)
	b.conn.OnMessage(b.handleMsg)
	b.conn.OnServerNotice(b.handleNotice)
	b.conn.OnChannelUserNotice(b.handleUserNotice)
	b.conn.OnDisconnect(b.handleDisconnect)
}

func (b *Bot) refreshAuth() {
	b.oauth = accessToken()
	b.conn.SetLogin(b.username, b.oauth)
}

func (b *Bot) Start() {
	if e := b.conn.Connect(); e != nil {
		log.Printf("ERROR: Failed to connect: %q\n", e)
	}
	if e := b.conn.Join(b.channels...); e != nil {
		log.Printf("ERROR: Failed to join %q: %q\n", b.channels, e)
	}
}

func (b *Bot) Stop() {
	b.stopping = true
	if !b.conn.IsConnected() {
		return
	}
	if e := b.conn.Leave(b.channels...); e != nil {
		log.Printf("ERROR: Failed to leave %q: %q\n", b.channels, e)
	}
	b.conn.Close()
}

func New(username string, channels []string) *Bot {
	var b = &Bot{
		channels: make([]string, len(channels)),
		username: username,
		oauth:    accessToken(),
		conn:     &irc.Conn{},
	}
	copy(b.channels, channels)
	b.setHandlers()
	b.conn.SetLogin(b.username, b.oauth)
	return b
}
