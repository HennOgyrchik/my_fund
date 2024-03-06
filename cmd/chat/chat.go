package chat

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"project1/cmd/service"
	"project1/db"
	"regexp"
)

type Chat struct {
	username string
	chatId   int64
	msg      *tgbotapi.Message
	*service.Service
}

func (c *Chat) test() error { // 3 попытки на отправку
	var err error

	for i := 0; i < 3; i++ {
		if err = c.send(tgbotapi.NewMessage(c.chatId, "тестовая кнопка")); err == nil {
			return err
		}
	}
	return err
}

func NewChat(username string, chatId int64, service *service.Service) *Chat {
	return &Chat{
		username: username,
		chatId:   chatId,
		msg:      nil,
		Service:  service,
	}
}

func (c *Chat) GetChatId() int64 {
	return c.chatId
}

func (c *Chat) GetMessage() *tgbotapi.Message {
	return c.msg
}

// send 3 попытки на отправку, иначе ошибка
func (c *Chat) send(message tgbotapi.MessageConfig) error {
	var err error

	for i := 0; i < 3; i++ {
		if _, err = c.Service.GetBot().Send(message); err == nil {
			return nil
		}
	}
	return err
}

func (c *Chat) CommandSwitcher(query string) {
	var paymentPat = regexp.MustCompile(`^оплатить\s\d*.`)
	var rejectionPat = regexp.MustCompile(`^отказ\s\d*.`)
	var waitingPat = regexp.MustCompile(`^ожидание\s\d*.`)
	var acceptPat = regexp.MustCompile(`^подтвердить\s\d*.`)

	switch cmd := query; {
	case cmd == "start":
		c.startMenu()
	case cmd == "menu":
		c.showMenu()
	case cmd == "создать":
		//c.confirmationCreationNewFund()
	case cmd == "присоединиться":
		//c.join()
	case cmd == "создать новый фонд":
		//c.creatingNewFund()
	case cmd == "баланс":
		//c.showBalance()
	case cmd == "test":
		c.test()
	case cmd == "участники":
		//c.getMembers()
	case cmd == "новый сбор":
		//c.createCashCollection()
	case cmd == "новое списание":
		//c.createDebitingFunds()
	case paymentPat.MatchString(cmd): // оплата
		//cashCollectionId, err := strconv.Atoi(strings.Split(cmd, " ")[1])
		//if err != nil {
		//	c.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
		//	return
		//}
		//c.payment(cashCollectionId)
	case acceptPat.MatchString(cmd): // подтверждение оплаты
		//idTransaction, err := strconv.Atoi(strings.Split(cmd, " ")[1])
		//if err != nil {
		//	c.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
		//	return
		//}
		//c.changeStatusOfTransaction(idTransaction, "подтвержден")
	case waitingPat.MatchString(cmd): // ожидание оплаты
		//idTransaction, err := strconv.Atoi(strings.Split(cmd, " ")[1])
		//if err != nil {
		//	c.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
		//	return
		//}
		//c.changeStatusOfTransaction(idTransaction, "ожидание")
	case rejectionPat.MatchString(cmd): // отказ оплаты
		//idTransaction, err := strconv.Atoi(strings.Split(cmd, " ")[1])
		//if err != nil {
		//	c.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
		//	return
		//}
		//c.changeStatusOfTransaction(idTransaction, "отказ")
	default:
		_, _ = c.GetBot().Send(tgbotapi.NewMessage(c.chatId, "Я не знаю такую команду"))
	}

}

func (c *Chat) startMenu() {
	var startKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать фонд", "создать"),
			tgbotapi.NewInlineKeyboardButtonData("Присоединиться", "присоединиться"),
			tgbotapi.NewInlineKeyboardButtonData("Тест", "test"),
		),
	)

	msg := tgbotapi.NewMessage(c.chatId, "Приветствую! Выберите один из вариантов")
	msg.ReplyMarkup = &startKeyboard

	err := c.send(msg)
	fmt.Println(err)
}

