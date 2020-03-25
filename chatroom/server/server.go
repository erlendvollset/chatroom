package server

import "context"
import "fmt"
import "github.com/erlendvollset/chatroom/proto"
import "google.golang.org/grpc"
import "google.golang.org/grpc/grpclog"
import "log"
import "net"
import "os"
import "sync"

type Server interface {
	Start(host string, port string)
}

type connection struct {
	stream proto.Chatroom_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type server struct {
	proto.UnimplementedChatroomServer
	Connections []*connection
	grpcLog     grpclog.LoggerV2
}

func New() Server {
	var connections []*connection
	return &server{Connections: connections, grpcLog: grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)}
}

func (s *server) CreateStream(pconn *proto.Connect, stream proto.Chatroom_CreateStreamServer) error {
	conn := &connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}
	s.Connections = append(s.Connections, conn)
	return <-conn.error
}

func (s *server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connections {
		wait.Add(1)
		go func(msg *proto.Message, conn *connection) {
			defer wait.Done()
			if conn.active {
				err := conn.stream.Send(msg)
				s.grpcLog.Info("Sending message to:  ", conn.stream)
				if err != nil {
					s.grpcLog.Errorf("Error with Stream: %s - Error: %v", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &proto.Empty{}, nil
}

func (s *server) Start(host string, port string) {
	grpcServer := grpc.NewServer()
	address := fmt.Sprintf("%v:%v", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}
	s.grpcLog.Infof("Starting server at address: %s", address)
	proto.RegisterChatroomServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("error occurred while serving %v", err)
	}
}
