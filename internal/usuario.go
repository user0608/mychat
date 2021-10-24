package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/user0608/mychat/models"
	"github.com/user0608/mychat/utils"
)

type UserChat struct {
	chanel *Chanel
	models.User
	Connection *websocket.Conn
}

func NewUserChat(chanel *Chanel, UserName string, Connection *websocket.Conn) *UserChat {
	return &UserChat{chanel, models.User{UserName: UserName}, Connection}
}

func (u *UserChat) OnLine() {
	for {
		if _, message, err := u.Connection.ReadMessage(); err != nil {
			log.Println("Error on read message::", err.Error())
			break
		} else {
			sms := &models.Message{}
			fmt.Println("Data: ", string(message))
			if err := json.Unmarshal(message, sms); err != nil {
				log.Printf("No se pudo leer el mensaje: user %s\n, err: %s", u.UserName, err.Error())
			} else {
				log.Println(sms)
				u.chanel.messageChanel <- sms
			}
		}
	}
	u.chanel.leaveChanel <- u
}
func (u *UserChat) SendMessage(message *models.Message) error {
	message.ID = utils.GetRandonInt()
	if data, err := json.Marshal(message); err != nil {
		return err
	} else {
		err = u.Connection.WriteMessage(websocket.TextMessage, data)
		log.Printf("Message send: from %s to %s", message.SenderUserName, message.TargetUserName)
		return err
	}
}
