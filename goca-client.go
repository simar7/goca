package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	connection, _ := net.Dial("tcp", "localhost:8081")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the user identity: ")
	user_id, _ := reader.ReadString('\n')

	fmt.Fprintf(connection, user_id+"\n")

	msg, _ := bufio.NewReader(connection).ReadString('\n')
	fmt.Print("Message from server: " + msg)
}
