package main

import "flag"
import "fmt"
import "github.com/erlendvollset/chatroom/client"
import "github.com/erlendvollset/chatroom/server"
import "os"
import "strings"

type NetworkConfig struct {
	Host *string
	Port *string
}

type ClientConfig struct {
	NetworkConfig
	userName *string
}

type ServerConfig struct {
	NetworkConfig
}

func parseClientCommand() ClientConfig {
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)
	config := ClientConfig{
		NetworkConfig{
			clientCommand.String("h", "localhost", "The host for the chatroom server"),
			clientCommand.String("p", "8080", "The port of the chatroom server")},
		clientCommand.String("n", "Anon", "The name of the user"),
	}
	clientCommand.Parse(os.Args[2:])
	return config
}

func parseServerCommand() ServerConfig {
	serverCommand := flag.NewFlagSet("server", flag.ExitOnError)
	config := ServerConfig{
		NetworkConfig{
		serverCommand.String("h", "localhost", "The host for the chatroom server"),
		serverCommand.String("p", "8080", "The port of the chatroom server")},
	}
	serverCommand.Parse(os.Args[2:])
	return config
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: chatroom {client|server} [args]")
		os.Exit(0)
	}
	switch strings.ToLower(os.Args[1]) {
	case "client":
		config := parseClientCommand()
		client.New(*config.userName).Start(*config.Host, *config.Port)
	case "server":
		config := parseServerCommand()
		server.New().Start(*config.Host, *config.Port)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
