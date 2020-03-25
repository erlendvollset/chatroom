package server

import "context"
import "github.com/erlendvollset/chatroom/proto"
import "google.golang.org/grpc"
import "google.golang.org/grpc/grpclog"
import "log"
import "net"
import "os"
import "sync"

var grpcLog grpclog.LoggerV2

func init() {
	grpcLog = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream proto.Chatroom_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	proto.UnimplementedChatroomServer
	Connection []*Connection
}

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Chatroom_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}
	s.Connection = append(s.Connection, conn)

	return <-conn.error
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()
			if conn.active {
				err := conn.stream.Send(msg)
				grpcLog.Info("Sending message to:  ", conn.stream)
				if err != nil {
					grpcLog.Errorf("Error with Stream: %s - Error: %v", conn.stream, err)
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

func StartServer(port string) {
	var connections []*Connection
	chatroomServer := &Server{Connection: connections}
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error creating the Server %v", err)
	}
	grpcLog.Info("Starting Server at port :8080")
	proto.RegisterChatroomServer(grpcServer, chatroomServer)
	grpcServer.Serve(listener)
}