func (c *Chat) showMenu() {
	ok, err := db.IsMember(c.chatId)
	if err != nil {
		err = c.send(tgbotapi.NewMessage(c.chatId, "Произошла ошибка. Попробуйте еще раз."))
		return
	}
	if !ok {
		if err = c.send(tgbotapi.NewMessage(c.chatId, "Вы не являетесь участником фонда. Создайте новый фонд или присоединитесь к существующему.")); err != nil {
			return
		}
		c.startMenu()
		return
	}

	var menuKeyboard = tgbotapi.NewInlineKeyboardMarkup( //меню для обычного пользователя
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Баланс", "баланс"),
			tgbotapi.NewInlineKeyboardButtonData("Оплатить", "1"), // реализовать
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Покинуть фонд", "3"), // реализовать
		),
	)

	msg := tgbotapi.NewMessage(c.chatId, "Приветствую! Выберите один из вариантов")

	ok, err = db.IsAdmin(c.chatId)
	if err != nil {
		err = c.send(tgbotapi.NewMessage(c.chatId, "Произошла ошибка. Попробуйте еще раз."))
		return
	}

	if ok { // если админ, то дополнить меню
		menuKeyboard.InlineKeyboard = append(menuKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Новый сбор", "новый сбор"),
				tgbotapi.NewInlineKeyboardButtonData("Новое списание", "новое списание")),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Участники", "участники"),
				tgbotapi.NewInlineKeyboardButtonData("Статистика", "2"))) // реализовать
	}

	msg.ReplyMarkup = &menuKeyboard
	_ = c.send(msg)
}

