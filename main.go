package main

import (
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("TG_TOKEN")
	chatID := os.Getenv("TG_CHAT_ID")

	if token == "" || chatID == "" {
		fmt.Println("Error: TG_TOKEN and TG_CHAT_ID required")
		os.Exit(1)
	}

	fmt.Printf("Bot started. Token: %s, ChatID: %s\n", token[:10]+"...", chatID)
}
