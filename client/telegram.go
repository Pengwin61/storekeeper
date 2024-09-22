package client

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	UPDATE_CONFIG_TIMEOUT = 60
)

type UserState struct {
	Action      string
	Product     string
	Description string
	Count       int
	Price       float64
}

var userStates = make(map[int64]*UserState)

type Client struct {
	bot          *tgbotapi.BotAPI
	updateConfig *tgbotapi.UpdateConfig
}

func NewClient(tgToken string) *Client {
	b, uc := initClient(tgToken)
	return &Client{bot: b, updateConfig: uc}
}

func initClient(tgToken string) (*tgbotapi.BotAPI, *tgbotapi.UpdateConfig) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	return bot, &updateConfig
}

func (c *Client) Start(db *sql.DB) {

	for update := range c.bot.GetUpdatesChan(*c.updateConfig) {
		if update.Message != nil {
			c.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			c.handleCallback(update.CallbackQuery, db)
		}
	}
}

func (c *Client) handleMessage(message *tgbotapi.Message) {
	fmt.Println("ID:", message.Chat.ID)

	if c.isAdmins(message.Chat.ID) {
		// c.sendMsg(message.Chat.ID, "Панель администратора")
		if userStates[message.Chat.ID] == nil {
			userStates[message.Chat.ID] = &UserState{}
			c.adminShowMainMenu(message.Chat.ID)
		}
		c.handleAdditions(message.Chat.ID, message.Text)
	} else {
		// c.sendMsg(message.Chat.ID, "Пользовательская панель")
		c.userShowMainMenu(message.Chat.ID)
	}

}

func (c *Client) userShowMainMenu(chatID int64) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть все", "view_all"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Что вы хотите сделать?")
	msg.ReplyMarkup = inlineKeyboard
	c.sendMsgKB(msg)
}

func (c *Client) adminShowMainMenu(chatID int64) {
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить", "add"),
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть все", "view_all"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Что вы хотите сделать?")
	msg.ReplyMarkup = inlineKeyboard
	c.sendMsgKB(msg)
}

func (c *Client) handleCallback(callback *tgbotapi.CallbackQuery, db *sql.DB) {
	userID := callback.Message.Chat.ID
	state := userStates[userID]

	if state == nil {
		state = &UserState{}
		return
	}

	switch callback.Data {
	case "add":
		state.Product = ""
		state.Description = ""
		state.Price = 0
		state.Action = "waiting_for_product"
	case "view_all":
		res, err := c.getAllProducts(db)
		if err != nil {
			fmt.Println(err)
		}
		c.sendMsg(userID, res)
	case "cancel":
		if state.Action == "waiting_for_confirm" {
			c.sendMsg(userID, "Операция отменена.")
			state.Action = ""
			state.Product = ""
			state.Description = ""
			state.Price = 0
		} else {
			c.sendMsg(userID, "Неверная операция.")
		}

	case "confirm":
		if state.Action == "waiting_for_confirm" {
			_, err := db.Exec("INSERT INTO products (name, description, count, price) VALUES (?, ?, ?, ?)", state.Product, state.Description, state.Count, state.Price)
			if err != nil {
				c.sendMsg(userID, "Ошибка при сохранении данных.")
			}
			c.sendMsg(userID, "Товар добавлен")
			state.Action = ""
		} else {
			c.sendMsg(userID, "Неверная операция.")
		}
	}

	switch state.Action {
	case "":
		return
	case "waiting_for_product":
		c.sendMsg(userID, "Введите наименование:")
	default:
		c.adminShowMainMenu(userID)
	}
}

func (c *Client) handleAdditions(userID int64, input string) {
	state := userStates[userID]

	switch state.Action {
	case "waiting_for_product":
		state.Product = input
		c.sendMsg(userID, "Введите описание:")
		state.Action = "waiting_for_description"

	case "waiting_for_description":
		state.Description = input
		c.sendMsg(userID, "Введите количество:")
		state.Action = "waiting_for_count"

	case "waiting_for_count":
		count, err := strconv.Atoi(input)
		if err != nil {
			c.sendMsg(userID, "Введите корректное количество:")
			return
		}
		if count > 0 {
			state.Count = count
			c.sendMsg(userID, "Введите стоимость:")
			state.Action = "waiting_for_price"
		} else {
			c.sendMsg(userID, "Введите корректное количество:")
		}

	case "waiting_for_price":
		price, err := strconv.ParseFloat(input, 64)
		if err != nil {
			c.sendMsg(userID, "Введите корректную стоимость:")
			return
		}
		if price > 0 {
			state.Price = price
			// Создаем итоговое сообщение
			summary := fmt.Sprintf("Вы добавляете товар:\nНазвание: %s\nОписание: %s\nКол-во: %d\nСтоимость: %.2f рублей\nПодтверждаете данные? (да/нет)", state.Product, state.Description, state.Count, state.Price)
			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm"),
					tgbotapi.NewInlineKeyboardButtonData("Отклонить", "cancel"),
				),
			)
			msg := tgbotapi.NewMessage(userID, summary)
			msg.ReplyMarkup = inlineKeyboard
			c.sendMsgKB(msg)
			state.Action = "waiting_for_confirm"
		} else {
			c.sendMsg(userID, "Введите положитеную стоимость:")
			return
		}
	}
}

func (c *Client) getAllProducts(db *sql.DB) (string, error) {
	rows, err := db.Query("SELECT name, description, count, price FROM products")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result strings.Builder
	for rows.Next() {
		var name, description string
		var price float64
		var count int
		if err := rows.Scan(&name, &description, &count, &price); err != nil {
			return "", err
		}
		_, err := result.WriteString(fmt.Sprintf("Название: %s\nОписание: %s\nКол-во: %d\nСтоимость: %.2f\n\n", name, description, count, price))
		if err != nil {
			return "", err
		}
	}
	if result.Len() == 0 {
		return "Нет позиций в базе данных.", nil
	}
	return result.String(), nil
}

func (c *Client) sendMsg(userID int64, msg string) {
	_, err := c.bot.Send(tgbotapi.NewMessage(userID, msg))
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) sendMsgKB(msg tgbotapi.MessageConfig) {
	_, err := c.bot.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) isAdmins(userID int64) bool {
	admins := make([]int64, 0)
	admins = append(admins, 329159577)
	admins = append(admins, 1188924200)

	for _, admin := range admins {
		if userID == admin {
			return true
		}
	}

	return false
}
