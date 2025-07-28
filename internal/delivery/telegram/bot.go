package telegram

import (
	"context"
	"currency-bot/internal/handler"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type Bot struct {
	bot     *telego.Bot
	handler *handler.CurrencyHandler
}

func NewBot(bot *telego.Bot, handler *handler.CurrencyHandler) *Bot {
	return &Bot{bot, handler}
}

func (b *Bot) Start(updates <-chan telego.Update) {

	stopChan := make(chan struct{})
	go b.startCurrencyChecker(stopChan)
	defer close(stopChan)

	for update := range updates {
		if update.Message != nil {

			switch true {
			case update.Message.Text == "/start":

				chatID := update.Message.Chat.ID
				text := b.handler.HandleStart(update.Message.Chat.ID)
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}

				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Println("Помилка надсилання повідомлення до chatID %d: %v (HandleStart)", chatID, err)
				}

			case strings.HasPrefix(update.Message.Text, "/rates"):

				var money string
				chatID := update.Message.Chat.ID
				text, errMoney := b.handler.HandleRates(chatID, money)
				if errMoney != nil {
					text = "Помилка: " + errMoney.Error()
				}
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}

				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleRates)", chatID, err)
				}

			case update.Message.Text == "/help":
				chatID := update.Message.Chat.ID
				text := b.handler.HandleHelp(chatID)
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}

				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleHelp)", chatID, err)
				}

			case update.Message.Text == "/about":
				chatID := update.Message.Chat.ID
				text := b.handler.HandleAbout(chatID)
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}
				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleAbout)", chatID, err)
				}

			case strings.HasPrefix(update.Message.Text, "/convert"):
				parts := strings.Fields(update.Message.Text)
				chatID := update.Message.Chat.ID
				if len(parts) != 4 {
					msgParams := &telego.SendMessageParams{
						ChatID:    tu.ID(chatID),
						Text:      "Використовуйте: /convert <сума> <валюта1> <валюта2>, наприклад, /convert 100 USD UAN ",
						ParseMode: "Markdown",
					}
					_, err := b.bot.SendMessage(context.Background(), msgParams)
					if err != nil {
						log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleConvert)", chatID, err)
					}

				} else {
					amount, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						text := "Невірний формат суми"
						log.Printf("chatID %d: Невірний формат суми: %s", parts[1])
						msgParams := &telego.SendMessageParams{
							ChatID:    tu.ID(chatID),
							Text:      text,
							ParseMode: "Markdown",
						}
						_, err := b.bot.SendMessage(context.Background(), msgParams)
						if err != nil {
							log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleConvert)", chatID, err)
						}
					} else {
						log.Printf("Викликаємо HandleConvert: amount= %v, from= %v, to= %s", amount, parts[2], parts[3])
						text, err := b.handler.HandleConvert(chatID, amount, parts[2], parts[3])
						if err != nil {
							text = "Помилка: " + err.Error()
						}
						msgParams := &telego.SendMessageParams{
							ChatID:    tu.ID(chatID),
							Text:      text,
							ParseMode: "Markdown",
						}
						_, err = b.bot.SendMessage(context.Background(), msgParams)
						if err != nil {
							log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleConvert)", chatID, err)
						}
					}
				}
			case update.Message.Text == "/usd":
				chatID := update.Message.Chat.ID
				text, errMoney := b.handler.HandleCurrency(chatID, "USD")
				if errMoney != nil {
					text = "Помилка: " + errMoney.Error()
				}
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}
				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleCurrency USD)", chatID, err)
				}

			case update.Message.Text == "/eur":
				chatID := update.Message.Chat.ID
				text, errMoney := b.handler.HandleCurrency(chatID, "EUR")
				if errMoney != nil {
					text = "Помилка: " + errMoney.Error()
				}
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}
				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleCurrency EUR)", chatID, err)
				}
			case update.Message.Text == "/history":
				chatID := update.Message.Chat.ID
				log.Printf("chatID %d: Викликано /history", chatID)
				text, err := b.handler.HandleHistory(chatID)
				if err != nil {
					text = "Помилка: " + err.Error()
				}

				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "Markdown",
				}

				_, err = b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleHistory)", chatID, err)
				}

			case update.Message.Text == "/stats":
				chatID := update.Message.Chat.ID
				log.Printf("chatID %d: Викликано /stats", chatID)
				text, err := b.handler.HandleStats(chatID)
				if err != nil {
					text = "Помилка: " + err.Error()
				}

				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "",
				}
				_, err = b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleStats)", chatID, err)
				}
			case update.Message.Text == "/clearhistory":
				chatID := update.Message.Chat.ID
				log.Printf("chatID %d: Викликано /clearhistory", chatID)
				text, err := b.handler.HandleClearHistory(chatID)
				if err != nil {
					text = "Помилка: " + err.Error()
				}

				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "",
				}

				_, err = b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleClearHistory)", chatID, err)
				}

			case strings.HasPrefix(update.Message.Text, "/lang"):
				chatID := update.Message.Chat.ID
				log.Printf("chatID %d: Викликано /lang", chatID)
				lang := strings.TrimPrefix(update.Message.Text, "/lang")
				lang = strings.TrimSpace(lang)
				text, err := b.handler.HandleLang(chatID, lang)
				if err != nil {
					text = "Помилка: " + err.Error()
				}

				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      text,
					ParseMode: "",
				}

				_, err = b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (HandleLang)", chatID, err)
				}
			case true:
				chatID := update.Message.Chat.ID
				msgParams := &telego.SendMessageParams{
					ChatID:    tu.ID(chatID),
					Text:      "Невідома команда, спробуй /start, /rates або /help",
					ParseMode: "Markdown",
				}
				_, err := b.bot.SendMessage(context.Background(), msgParams)
				if err != nil {
					log.Printf("Помилка надсилання повідомлення до chatID %d: %v (Default)", chatID, err)
				}
			}
		}
	}
}

