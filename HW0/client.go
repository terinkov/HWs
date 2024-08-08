package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	serv := os.Args[1]               // берем адрес сервера из аргументов командной строки
	conn, _ := net.Dial("tcp", serv) // открываем TCP-соединение к серверу
	defer conn.Close()
	checkServAns(conn)

}

func checkServAns(src io.Reader) {
	message, _ := bufio.NewReader(src).ReadString('\n')
	// fmt.Print("Message from server: "+message)
	if message != "OK" {
		fmt.Println("Некорректное сообщение от сервера:", message, '\n')
		return
	}

	fmt.Println("Получено сообщение от сервера:", message, '\n')
}
