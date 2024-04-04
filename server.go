// server.go
package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Addr     *net.UDPAddr
	Msgs     map[string]bool
	LastSeen time.Time
}

var clients = make(map[string]*Client)

func cleanupClients() {
	for {
		time.Sleep(5 * time.Minute)
		for addr, client := range clients {
			if time.Since(client.LastSeen) > 5*time.Minute {
				delete(clients, addr)
			}
		}
	}
}

func handleClient(conn *net.UDPConn, clientAddr *net.UDPAddr, msg string) {
	if _, ok := clients[clientAddr.String()]; !ok {
		clients[clientAddr.String()] = &Client{
			Addr:     clientAddr,
			Msgs:     make(map[string]bool),
			LastSeen: time.Now(),
		}
	} else {
		clients[clientAddr.String()].LastSeen = time.Now()
	}

	clients[clientAddr.String()].Msgs[msg] = true

	fmt.Printf("Received message from %s: %s\n", clientAddr.String(), msg)

	for _, client := range clients {
		if client.Addr.String() != clientAddr.String() {
			if !client.Msgs[msg] {
				buffer := []byte(msg)
				if _, err := conn.WriteToUDP(buffer, client.Addr); err != nil {
					fmt.Println(err)
					delete(clients, client.Addr.String())
				} else {
					client.Msgs[msg] = true
				}
			}
		}
	}
}

func main() {
	ipStr := flag.String("ip", "127.0.0.1", "The IP address to listen on")
	port := flag.Int("port", 12345, "The port to listen on")
	flag.Parse()

	go cleanupClients()

	ip := net.ParseIP(*ipStr)
	if ip == nil {
		fmt.Println("Invalid IP address")
		return
	}

	addr := net.UDPAddr{
		Port: *port,
		IP:   ip,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := string(buffer[:n])
		go handleClient(conn, clientAddr, string(msg)) // Pass a copy of the message
	}
}
