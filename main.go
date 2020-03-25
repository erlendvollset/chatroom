package main

import "flag"
import "fmt"
import "github.com/erlendvollset/chatroom/client"
import "github.com/erlendvollset/chatroom/server"
import "os"
import "strings"

func main() {
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)

	userName := clientCommand.String("N", "Anon", "The name of the user")

	switch strings.ToLower(os.Args[1]) {
		case "client":
			clientCommand.Parse(os.Args[2:])
			client.StartClient(userName)
		case "server":
			serverCommand.Parse(os.Args[2:])
			server.StartServer(":8080")
		default:
			fmt.Printf("%q is not valid command.\n", os.Args[1])
			os.Exit(2)
	}
}
