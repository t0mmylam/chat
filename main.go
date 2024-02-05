package main

func main() {
	server := newServer()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("unable to start server: ", err)
	}

	defer listener.Close()
	log.Printf("started server on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handle(conn, server)
	}
}