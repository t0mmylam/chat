package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_LEAVE
	// TODO: Add more commands
	// CMD_ADMIN
	// CMD_HELP
	// CMD_KICK
	// CMD_BAN
	// CMD_UNBAN
	// CMD_QUIT
	// CMD_USERS
	// CMD_ADMINS
	// CMD_BROADCAST
	// CMD_WHISPER
	// CMD_MUTE
	// CMD_UNMUTE
)

type command struct {
	id     commandID
	client *client
	args   []string
}
