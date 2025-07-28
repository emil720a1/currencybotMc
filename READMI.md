# CurrencyBot

## Description
CurrencyBot is a Telegram bot that converts currencies in real time. It supports currencies like USD, EUR, and UAN. Making it easy for user to convert currencies and calculate exchange rates.

## Installation
1. Ensure you have GO (version 1.18 or higher) installed.
2. Clone the repository: ``
3. Navigate to the folder `cd currencybot `.
4. Install dependepcies: `go mod tidy`.
5. Configure the bot (see Configuration section).
6. Run the bot: `go run main.go`.

## Configuration
1. Create a `config.yaml` file in the root folder with this content:
```yaml
bot_token: "YOUR_BOT_TOKEN_FROM_TELEGRAM"
db_host: "localhost"
db_user: "your_username"
db_pass: "your_password"
db_name: "currency_bot"

## Usage
- Add the bot to Telegram and type '/start'.
- Use the command `/convert <amount> <from_currency> to <to_currency>` (e.g., `/convert 100 USD to UAN`).
- The bot will return the conversion result and save it to history

## Features
- Converts currencies (e.g., USD to UAN, EUR to UAN, UAN to USD, UAN to EUR).
- Supports multiple languages (English, Ukrainian).
- Saves conversion history in MySQL.
- Handles large amounts and reverse conversions.

## Testing
1. Run the bot with a test 'course' (e.g., `"*USD/UAN*: Buy 41.50, Sell 42.01"`).
2. Test commands:
    - `convert 100 USD to UAN` (expected: 4175.50).
    - `convert 3 UA to USD` (expected error).
3. Check logs (`bot.log`) and database for verification.    

## Contributing
1. Fork the repository
2. Create a branch: `git checkout -b feature/new-feature`.
3. Make changes and commit: `git commit -m "Add new feature"`.
4. Submit a pull request.

## License
This project is licensed under the MIT License, See the `LICENSE` file for details.

## Author
- [emil720a.feat[McOleg]](