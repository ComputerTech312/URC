// client.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// getNickname prompts the user for their nickname and returns it.
func getNickname() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	return strings.TrimSpace(nickname)
}

// readFromServer continuously reads messages from the server and prints them.
func readFromServer(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		// Get the current time and format it
		timestamp := time.Now().Format("15:04:05")

		// Print the message with the timestamp
		fmt.Printf("[%s] %s\n", timestamp, string(buffer[:n]))
	}
}

// writeToServer continuously reads messages from the user and sends them to the server.
// writeToServer continuously reads messages from the user and sends them to the server.
func writeToServer(conn *net.UDPConn, nickname string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		// Include the nickname in the message
		fullMsg := fmt.Sprintf("%s: %s", nickname, msg)

		if _, err := conn.Write([]byte(fullMsg)); err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}
	}
}

func main() {
	// Parse command line arguments
	serverIP := flag.String("server-ip", "127.0.0.1", "The IP address of the server")
	serverPort := flag.Int("server-port", 12345, "The port of the server")
	flag.Parse()

	// Get the user's nickname
	nickname := getNickname()

	// Parse the server IP
	ip := net.ParseIP(*serverIP)
	if ip == nil {
		fmt.Println("Invalid IP address")
		return
	}

	// Create the server address
	addr := net.UDPAddr{
		Port: *serverPort,
		IP:   ip,
	}

	// Connect to the server
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start reading from the server
	go readFromServer(conn)

	// Start writing to the server
	writeToServer(conn, nickname)
}
