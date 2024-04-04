// client.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func getNickname() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	return strings.TrimSpace(nickname)
}

func main() {
	serverIP := flag.String("server-ip", "127.0.0.1", "The IP address of the server")
	serverPort := flag.Int("server-port", 12345, "The port of the server")
	flag.Parse()

	nickname := getNickname()

	ip := net.ParseIP(*serverIP)
	if ip == nil {
		fmt.Println("Invalid IP address")
		return
	}

	addr := net.UDPAddr{
		Port: *serverPort,
		IP:   ip,
	}

	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(buffer[:n]))
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if _, err := conn.Write([]byte(nickname + ": " + msg)); err != nil {
			fmt.Println(err)
			return
		}
	}
}
