package main

import (
	"bufio"
	"fmt"
	"github.com/git-hyagi/janaina/pkg/chat"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
)

const (
	servername = "golab"
	idPerson   = 2
	username   = "jose"
)

func main() {

	var wg sync.WaitGroup
	var conn *grpc.ClientConn

	// make connection
	conn, err := grpc.Dial(servername+":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer conn.Close()

	// create GRPC chat service
	chatGRPC := chat.NewChatClient(conn)

	// call the sendMessage method from grpc
	stream, _ := chatGRPC.SendMessage(context.Background())

	// this goroutine is for "watching" messages being received
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			newMsg, _ := stream.Recv()
			fmt.Printf("[%s] %s: %s\n", newMsg.Timestamp, newMsg.Message.Username, newMsg.Message.Content)
		}
	}()

	// this other goroutine is to send a message and let the program in a non-blocking state while a message is not gathered yet(scanf)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			msgBody := bufio.NewScanner(os.Stdin)
			msgBody.Scan()

			msg := &chat.Message{IdPerson: idPerson, Username: username, Content: msgBody.Text()}

			stream, _ := chatGRPC.SendMessage(context.Background())
			stream.Send(msg)
		}
	}()

	wg.Wait()

}
