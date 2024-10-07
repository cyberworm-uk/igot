package main

import (
	"flag"
	"igot/twitch"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	var username = flag.String("username", os.Getenv("TWITCH_USERNAME"), "default username to login to twitch")
	var channels = flag.String("channels", os.Getenv("TWITCH_CHANNELS"), "comma delimited list of channels to serve")
	flag.Parse()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	var b = twitch.New(*username, strings.Split(*channels, ","))
	b.Start()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		log.Println("signal received, stopping...")
		b.Stop()
		wg.Done()
		log.Println("done...")
	}()
}
