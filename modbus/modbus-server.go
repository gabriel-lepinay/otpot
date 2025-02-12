package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Create a logs directory if it doesn't exist
	err := os.MkdirAll("/logs", 0777)
	if err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		return
	}

	// Open the log file
	logFile, err := os.OpenFile("/logs/modbus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Set up multi-writer to log to both the terminal and the file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// Create a TCP listener
	listener, err := net.Listen("tcp", "0.0.0.0:502")
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer listener.Close()
	log.Println("Modbus server listening on port 502")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("Connection established from %s", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	originIP := conn.RemoteAddr().String()

	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			// If the error is EOF, the connection was closed by the client, which is expected.
			if err.Error() == "EOF" {
				log.Printf("Client %s disconnected.", originIP)
				break
			} else {
				log.Printf("Error reading data from %s: %v", originIP, err)
				break
			}
		}
		log.Printf("Received data from %s: %x\n", originIP, data[:n])

		// Process the received data
		response := processData(data[:n])

		// Echo the processed data back to the client
		_, err = conn.Write(response)
		if err != nil {
			log.Printf("Error writing data to %s: %v", originIP, err)
			return
		}
	}
}

func processData(data []byte) []byte {
	// Implement your data processing logic here
	// For demonstration, let's just return the received data as-is
	return data
}
