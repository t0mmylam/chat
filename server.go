package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_LEAVE:
			s.quit(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {
	now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
	log.Printf("[%s] SERVER: %s has joined the chat.", now, conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
}

func (s *server) nick(c *client, args []string) {
	now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
	if len(args) < 2 {
		c.msg(fmt.Sprintf("[%s] SERVER: must provide nickname | Ex: usage: /nick <name>", now))
		return
	} else if c.nick == args[1] {
		c.msg(fmt.Sprintf("[%s] SERVER: nickname is already %s", now, c.nick))
		return
	} else if strings.ToLower(args[1]) == "server" {
		c.msg(fmt.Sprintf("[%s] SERVER: nickname cannot be: %s", now, c.nick))
		return
	}
	c.nick = args[1]
	c.msg(fmt.Sprintf("[%s] SERVER: nickname changed to: %s", now, c.nick))
}

func (s *server) join(c *client, args []string) {
	now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
	if len(args) < 2 {
		c.msg(fmt.Sprintf("[%s] -- must provide room name | Ex: usage: /join ROOM_NAME", now))
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("[%s] SERVER: %s joined the room", now, c.nick))

	c.msg(fmt.Sprintf("[%s] SERVER: welcome to %s", now, roomName))
}

func (s *server) listRooms(c *client) {
	now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("[%s] SERVER: available rooms: %s", now, strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
	if c.room == nil {
		c.msg(fmt.Sprintf("[%s] SERVER: must join a room first | Ex: usage: /join ROOM_NAME", now))
		return
	}
	msg := strings.Join(args, " ")
	c.room.broadcast(c, "["+now+"]"+" "+c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
	log.Printf("[%s] SERVER: %s has disconnected.", now, c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		now := time.Now().Format("YYYY-MM-DD HH:MM:SS")
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("[%s] SERVER: %s has left the room.", now, c.nick))
	}
}
