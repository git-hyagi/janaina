package chat

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"
)

type Server struct {
	Connections []*Connection
}

type Connection struct {
	chatConn Chat_SendMessageServer

	// This channel is the only way I found to avoid the stream channel from flutter/dart not getting "closed" before receive the message
	// without it I was getting 'rpc: context canceled' for every message trying to be delivered to the dart client
	wchan chan error
}

func (s *Server) SendMessage(stream Chat_SendMessageServer) error {

	// wait group to sync with the go routine used to send message to clients
	var wg sync.WaitGroup

	con := &Connection{chatConn: stream, wchan: make(chan error)}

	// for each client connected, add them to the list of connections
	s.Connections = append(s.Connections, con)

	log.Printf("connected! %v", stream)
	log.Printf("connections! %v", s.Connections)

	// gather the message from client
	msg, err := con.chatConn.Recv()
	if err != nil {
		log.Printf("error receiving msg!!\n")
		<-con.wchan
		return status.Error(codes.Unknown, "Failed to receive client data!")
	}

	log.Printf("Message received from: [%v] %v", con.chatConn, msg)

	// format it according to server response message type
	rsp_message := &ServerMessage{Message: msg, Timestamp: time.Now().Format("03:04")}

	// for each client connected
	for i, connected := range s.Connections {

		// create a goroutine to each connected client to send the messages
		wg.Add(1)
		go func(i int, connected *Connection) {
			defer wg.Done()
			log.Printf("sending message to: %v", connected.chatConn)

			// send the message
			err = connected.chatConn.Send(rsp_message)
			if err != nil {
				log.Printf("Error trying to send msg to %v: %v", connected.chatConn, err)
				s.removeFromSlice(i)
				<-connected.wchan
			}
		}(i, connected)
	}
	wg.Wait()

	return <-con.wchan

}

// function to remove a "closed" grpc connection
// this is kind of a "garbage collector" as the flutter/dart program get into a block state when it tried to send a message to a closed connection
func (s *Server) removeFromSlice(i int) {
	log.Printf("Removing connection from list: %v", s.Connections[i].chatConn)
	s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
	log.Printf("connections: %v", s.Connections)
}
