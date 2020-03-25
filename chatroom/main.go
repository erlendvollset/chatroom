package main

import "flag"
import "fmt"
import "github.com/erlendvollset/chatroom/client"
import "github.com/erlendvollset/chatroom/server"
import "os"
import "strings"

var userName *string
var port *string
var host *string

func parseClientCommand() {
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)
	userName = clientCommand.String("n", "Anon", "The name of the user")
	host = clientCommand.String("h", "localhost", "The host for the chatroom server")
	port = clientCommand.String("p", "8080", "The port of the chatroom server")
	clientCommand.Parse(os.Args[2:])
}

func parseServerCommand() {
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	host = serverCommand.String("h", "localhost", "The host for the chatroom server")
	port = serverCommand.String("p", "8080", "The port of the chatroom server")
	serverCommand.Parse(os.Args[2:])
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: chatroom {client|server} [args]")
		os.Exit(0)
	}
	switch strings.ToLower(os.Args[1]) {
	case "client":
		parseClientCommand()
		client.New(userName).Start(*host, *port)
	case "server":
		parseServerCommand()
		server.New().Start(*host, *port)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
