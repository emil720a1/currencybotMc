package handler

import (
	"currency-bot/internal/service"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var translations = map[string]map[string]string{
	"en": {
		"start":                       "Hello! Use /rates to get currency rates",
		"help":                        "Use /start to start the bot\nUse /rates to get currency rates\nUse /convert <amount> <from> <to> to convert currency",
		"about":                       "This bot shows currency rates from Monobank\nAuthor: McAleg - @emil720a",
		"history":                     "History:",
		"history_empty":               "History is empty",
		"stats":                       "Statistics:",
		"stats_empty":                 "Statistics is empty",
		"clear_history":               "History cleared",
		"lang":                        "Language:",
		"currency":                    "Currency:",
		"lang_success":                "Language saved",
		"error_invalid_currencyRATES": "Unknown currency, use USD or EUR",
		"error_receiving_ratesRATES":  "Error receiving currency rates",
		"error_parse_Convert":         "Error parsing currency rates for",
		"error_save_Convert":          "Error saving conversion",
		"error_find_rates":            "Currency rates not found",
		"error_handle_currency":       "Currency not supported, use USD or EUR",
		"error_parse_handleCurrency":  "Error parsing currency",
		"error_scan_history":          "Error scanning history",
		"error_receive_history":       "Error receiving history",
		"error_receive_stats":         "Error receiving statistics",
		"log_stats":                   "Received statistics:",
		"error_clear_history":         "Error clearing history",
		"error_lang":                  "Language not supported, use en or uk",
		"error_save_language":         "Error saving language",
		"save_conversation":           "Saved conversion:",
		"error_format_row":            "Invalid row format:",
		"getRates_return":             "GetRates returned for:",
		"course_not_found":            "Course not found",
		"row_is_empty":                "Row is empty",
		"row_is_invalid":              "Invalid row format",
		"course_is_invalid":           "Invalid course format for",
		"course_parse_error":          "Error parsing at buy for",
		"course_parse_error_sell":     "Error parsing at sell for",
		"buy":                         "Buy",
		"sell":                        "Sell",
		"error_format_course":         "Currency not in format",
	},
	"uk": {
		"start":                         "Привіт! Використовуйте /rates для отримання курсів валют",
		"help":                          "Використовуйте /start для старту роботи з ботом\nВикористовуйте /rates для отримання курсів валют\nВикористовуйте /convert <сума> <з_валюти> <в_валюту> для конвертування валюти",
		"about":                         "Цей бот показує курси валют від Monobank\nАвтор: McAleg - @emil720a",
		"history":                       "Історія:",
		"history_empty":                 "Історія порожня",
		"stats":                         "Статистика:",
		"stats_empty":                   "Статистика порожня",
		"clear_history":                 "Історія очищена",
		"lang":                          "Мова:",
		"currency":                      "Валюта:",
		"lang_success":                  "Мова збережена",
		"error_invalid_currencyRATES":   "Невідома валюта, використовуйте USD або EUR",
		"error_invalid_currency-handle": "Невідома валюта, використовуйте USD, EUR, UAN",
		"error_receiving_ratesRATES":    "Помилка при отриманні курсів валют",
		"error_parse_Convert":           "Помилка при парсі курсів валют для ",
		"error_save_Convert":            "Помилка при збереженні конвертації",
		"error_find_rates":              "Не знайдені курси валют",
		"error_handle_currency":         "Помилка валюта не підтримується використовуй USD або EUR",
		"error_parse_handleCurrency":    "Помилка при парсингу валют",
		"error_scan_history":            "Помилка при скануванні історії",
		"error_receive_history":         "Помилка при отриманні історії",
		"error_receive_stats":           "Помилка при отриманні статистики",
		"log_stats":                     "Отримано статистику:",
		"error_clear_history":           "Помилка при очищенні історії",
		"error_lang":                    "Помилка мова не підтимується використовуйте en або uk",
		"error_save_language":           "Помилка при збереженні мови",
		"save_conversation":             "Збережено конвертацію:",
		"error_format_row":              "Невірний формат рядка:",
		"getRates_return":               "GetRates повернув для:",
		"course_not_found":              "Курс не знайдено",
		"row_is_empty":                  "Рядок порожній",
		"row_is_invalid":                "Невірний формат рядка",
		"course_is_invalid":             "Невірний формат курсу для",
		"course_parse_error":            "Помилка парсингу при купівлі для",
		"course_parse_error_sell":       "Помилка парсингу при продажі для",
		"buy":                           "Купівля",
		"sell":                          "Продаж",
		"error_format_course":           "Валюта не відповідає формату",
	},
}

type CurrencyHandler struct {
	Service *service.CurrencyService
}

func NewCurrencyHandler(service *service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{service}
}

func (h *CurrencyHandler) HandleStart(chatID int64) string {
	return h.getTranslation(chatID, "start")
}

func (h *CurrencyHandler) HandleRates(chatID int64, money string) (string, error) {

	if money != "USD" && money != "EUR" && money != "PLN" && money != "GBP" && money != "" {
		return "", fmt.Errorf("%s", h.getTranslation(chatID, "error_invalid_currencyRATES"))
	}
	rates, err := h.Service.GetRates(money)
	if err != nil {
		return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_receiving_ratesRATES"), err)
	}
	return rates, nil
}

func (h *CurrencyHandler) HandleHelp(chatID int64) string {
	return h.getTranslation(chatID, "help")
}

func (h *CurrencyHandler) HandleAbout(chatID int64) string {
	return h.getTranslation(chatID, "about")
}

func (h *CurrencyHandler) HandleConvert(chatID int64, amount float64, from, to string) (string, error) {
	//check value
	if (from != "USD" && from != "EUR" && from != "UAN") || (to != "USD" && to != "EUR" && to != "UAN") || from == to {
		return "", fmt.Errorf("%s", h.getTranslation(chatID, "error_invalid_currency-handle"))
	}
	var course string
	var err error
	if to == "UAN" {

		// Get course
		course, err = h.Service.GetRates(from)
		if err != nil {
			log.Println("%s: %v", h.getTranslation(chatID, "error_invalid_currencyRATES"), err)
			return "", err
		}
		log.Printf("GetRates повернув для %s: %q", from, course)
		//Parse course
		buy, sell, err := h.parseRate(course, from, chatID)
		if err != nil {
			log.Println("%s %s: %v", h.getTranslation(chatID, "error_parse_convert"), from, err)
			return "", fmt.Errorf("%s %s: %v", h.getTranslation(chatID, "error_parse_convert"), from, err)
		}

		log.Printf("Курс для %s: buy= %.2f, sell=%.2f", from, buy, sell)
		result := amount * (buy + sell) / 2

		errDb := h.Service.SaveConversion(chatID, amount, from, to, result, time.Now())
		if errDb != nil {
			log.Printf("%s: %v", h.getTranslation(chatID, "error_save_conversion"), errDb)
			return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_save_conversion"), errDb)
		}

		log.Printf("%s: %.2f %s = %.2f %s", h.getTranslation(chatID, "save_conversation"), amount, from, result, to)
		return fmt.Sprintf("%.2f %s = %.2f %s", amount, from, result, to), nil

	} else if from == "UAN" {
		//Get course
		course, err = h.Service.GetRates(to)
		if err != nil {
			log.Println("%s: %v", h.getTranslation(chatID, "error_invalid_currencyRATES"), err)
			return "", err
		}
		log.Printf("%s %s: %q", h.getTranslation(chatID, "getRates_return"), to, course)

		//Parse course
		buy, sell, err := h.parseRate(course, to, chatID)
		if err != nil {
			log.Println("%s %s: %v", h.getTranslation(chatID, "error_parse_convert"), from, err)
			return "", fmt.Errorf("%s %s: %v", h.getTranslation(chatID, "error_parse_convert"), from, err)
		}
		log.Printf("Курс для %s: buy= %.2f, sell=%.2f", to, buy, sell)
		rate := (buy + sell) / 2
		result := amount / rate

		errDb := h.Service.SaveConversion(chatID, amount, from, to, result, time.Now())
		if errDb != nil {
			log.Printf("chatID %d: %s: %v", chatID, h.getTranslation(chatID, "error_save_conversion"), errDb)
			return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_save_conversion"), errDb)
		}

		log.Printf("%s: %.2f %s = %.2f %s", h.getTranslation(chatID, h.getTranslation(chatID, "save_conversation")), amount, from, result, to)

		return fmt.Sprintf("%.2f %s = %.2f %s", amount, from, result, to), nil
	}

	//Debug with return
	log.Printf("%s: %q", h.getTranslation(chatID, "error_find_rates"), from, course)

	return "", fmt.Errorf("%s ", h.getTranslation(chatID, "error_find_rates"), from)

}

func (h *CurrencyHandler) HandleCurrency(chatID int64, money string) (string, error) {
	if money != "USD" && money != "EUR" {
		return "", fmt.Errorf("%s", h.getTranslation(chatID, "error_find_rates"), money)
	}
	rates, err := h.Service.GetRates(money)
	if err != nil {
		return "", fmt.Errorf("%s %s : %v", h.getTranslation(chatID, "error_receiving_ratesRATES"), money, err)
	}

	buy, sell, err := h.parseRate(rates, money, chatID)
	if err != nil {
		log.Printf("%s %s: %v", h.getTranslation(chatID, "error_parse_handleCurrency"), money, err)
		return "", fmt.Errorf("%s %s: %v", h.getTranslation(chatID, "error_parse_Convert"), money, err)
	}

	log.Printf("Курс для %s : buy=%.2f, sell=%.2f", money, buy, sell)
	return rates, nil
}

func (h *CurrencyHandler) HandleHistory(chatID int64) (string, error) {
	rows, err := h.Service.GetHistory(chatID)
	if err != nil {
		log.Printf("chatID %d: %v", chatID, err)
		return "", fmt.Errorf("%s %v", h.getTranslation(chatID, "error_receive_history"), err)
	}

	defer rows.Close()

	var amount float64
	var from_currency string
	var to_currency string
	var result float64
	var timestamp time.Time
	records := []string{}
	for rows.Next() {
		if errRows := rows.Scan(&amount, &from_currency, &to_currency, &result, &timestamp); errRows != nil {
			log.Printf("%s: %v", h.getTranslation(chatID, "error_scan_history"), errRows)
			return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_scan_history"), errRows)
		}
		//above must be 100.00 USD = 4150.00 UAN (2025-07-19 15:24:00)
		record := fmt.Sprintf("%.2f %s = %.2f %s (%s)", amount, from_currency, result, to_currency, timestamp.Format("2006-01-02 15:04:05"))
		records = append(records, record)
	}
	if len(records) == 0 {
		log.Println(" %s", h.getTranslation(chatID, "history_empty"))
		return h.getTranslation(chatID, "history_empty"), nil
	}
	return fmt.Sprintf("%s:\n%s", h.getTranslation(chatID, "history"), strings.Join(records, "\n")), nil

}

func (h *CurrencyHandler) HandleStats(chatID int64) (string, error) {
	log.Printf("chatID %d: Викликано /stats", chatID)
	records, isEmpty, err := h.Service.GetStats(chatID)
	if err != nil {
		log.Printf("%s: %v", h.getTranslation(chatID, "error_receive_stats"), err)
		return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_receive_stats"), err)
	}

	if isEmpty {
		log.Printf("chatID %d: %s", chatID, h.getTranslation(chatID, "stats_empty"))
		return h.getTranslation(chatID, "stats_empty"), nil
	}

	log.Printf("%s: %q", h.getTranslation(chatID, "log_stats"), strings.Join(records, "\n"))
	return fmt.Sprintf("%s\n%s", h.getTranslation(chatID, "stats"), strings.Join(records, "\n")), nil
}

func (h *CurrencyHandler) HandleClearHistory(chatID int64) (string, error) {
	err := h.Service.ClearHistory(chatID)
	if err != nil {
		log.Printf("%s: %v", h.getTranslation(chatID, "error_clear_history"), err)
		return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_clear_history"), err)
	}
	log.Printf("%s", h.getTranslation(chatID, "clear_history"))
	return h.getTranslation(chatID, "clear_history"), nil
}

func (h *CurrencyHandler) HandleLang(chatID int64, lang string) (string, error) {
	if lang != "en" && lang != "uk" {
		log.Printf("%s,", h.getTranslation(chatID, h.getTranslation(chatID, "error_lang")), lang)
		return "", fmt.Errorf("%s: %s", h.getTranslation(chatID, h.getTranslation(chatID, "error_lang")), lang)
	}

	text := h.getTranslation(chatID, "lang_success")
	err := h.Service.SetLanguage(chatID, lang)
	if err != nil {
		log.Printf("%s: %v", h.getTranslation(chatID, "error_save_language"), err)
		return "", fmt.Errorf("%s: %v", h.getTranslation(chatID, "error_save_language"), err)
	}

	log.Printf("%s", h.getTranslation(chatID, "lang_success"))
	return text, nil

}

func (h *CurrencyHandler) getTranslation(chatID int64, key string) string {
	lang := h.Service.GetLanguage(chatID)
	if trans, ok := translations[lang][key]; ok {
		return trans
	}
	return translations["uk"][key]
}

func (h *CurrencyHandler) parseRate(course, currency string, chatID int64) (buy, sell float64, err error) {
	if course == "" {
		log.Printf("%s : %s", h.getTranslation(chatID, "course_not_found"), currency)
		return 0, 0, fmt.Errorf("%s : %s", h.getTranslation(chatID, "course_not_found"), currency)
	}
	lines := strings.Split(strings.TrimSpace(course), "\n")
	for _, line := range lines {
		log.Printf("Обробка рядка: %q", line)
		if line == "" {
			log.Printf("%s %s", h.getTranslation(chatID, "row_is_empty"), currency)
			return 0, 0, fmt.Errorf("%s %s", h.getTranslation(chatID, "row_is_empty"), currency)
		}
		if strings.Contains(line, "*"+currency+"/UAN*") {
			line = strings.TrimSpace(line)
			line = strings.ReplaceAll(line, "\r", "")
			line = strings.ReplaceAll(line, "  ", " ")

			parts := strings.Split(line, ": ")
			if len(parts) != 2 {
				log.Printf("%s %s: %q", h.getTranslation(chatID, "row_is_invalid"), currency, line)
				return 0, 0, fmt.Errorf("%s %s: %q", h.getTranslation(chatID, "row_is_invalid"), currency, line)
			}

			if !strings.HasPrefix(parts[0], "*"+currency+"/UAN*") {
				log.Printf("%s: %q", h.getTranslation(chatID, "error_format_course"), parts[0])
				return 0, 0, fmt.Errorf("%s: %q", h.getTranslation(chatID, "error_format_course"), parts[0])
			}

			rateParts := strings.Split(parts[1], ", ")
			if len(rateParts) != 2 {
				log.Printf("%s %s: %q", h.getTranslation(chatID, "course_is_invalid"), currency, parts[1])
				return 0, 0, fmt.Errorf("%s %s: %q", h.getTranslation(chatID, "course_is_invalid"), currency, parts[1])
			}

			buyStr := strings.TrimPrefix(rateParts[0], h.getTranslation(chatID, "buy")+" ")
			log.Printf("Translation buy: %q", h.getTranslation(chatID, "buy"))
			buy, err := strconv.ParseFloat(buyStr, 64)
			if err != nil {
				log.Printf("%s %s: %v, str: %q", h.getTranslation(chatID, "course_parse_error"), currency, err, buyStr)
				return 0, 0, fmt.Errorf("%s %s: %v, str: %q", h.getTranslation(chatID, "course_parse_error"), currency, err, buyStr)
			}

			sellStr := strings.TrimPrefix(rateParts[1], h.getTranslation(chatID, "sell")+" ")
			sell, err := strconv.ParseFloat(sellStr, 64)
			if err != nil {
				log.Printf("%s %s: %v, str: %q", h.getTranslation(chatID, "course_parse_error_sell"), currency, err, sellStr)
				return 0, 0, fmt.Errorf("%s %s: %v, str: %q", h.getTranslation(chatID, "course_parse_error_sell"), currency, err, sellStr)
			}

			if buy == 0.0 && sell == 0.0 {
				return 0, 0, fmt.Errorf("%s: %s", h.getTranslation(chatID, "course_not_found"), currency)
			}
			return buy, sell, nil
		}

	}
	return 0, 0, fmt.Errorf("%s %s", h.getTranslation(chatID, "course_not_found"), currency)
}
