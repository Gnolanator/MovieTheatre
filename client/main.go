package main

import (
  "grpcPractise/proto"
  "context"
  "os"
  "bufio"
  "fmt"
  "time"
  "flag"
  "sync"
  "github.com/gin-gonic/gin"
  grpc "google.golang.org/grpc"
  glog "google.golang.org/grpc/grpclog"
)


func ginStuffs() {
  r := gin.Default()
  r.GET("ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  r.Run("0.0.0.0:8000")
}


var client proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
  wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
  var streamerror error

  stream, err := client.CreateStream(context.Background(), &proto.Connect{
    User: user,
    Active: true,
  })

  if err != nil {
    return fmt.Errorf("Connection failed %v", err)
  }

  wait.Add(1)
  go func(str proto.Broadcast_CreateStreamClient) {
    defer wait.Done()

    for {
      msg, err := str.Recv()
      if err != nil {
        streamerror = fmt.Errorf("Error reading the message %v", err)
        break
      }

      fmt.Printf("%v : %s\n", msg.Id, msg.Content)
    }
  }(stream)

  return streamerror
}

func main() {
  ginStuffs()
  timestamp := time.Now()
  done := make(chan int)

  name := flag.String("N", "Anon", "The name of the user")
  flag.Parse()

  id := "1" + *name

  conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
  if err != nil {
    glog.Fatalf("Could not connect to service: %v", err)
  }

  client = proto.NewBroadcastClient(conn)
  user := &proto.User{
    Id: id,
    Name: *name,
  }

  connect(user)

  wait.Add(1)
  go func(){
    defer wait.Done()

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
      msg := &proto.Message{
        Id: user.Id,
        Content: scanner.Text(),
        Timestamp: timestamp.String(),
      }

      _, err := client.BroadcastMessage(context.Background(), msg)
      if err != nil {
        fmt.Printf("Error Sending Message: %v", err)
        break
      }
    }
  }()

  go func() {
    wait.Wait()
    close(done)
  }()

  <-done
}
