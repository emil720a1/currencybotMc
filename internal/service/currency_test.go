package service

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestGetStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Помилка створення mock: ", err)
	}
	defer db.Close()
	service := CurrencyService{db: db}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT .*FROM conversions.*chat_id").
			WithArgs(927118456).
			WillReturnRows(
				sqlmock.NewRows([]string{"from_currency", "to_currency", "count", "avg"}).AddRow("USD", "UAN", 5, 150.0).AddRow("EUR", "UAN", 2, 200.0))

		result, errResult := Service.GetStats(927118456)
		expected := "Статистика:\nUSD -> UAN: 5 конвертацій, середня сума: 150.00 USD\nEUR -> UAN: 2 конвертацій, середня сума: 200.00 EUR"
		if errResult != nil {
			t.Logf("Отриманий результат: %q", result)
			t.Fatal("Помилка при обробці результатів: ", errResult)
		}
		if result != expected {
			t.Logf("Отриманий результат: %q", result)
			t.Errorf("Очікувалося %q, отримано %q", expected, result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Logf("Отриманий результат: %q", result)
			t.Errorf("Очікування моку не виконані: %v", err)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		mock.ExpectQuery("SELECT .*FROM conversions.*chat_id").
			WithArgs(927118456).
			WillReturnRows(sqlmock.NewRows([]string{
				"from_currency", "to_currency", "count", "avg",
			}))
		result, errResult := Service.GetStats(927118456)
		expected := "Статистика відсутня"
		if errResult != nil {
			t.Errorf("Очікувалось err == nil, отримано %v", errResult)
		}
		if result != expected {
			t.Logf("Отриманий результат: %q", result)
			t.Errorf("Очікувалося %q, отримано %q", expected, result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Очікування моку не виконані: %v", err)
		}
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .*FROM conversions.*chat_id").
			WithArgs(927118456).
			WillReturnError(fmt.Errorf("db error"))
		result, errResult := service.GetStats(927118456)
		if errResult == nil {
			t.Error("Очікувалась помилка, отримано err == nil", errResult)
		}
		if result != "" {
			t.Errorf("Очікувалося result == %q, отримано %q", "", result)
		}
		if errResult.Error() != "db error" {
			t.Logf("Отримана помилка: %v", errResult)
			t.Errorf("Очікувалось err.Error() = %q, отримано %q", "db error", errResult.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Очікування моку не виконані: %v", err)
		}
	})

}
