package main

import (
	"fmt"
	"net"
	"os"
)

// Сервер отправляет клиенту с четным номером "OK", а с нечетным "NOT OK"
func main() {
	listener, _ := net.Listen("tcp", "localhost:8080") // открываем слушающий сокет
	i := 0
	for {
		fmt.Println("Клиент под номером:", i)
		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while client connection")
			return
		}
		var phrase string
		if i%2 == 0 {
			phrase = "OK"
		} else {
			phrase = "Not ok"
		}
		handleClient(conn, phrase)
		i += 1
	}
}

func handleClient(conn net.Conn, phrase string) {
	defer conn.Close()
	fmt.Println("Connection closed")
	conn.Write([]byte(phrase)) // пишем в сокет
	fmt.Println(phrase + " send to client")
}
