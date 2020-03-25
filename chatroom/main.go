package main

import "flag"
import "fmt"
import "github.com/erlendvollset/chatroom/client"
import "github.com/erlendvollset/chatroom/server"
import "os"
import "strings"

func main() {
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)
	userName := clientCommand.String("N", "Anon", "The name of the user")

	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)

	switch strings.ToLower(os.Args[1]) {
	case "client":
		clientCommand.Parse(os.Args[2:])
		client.New(userName).Start("8080")
	case "server":
		serverCommand.Parse(os.Args[2:])
		server.New().Start("8080")
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
