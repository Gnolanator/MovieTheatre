package main

import (
  "os"
  "sync"
  "net"
  "context"
  "log"
  proto "grpcPractise/proto"
  grpc "google.golang.org/grpc"
  glog "google.golang.org/grpc/grpclog"
  vlc "github.com/adrg/libvlc-go"
)

//================================================
func streamVideo(videoPath string) {

  if err := vlc.Init("-vvv", "--sout udp:192.168.0.42", "--ttl 12"); err != nil {
    log.Fatal(err)
  }
  defer vlc.Release()

  player, err := vlc.NewPlayer()
  if err != nil {
    log.Fatal(err)
  }
  defer func() {
    player.Stop()
    player.Release()
  }()

  media, err := player.LoadMediaFromURL(videoPath)
  if err != nil {
    log.Fatal(err)
  }
  defer media.Release()

  err = player.Play()
  if err != nil {
    log.Fatal(err)
  }

  manager, err := player.EventManager()
  if err != nil {
    log.Fatal(err)
  }

  quit := make(chan struct{})
  eventCallback := func(event vlc.Event, userData interface{}) {
    close(quit)
  }

  eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
  if err != nil {
    log.Fatal(err)
  }
  defer manager.Detach(eventID)

  <-quit
}
//================================================

var grpcLog glog.LoggerV2

func init() {
  grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
  stream proto.Broadcast_CreateStreamServer
  id string
  active bool
  error chan error
}

type Server struct {
  Connection []*Connection
}

func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
  conn := &Connection{
    stream: stream,
    id: pconn.User.Id,
    active: true,
    error: make(chan error),
  }

  s.Connection = append(s.Connection, conn)

  return <-conn.error
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
  wait := sync.WaitGroup{}
  done := make(chan int)

  for _, conn := range s.Connection {
    wait.Add(1)

    go func(msg *proto.Message, conn *Connection){
      defer wait.Done()

      if conn.active {
        err := conn.stream.Send(msg)
        grpcLog.Info("Sending message to: ", conn.stream)

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
  return &proto.Close{}, nil
}

func main() {
  var connections []*Connection

  server := &Server{connections}

  grpcServer := grpc.NewServer()
  listener, err := net.Listen("tcp", ":8080")
  if err != nil {
    glog.Fatalf("error creating the server %v", err)
  }

  grpcLog.Info("Starting server at port 8080")
  proto.RegisterBroadcastServer(grpcServer, server)
  grpcServer.Serve(listener)
}
