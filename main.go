package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Fprintln(conn, "[СЕРВЕР] Ласкаво просимо до мережевого сервісу диспетчеризації!")

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		command := strings.TrimSpace(message)

		switch {
		case strings.HasPrefix(command, "ECHO "):
			text := strings.TrimPrefix(command, "ECHO ")
			fmt.Fprintf(conn, "[ECHO]: %s\n", text)

		case command == "TIME":
			fmt.Fprintf(conn, "[TIME]: %s\n", time.Now().Format(time.RFC3339))

		case command == "CLOSE" || command == "QUIT":
			fmt.Fprintln(conn, "[СЕРВЕР]: До побачення!")
			return

		default:
			fmt.Fprintln(
				conn,
				"[ПОМИЛКА]: Невідома команда. Доступні команди: ECHO <повідомлення>, TIME, CLOSE, QUIT",
			)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("=== Мережевий сервіс диспетчеризації ===")
	fmt.Println("[СЕРВЕР] TCP-слухач успішно відкрито на порту 9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Помилка прийняття з'єднання:", err)
			continue
		}

		fmt.Println("[ПІДКЛЮЧЕННЯ] Новий клієнт:", conn.RemoteAddr())

		go handleConnection(conn)
	}
}
