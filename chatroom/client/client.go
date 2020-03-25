package client

import "bufio"
import "context"
import "crypto/sha256"
import "encoding/hex"
import "fmt"
import "github.com/erlendvollset/chatroom/proto"
import "google.golang.org/grpc"
import "log"
import "os"
import "sync"
import "time"

type Client interface {
	Start(host string, port string)
}

type client struct {
	chatroomClient proto.ChatroomClient
	user           *proto.User
	wait           *sync.WaitGroup
}

func generateId(userName string) string {
	id := sha256.Sum256([]byte(time.Now().String() + userName))
	return hex.EncodeToString(id[:])
}

func New(userName *string) Client {
	user := &proto.User{Id: generateId(*userName), Name: *userName}
	wait := &sync.WaitGroup{}
	return &client{user: user, wait: wait}
}

func (c *client) streamMessages() error {
	stream, err := c.chatroomClient.CreateStream(context.Background(), &proto.Connect{
		User:   c.user,
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("Connections failed: %w", err)
	}
	c.wait.Add(1)
	go func(str proto.Chatroom_CreateStreamClient) {
		defer c.wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				break
			}
			fmt.Printf("%v: %s\n", msg.User.Name, msg.Content)
		}
	}(stream)
	return nil
}

func (c *client) scanAndBroadcast() {
	c.wait.Add(1)
	go func() {
		defer c.wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &proto.Message{
				User:    c.user,
				Content: scanner.Text(),
			}
			_, err := c.chatroomClient.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error sending message: %v", err)
				break
			}
		}
	}()
}

func (c *client) Start(host string, port string) {
	done := make(chan int)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to service: %v", err)
	}
	c.chatroomClient = proto.NewChatroomClient(conn)
	if err := c.streamMessages(); err != nil {
		log.Fatalf("Error connecting to the server: %v", err)
	}
	c.scanAndBroadcast()
	go func() {
		c.wait.Wait()
		close(done)
	}()
	<-done
}
