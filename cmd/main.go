package main

import (
	"fmt"
	"lesta-start-battleship/cli/internal/app"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("app.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	app, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}

	fmt.Println("CLI клиент запущен.")
	log.Print("CLI клиент запущен.")
	if err := app.Run(); err != nil {
		log.Printf("Ошибка выполнения: %v", err)
		os.Exit(1)
	}
}
