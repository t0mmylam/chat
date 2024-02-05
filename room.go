package main

type room struct {
	name string
	members map[net.Conn]*client
}