func (b *Bot) startCurrencyChecker(stopChan <-chan struct{}) {
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	sliceCourse := []string{"USD", "EUR"}

	go func() {
		rates := make(map[string]string)
		prevRates := make(map[string]string)
		for range ticker.C {
			for _, t := range sliceCourse {
				course, err := b.handler.Service.GetRates(t)
				if err != nil {
					log.Printf("Помилка оновлення курсу %v", err)
					continue
				} else {
					rates[t] = course
				}
				log.Printf("Course: %s", course)
				parseValue := strings.TrimSpace(strings.Split(strings.Split(strings.Split(course, ": ")[1], ",")[0], " ")[1])
				log.Printf("ParseValue: %s", parseValue)
				parseValueFloat, err := strconv.ParseFloat(parseValue, 64)
				if err != nil {
					log.Printf("Помилка парсингу курсу %v", err)
					continue
				}

				if parseValueFloat == 0.0 {
					continue
				}
				log.Printf("Курс числа: %f", parseValueFloat)

				if prevRates[t] == "" {
					prevRates[t] = course
				} else {
					parsePrev := strings.TrimSpace(strings.Split(strings.Split(strings.Split(prevRates[t], ": ")[1], ",")[0], " ")[1])
					log.Printf("ParsePrev: %s", parsePrev)
					if parsePrev == "" {
						continue
					}
					parsePrevFloat, err := strconv.ParseFloat(parsePrev, 64)
					if err != nil {
						log.Printf("Помилка парсингу курсу %v", err)
						continue
					}
					if parsePrevFloat == 0 {
						continue
					}
					res := (parseValueFloat - parsePrevFloat) / parsePrevFloat * 100
					log.Printf("Різниця курсу: %.2f", res)
					if res > 5 {
						msgParams := &telego.SendMessageParams{
							ChatID:    tu.ID(927118456),
							Text:      "Курс " + t + " змінився на " + course + " " + "різниця: " + fmt.Sprintf("%.2f%%", res),
							ParseMode: "Markdown",
						}
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						_, err := b.bot.SendMessage(ctx, msgParams)
						if err != nil {
							log.Printf("Помилка надсилання повідомлення до chatID %d: %v (Default)", 927118456, err)
						}

					}
					prevRates[t] = course

				}
				log.Printf("Оновлено %s: %s", t, rates[t])

			}
			select {
			case <-stopChan:
				break
			}
		}
	}()

}