//
//func (r *response) downloadAttachment(fileId string) (fileName string, err error) {
//	_, err = r.bot.GetFile(tgbotapi.FileConfig{FileID: fileId})
//	if err != nil {
//		return
//	}
//
//	pathFile, err := r.bot.GetFileDirectURL(fileId)
//	if err != nil {
//		return
//	}
//
//	resp, err := http.Get(pathFile)
//	defer resp.Body.Close()
//	if err != nil {
//		return
//	}
//
//	fileName = strconv.FormatInt(r.chatId, 10) + "_" + path.Base(pathFile)
//	ok, err := ftp.StoreFile(fileName, resp.Body)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Print(ok)
//
//	return
//}
//
//func (r *response) createDebitingFunds() {
//	sum, err := r.getFloatFromUser("Введите сумму списания.")
//	if err != nil {
//		return
//	}
//
//	msg := tgbotapi.NewMessage(r.chatId, "Укажите причину списания")
//	if _, err = r.bot.Send(msg); err != nil {
//		return
//	}
//
//	answer, err := r.waitingResponse("text")
//	if err != nil {
//		return
//	}
//	purpose := answer.Text
//
//	tag, err := db.GetTag(r.chatId)
//	if err != nil {
//		return
//	}
//
//	msg.Text = "Прикрепите чек файлом"
//	if _, err = r.bot.Send(msg); err != nil {
//		return
//	}
//	////////////////////////////////////////////ожидание чека файлом или картинкой
//	answer, err = r.waitingResponse("attachment")
//	if err != nil {
//		return
//	}
//
//	var file string
//	if answer.Photo != nil {
//		file = answer.Photo[len(answer.Photo)-1].FileID
//
//	} else {
//		file = answer.Document.FileID
//	}
//	fileName, err := r.downloadAttachment(file)
//	if err != nil {
//		return
//	}
//
//	///////////////////////////Создание транзакции////////////////////////////////////////////////
//	ok, err := db.CreateDebitingFunds(r.chatId, tag, sum, fmt.Sprintf("Инициатор: %s", r.username), purpose, fileName)
//	if err != nil || !ok {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	///////////////////////////////////////////////////////////////////////////
//	msg.Text = "Списание проведено успешно."
//	_, _ = r.bot.Send(msg)
//
//}
//
//func (r *response) changeStatusOfTransaction(idTransaction int, status string) {
//	err := db.ChangeStatusTransaction(idTransaction, status)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	msg := tgbotapi.NewMessage(r.chatId, fmt.Sprintf("Статус оплаты: %s", status))
//	_, _ = r.bot.Send(msg)
//
//	r.paymentChangeStatusNotification(idTransaction)
//}
//
//func (r *response) paymentChangeStatusNotification(idTransaction int) {
//	status, _, _, memberId, _, err := db.InfoAboutTransaction(idTransaction)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	msg := tgbotapi.NewMessage(memberId, fmt.Sprintf("Статус оплаты изменен на: %s", status))
//	_, _ = r.bot.Send(msg)
//}
//
//func (r *response) payment(cashCollectionId int) {
//	target, _, err := db.InfoAboutCashCollection(cashCollectionId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	sum, err := r.getFloatFromUser("Введите сумму пополнения.")
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	if sum < target {
//		_, _ = r.bot.Send(tgbotapi.NewMessage(r.chatId, "Вы не можете оплатить сумму меньше необходимой."))
//		return
//	}
//
//	idTransaction, err := db.InsertInTransactions(cashCollectionId, sum, "пополнение", "ожидание", "", r.chatId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	msg := tgbotapi.NewMessage(r.chatId, "Ваша оплата добавлена в очередь на подтверждение")
//	_, _ = r.bot.Send(msg)
//	r.paymentNotification(idTransaction)
//}
//
//func (r *response) paymentNotification(idTransaction int) { //доделать
//	tag, err := db.GetTag(r.chatId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	adminId, err := db.GetAdminFund(tag)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	var okKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//		tgbotapi.NewInlineKeyboardRow(
//			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", fmt.Sprintf("подтвердить %d", idTransaction)),
//			tgbotapi.NewInlineKeyboardButtonData("Отказ", fmt.Sprintf("отказ %d", idTransaction)),
//			tgbotapi.NewInlineKeyboardButtonData("Ожидание", fmt.Sprintf("ожидание %d", idTransaction)),
//		),
//	)
//
//	_, _, _, memberId, sum, err := db.InfoAboutTransaction(idTransaction)
//
//	_, _, name, err := db.GetInfoAboutMember(memberId)
//
//	msg := tgbotapi.NewMessage(adminId, fmt.Sprintf("Подтвердите зачисление средств на счет фонда.\nСумма: %.2f\nОтправитель: %s", sum, name))
//	msg.ReplyMarkup = &okKeyboard
//	_, _ = r.bot.Send(msg)
//
//}
//
//func (r *response) getFloatFromUser(message string) (sum float64, err error) {
//	msg := tgbotapi.NewMessage(r.chatId, message)
//	if _, err = r.bot.Send(msg); err != nil {
//		return
//	}
//
//	for {
//		var answer *tgbotapi.Message
//
//		answer, err = r.waitingResponse("text")
//		if err != nil {
//			return
//		}
//
//		sum, err = strconv.ParseFloat(answer.Text, 64)
//		if err != nil {
//			r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//			continue
//		}
//		break
//	}
//	return
//}
//
//func (r *response) createCashCollection() {
//	sum, err := r.getFloatFromUser("Введите сумму сбора с одного участника.")
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	msg := tgbotapi.NewMessage(r.chatId, "Укажите назначение сбора")
//	if _, err = r.bot.Send(msg); err != nil {
//		return
//	}
//
//	answer, err := r.waitingResponse("text")
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	purpose := answer.Text
//
//	tag, err := db.GetTag(r.chatId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	id, err := db.CreateCashCollection(tag, sum, "открыт", fmt.Sprintf("Инициатор: %s", r.username), purpose, "")
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	msg.Text = "Сбор создан. Сообщение о сборе будет отправлено всем участникам."
//	_, _ = r.bot.Send(msg)
//
//	r.collectionNotification(id, tag)
//}
//
//func (r *response) collectionNotification(idCollection int, tagFund string) {
//	members, err := db.SelectMembers(tagFund)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	sum, purpose, err := db.InfoAboutCashCollection(idCollection)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	for _, value := range members {
//		var paymentKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//			tgbotapi.NewInlineKeyboardRow(
//				tgbotapi.NewInlineKeyboardButtonData("Оплатить", fmt.Sprintf("оплатить %d", idCollection)),
//			),
//		)
//		msg := tgbotapi.NewMessage(value, fmt.Sprintf("Иницирован новый сбор.\nСумма к оплате: %.2f\nНазначение: %s", sum, purpose))
//		msg.ReplyMarkup = &paymentKeyboard
//		_, _ = r.bot.Send(msg)
//	}
//}
//
//func (r *response) showBalance() {
//	tag, err := db.GetTag(r.chatId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	balance, err := db.ShowBalance(tag)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	msg := tgbotapi.NewMessage(r.chatId, fmt.Sprintf("Текущий баланс: %.2f руб", balance))
//	_, _ = r.bot.Send(msg)
//
//}
//
//
//func (r *response) join() {
//	msg := tgbotapi.NewMessage(r.chatId, "")
//	ok, err := db.IsMember(r.chatId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	if ok {
//		msg.Text = "Вы уже являетесь участником фонда"
//		_, _ = r.bot.Send(msg)
//		return
//	}
//	msg.Text = "Введите тег фонда. Если у вас нет тега, запросите его у администратора фонда."
//	if _, err = r.bot.Send(msg); err != nil {
//		return
//	}
//	answer, err := r.waitingResponse("text")
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//	}
//	tag := answer.Text
//
//	ok, err = db.ExistsFund(tag)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	if !ok {
//		msg.Text = "Фонд с таким тегом не найден."
//	} else {
//
//		name, err := r.getName()
//		if err != nil {
//			r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//			return
//		} else {
//			err = db.AddMember(tag, r.chatId, false, r.username, name)
//			if err != nil {
//				r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//				return
//			}
//			msg.Text = "Вы успешно присоединились к фонду."
//		}
//	}
//
//	_, _ = r.bot.Send(msg)
//}
//
//func (r *response) confirmationCreationNewFund() {
//	msg := tgbotapi.NewMessage(r.chatId, "")
//	ok, err := db.IsMember(r.chatId)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//	if ok {
//		msg.Text = "Вы уже являетесь участником фонда"
//	} else {
//		msg.Text = "Вы уверены, что хотите создать новый фонд?"
//		var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
//			tgbotapi.NewInlineKeyboardRow(
//				tgbotapi.NewInlineKeyboardButtonData("Да", "создать новый фонд"),
//				tgbotapi.NewInlineKeyboardButtonData("Нет", "start"),
//			),
//		)
//		msg.ReplyMarkup = numericKeyboard
//	}
//	_, _ = r.bot.Send(msg)
//}
//
//func (r *response) creatingNewFund() {
//	var err error
//	sum, err := r.getFloatFromUser("Введите начальную сумму фонда без указания валюты.")
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	var tag string
//	for i := 0; i < 10; i++ {
//		tag = newTag()
//
//		ok, err := db.DoesTagExist(tag)
//		if err != nil {
//			r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//			return
//		}
//		if !ok {
//			continue
//		}
//		break
//	}
//
//	err = db.CreateFund(tag, sum)
//	name, err := r.getName()
//
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	err = db.AddMember(tag, r.chatId, true, r.username, name)
//	if err != nil {
//		r.notificationAboutError("Произошла ошибка. Попробуйте еще раз.")
//		return
//	}
//
//	msg := tgbotapi.NewMessage(r.chatId, fmt.Sprintf("Новый фонд создан успешно! Присоединиться к фонду можно, используя тег: %s \nВнимание! Не показывайте этот тег посторонним людям.", tag))
//	_, _ = r.bot.Send(msg)
//
//}
//
//func (r *response) waitingResponse(obj string) (*tgbotapi.Message, error) {
//
//	waitingList[r.chatId] = make(chan *tgbotapi.Message)
//	defer func() {
//		close(waitingList[r.chatId])
//		delete(waitingList, r.chatId)
//	}()
//
//	var typeOfMessage string
//	var answer *tgbotapi.Message
//
//	for i := 0; i < 3; i++ {
//		answer = <-waitingList[r.chatId]
//		if answer.Photo != nil || answer.Document != nil {
//			typeOfMessage = "attachment"
//		} else {
//			typeOfMessage = "text"
//		}
//
//		if obj != typeOfMessage {
//			if i < 2 {
//				r.notificationAboutError(fmt.Sprintf("Вы ввели что-то не то. Количество доступных попыток: %d", 2-i))
//			}
//			continue
//		}
//		return answer, nil
//	}
//
//	return answer, errors.New("The number of attempts exceeded\n")
//}
//
//func newTag() string {
//	symbols := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
//	result := make([]byte, rand.Intn(5)+5)
//	for i := range result {
//		result[i] = symbols[rand.Intn(len(symbols))]
//	}
//	return string(result)
//}
//
//func (r *response) getName() (string, error) {
//	msg := tgbotapi.NewMessage(r.chatId, "Представьтесь, пожалуйста. Введите ФИО")
//	if _, err := r.bot.Send(msg); err != nil {
//		return "", err
//	}
//	answer, err := r.waitingResponse("text")
//	if err != nil {
//		return "", err
//	}
//	return answer.Text, nil
//}
//
//func (r *response) notificationAboutError(message string) {
//	if message == "" {
//		message = "Произошла ошибка. Попробуйте позже"
//	}
//
//	msg := tgbotapi.NewMessage(r.chatId, message)
//	_, _ = r.bot.Send(msg)
//	return
//}
//
//func (r *response) getMembers() {
//	tag, err := db.GetTag(r.chatId)
//	if err != nil {
//		r.notificationAboutError("")
//		return
//	}
//
//	id_members, err := db.SelectMembers(tag)
//	if err != nil {
//		r.notificationAboutError("")
//		return
//	}
//
//	message := "Список участников:\n"
//
//	for i, member := range id_members {
//		is_admin, login, name, err := db.GetInfoAboutMember(member)
//		if err != nil {
//			r.notificationAboutError("")
//			return
//		}
//		admin := ""
//		if is_admin {
//			admin = "Администратор"
//		}
//		message = message + fmt.Sprintf("%d. %s (@%s) %s\n", i+1, name, login, admin)
//	}
//
//	msg := tgbotapi.NewMessage(r.chatId, message)
//	_, _ = r.bot.Send(msg)
//}
