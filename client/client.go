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

var client proto.ChatroomClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("Connection failed: %w", err)
	}
	wait.Add(1)
	go func(str proto.Chatroom_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %w", err)
				break
			}
			fmt.Printf("%v : %s\n", msg.Id, msg.Content)
		}
	}(stream)
	return streamerror
}

func StartClient(userName *string) {
	timestamp := time.Now()
	done := make(chan int)
	id := sha256.Sum256([]byte(timestamp.String() + *userName))
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to service: %v", err)
	}

	client = proto.NewChatroomClient(conn)
	user := &proto.User{
		Id: hex.EncodeToString(id[:]),
		Name: *userName,
	}
	connect(user)
	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &proto.Message{
				Id:                   user.Id,
				Content:              scanner.Text(),
				Timestamp:            timestamp.String(),
			}
			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error sending message: %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()
	<- done
}
