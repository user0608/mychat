package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/user0608/mychat/models"
	"github.com/user0608/mychat/utils"
)

var upGradeWebSocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     check,
}

type messageChanel chan *models.Message
type userChanel chan *UserChat
type Chanel struct {
	messageChanel messageChanel
	leaveChanel   userChanel
}
type WebSocketChat struct {
	users      map[string]*UserChat
	joinChanel userChanel
	chanel     *Chanel
}

func NewWebSocketChat() *WebSocketChat {
	return &WebSocketChat{
		users:      make(map[string]*UserChat),
		joinChanel: make(userChanel),
		chanel: &Chanel{
			messageChanel: make(messageChanel),
			leaveChanel:   make(userChanel),
		},
	}
}
func check(r *http.Request) bool {
	log.Printf("%s %s%s %v", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

func (w *WebSocketChat) HandlerConnextions(rw http.ResponseWriter, r *http.Request) {
	connection, err := upGradeWebSocket.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("No se abrió la connextion")
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Connection failed.")
		return
	}
	keys := r.URL.Query()
	username := strings.TrimSpace(keys.Get("username"))
	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("user-%d", utils.GetRandonInt())
	}
	u := NewUserChat(w.chanel, username, connection)
	w.joinChanel <- u
	u.OnLine()
}

func (w *WebSocketChat) UsersManager() {
	for {
		select {
		case userChat := <-w.joinChanel:
			w.AddUser(userChat)
		case message := <-w.chanel.messageChanel:
			w.SendMessage(message)
		case user := <-w.chanel.leaveChanel:
			w.DisconnectUser(user.UserName)
		}
	}
}

func (w *WebSocketChat) AddUser(userchat *UserChat) {
	if user, ok := w.users[userchat.UserName]; ok {
		user.Connection = userchat.Connection
		log.Printf("Reconnection user: %s \n", userchat.UserName)
	} else {
		w.users[userchat.UserName] = userchat
		log.Printf("Connection user: %s \n", userchat.UserName)
	}
}
func (w *WebSocketChat) DisconnectUser(username string) {
	if user, ok := w.users[username]; ok {
		defer user.Connection.Close()
		delete(w.users, username)
		log.Printf("User: %s, left the chat.", username)
	}
}

func (w *WebSocketChat) SendMessage(message *models.Message) {
	if user, ok := w.users[message.TargetUserName]; ok {
		if err := user.SendMessage(message); err != nil {
			log.Printf("No se pudo mandar el mensaje to %s", message.TargetUserName)
		}

	}
}

func StartWebSocket(port string) {
	log.Printf("Chat listening on http://localhost:%s", port)
	ws := NewWebSocketChat()
	http.HandleFunc("/ws", ws.HandlerConnextions)
	go ws.UsersManager()
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil))

}
