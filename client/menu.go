package client

// func (c *Client) isStartMessage(update *tgbotapi.Update) bool {
// 	return update.Message != nil && update.Message.Text == "/start"
// }

// func (c *Client) Helloer(update *tgbotapi.Update) {
// 	str := fmt.Sprint("Привет, ", update.Message.From.FirstName, " ", emoji.WavingHand)

// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)

// 	_, err := c.bot.Send(msg)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Что хочешь сделать?")
// 	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData("добавить товар", "add"),
// 			tgbotapi.NewInlineKeyboardButtonData("списать товар", "/remove"),
// 			tgbotapi.NewInlineKeyboardButtonData("посмотреть список", "/list"),
// 		),
// 	)

// 	msg.ReplyMarkup = numericKeyboard
// 	_, err = c.bot.Send(msg)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func (c *Client) MsgAddingPosition(update *tgbotapi.Update) {
// 	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выбрали добавление позиции. Введите наименование")
// 	_, err := c.bot.Send(msg)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func (c *Client) printAddingPosition(update *tgbotapi.Update) {

//		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выбрали добавление позиции. Введите наименование")
//		_, err := c.bot.Send(msg)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
// func (c *Client) addingPosition(update *tgbotapi.Update) {
// 	var l []string
// 	l = append(l, update.Message.Text)
// 	fmt.Println(l)
// }
