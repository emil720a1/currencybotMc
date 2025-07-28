package main

import (
	"context"
	"currency-bot/internal/delivery/telegram"
	"currency-bot/internal/handler"
	"currency-bot/internal/service"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"log"
	"os"
	"time"
)

func main() {
	botFile, errBotFile := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errBotFile != nil {
		fmt.Fprintf(os.Stderr, "Не вдалося відкрити bot.log: %v\n", errBotFile)
		os.Exit(1)
	}
	log.SetOutput(botFile) //instead of os.Stdout we set botFile as output
	defer botFile.Close()
	log.SetFlags(log.LstdFlags) //add date and time to log
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Помилка завантаження .env файлу")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, errDb := sql.Open("mysql", dsn)
	if errDb != nil {
		log.Fatal(errDb)
	}

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal(errPing)
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Не встановлено значення BOT_TOKEN")
	}
	log.Println(botToken)
	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatal("Помилка створення бота:", err)
	}

	service := service.NewCurrencyService(db)
	service.RatesCache = ""
	service.CacheTime = time.Time{}
	handler := handler.NewCurrencyHandler(service)
	delivery := telegram.NewBot(bot, handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params := &telego.GetUpdatesParams{}
	updates, err := bot.UpdatesViaLongPolling(ctx, params)
	if err != nil {
		log.Fatal("Помилка запуску Long Polling: ", err)
	}

	delivery.Start(updates)

}
