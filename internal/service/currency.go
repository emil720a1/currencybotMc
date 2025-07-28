package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CurrencyService struct {
	db           *sql.DB
	RatesCache   string
	CacheTime    time.Time
	CacheMu      sync.Mutex
	langMap      map[int64]string
	translations map[string]map[string]string
}

func NewCurrencyService(db *sql.DB) *CurrencyService {
	return &CurrencyService{
		db:         db,
		RatesCache: "",
		CacheTime:  time.Time{},
		langMap:    make(map[int64]string),
	}
}

type Rate struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	RateBuy       float64 `json:"rateBuy"`
	RateSell      float64 `json:"rateSell"`
}

func (s *CurrencyService) GetRates(money string) (string, error) {
	s.CacheMu.Lock()
	defer s.CacheMu.Unlock()
	if time.Since(s.CacheTime) < 5*time.Minute && s.RatesCache != "" {
		return filterRates(s.RatesCache, money), nil
	}

	time.Sleep(100 * time.Millisecond)
	resp, err := http.Get("https://api.monobank.ua/bank/currency")
	defer resp.Body.Close()
	if err != nil {
		log.Println("Помилка при отриманні курсу валют: ", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		log.Println("неправильний статус код відповіді: ", resp.StatusCode)
		return "", fmt.Errorf("неправильний статус код відповіді: %d", resp.StatusCode)
	}

	var rates []Rate
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&rates)
	if err != nil {
		log.Println("Помилка при декодуванні json: ", err)
		return "", err
	}
	log.Printf("API повернуло: %d валют", len(rates))
	var result string
	for _, rate := range rates {
		log.Printf("API повернуло: currencyCodeA=%d, currencyCodeB=%d, rateBuy=%.2f, rateSell=%.2f", rate.CurrencyCodeA, rate.CurrencyCodeB, rate.RateBuy, rate.RateSell)
		if rate.CurrencyCodeB == 980 {

			if rate.CurrencyCodeA == 840 {
				result += fmt.Sprintf("*USD/UAN*: Купівля %.2f, Продаж %.2f\n", rate.RateBuy, rate.RateSell)

			} else if rate.CurrencyCodeA == 978 {
				result += fmt.Sprintf("*EUR/UAN*: Купівля %.2f, Продаж %.2f\n", rate.RateBuy, rate.RateSell)
			} else if money == "" {
				if rate.CurrencyCodeA == 840 {
					result += fmt.Sprintf("*USD/UAN*: Купівля %.2f, Продаж %.2f\n", rate.RateBuy, rate.RateSell)
				} else if rate.CurrencyCodeA == 978 {
					result += fmt.Sprintf("*EUR/UAN*: Купівля %.2f, Продаж %.2f\n", rate.RateBuy, rate.RateSell)
				}
			}
		}
	}
	if result == "" {
		log.Printf("Курси для %s не знайдені", money)
		return "", fmt.Errorf("Курси для вказаної валюти не знайдені")
	}

	s.RatesCache = result
	s.CacheTime = time.Now()

	return filterRates(result, money), nil
}

func filterRates(cached, money string) string {
	if money == "" {
		return cached
	}

	lines := strings.Split(cached, "\n")
	for _, line := range lines {
		log.Printf("Обробка рядка: %s", line)
		if strings.Contains(line, "*"+money+"/UAN*") {
			return line + "\n"
		}
		if line == "" {
			log.Printf("Порожній рядок у filterRates для %s", money)
			continue
		}
	}
	log.Printf("FilterRates: не знайдено курс для %s у кеші: %q", money, cached)
	return ""
}

func (s *CurrencyService) SaveConversion(chatID int64, amount float64, from, to string, result float64, timestamp time.Time) error {
	_, errDb := s.db.Exec("INSERT INTO conversions (chat_id, amount, from_currency, to_currency, result, timestamp) VALUES (?, ?, ?, ?, ?, ?)", chatID, amount, from, to, result, timestamp)
	if errDb != nil {
		return errDb
	}
	return nil
}

func (s *CurrencyService) GetHistory(chatID int64) (*sql.Rows, error) {
	rows, err := s.db.Query("SELECT amount, from_currency, to_currency, result, timestamp FROM conversions WHERE chat_id = ? ORDER BY timestamp DESC LIMIT 5", chatID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *CurrencyService) GetStats(chatID int64) ([]string, bool, error) {
	rows, err := s.db.Query("SELECT from_currency, to_currency, COUNT(*), AVG(amount) FROM conversions WHERE chat_id = ? GROUP BY from_currency, to_currency", chatID)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	var from_currency, to_currency string
	var count int
	var avg float64
	records := []string{}
	for rows.Next() {
		if errRows := rows.Scan(&from_currency, &to_currency, &count, &avg); errRows != nil {
			return nil, false, errRows
		}
		res := fmt.Sprintf("%s -> %s: %d конвертацій, середня сума: %.2f %s", from_currency, to_currency, count, avg, from_currency)
		records = append(records, res)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	return records, len(records) == 0, nil
}

func (s *CurrencyService) ClearHistory(chatID int64) error {
	_, errDb := s.db.Exec("DELETE FROM conversions WHERE chat_id = ?", chatID)
	if errDb != nil {
		return errDb
	}
	return nil
}

func (s *CurrencyService) SetLanguage(chatID int64, lang string) error {
	lang = strings.ToLower(lang)
	if lang == "" {
		return fmt.Errorf("Мова не може бути порожньою")
	}

	if lang != "uk" && lang != "en" {
		return fmt.Errorf("Непідтримувана мова, використовуйте uk або en")
	}
	log.Printf("SetLanguage: %d, %s", chatID, lang)
	s.langMap[chatID] = lang
	return nil
}

func (s *CurrencyService) GetLanguage(chatID int64) string {
	//return language for chatID or "uk" for default
	if lang, ok := s.langMap[chatID]; ok {
		return lang
	}
	return "uk"

}
