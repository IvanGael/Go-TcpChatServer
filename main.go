package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
	name string
}

var clients []Client

func main() {
	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	client := Client{conn: conn}

	client.name = readMessage(conn)

	clients = append(clients, client)

	fmt.Printf("%s has joined the chat.\n", client.name)

	for {
		msg := readMessage(conn)

		if msg == "/quit" {
			fmt.Printf("%s has left the chat.\n", client.name)
			removeClient(client)
			return
		}

		broadcastMessage(client, msg)
	}
}

func readMessage(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return ""
	}
	return strings.TrimSpace(message)
}

func broadcastMessage(sender Client, message string) {
	fmt.Printf("[%s]: %s\n", sender.name, message)
	for _, client := range clients {
		if client.name != sender.name {
			fmt.Fprintf(client.conn, "[%s]: %s\n", sender.name, message)
		}
	}
}

func removeClient(client Client) {
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// telnet 127.0.0.1 8080
