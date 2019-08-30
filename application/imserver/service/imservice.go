package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"
	"log"
	"micro-me/application/common/baseerror"
	"net/http"
	"sync"
)

type (
	ImService struct {
		rabbitMqService *RabbitMqService
		clients         map[string][]*websocket.Conn //key： token(uuid)
		Address         string
		lock            sync.Mutex
		upgrader        *websocket.Upgrader
	}

	SendMsgRequest struct {
		FromToken string `json:"fromToken"`
		ToToken   string `json:"toToken"`
		Body      string `json:"body"`
	}
	SendMsgResponse struct {
		FromToken string `json:"fromToken"`
		Body      string `json:"body"`
	}

	LoginRequest struct {
		Token string `json:"token"`
	}

	ImServiceOptions func(*ImService)
)

const (
	DefaultAddress = ":7272"
)

var (
	UserNotLoginErr = baseerror.NewBaseError("用户未登录")
)

func NewImService(topic string, rabbitMqService *RabbitMqService, opts ImServiceOptions) (*ImService, error) {
	if err := broker.Init(); err != nil {
		return nil, err
	}
	if err := broker.Connect(); err != nil {
		return nil, err
	}
	service := &ImService{
		rabbitMqService: rabbitMqService,
		clients:         make(map[string][]*websocket.Conn, 0),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			//跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	if opts != nil {
		opts(service)
	}

	if service.Address == "" {
		service.Address = DefaultAddress
	}
	return service, nil
}

func (s *ImService) SendMsg(r *SendMsgRequest) (*SendMsgResponse, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	conns := s.clients[r.ToToken]
	if conns == nil {
		return nil, UserNotLoginErr
	}
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(r.Body)); err != nil {
			log.Printf("conn send msg error: %+v", err)
			s.clients[r.ToToken] = nil
			conn.Close()
			return nil, err
		}
	}
	fmt.Printf("======>%+v\n", r)
	return &SendMsgResponse{}, nil
}

func (s *ImService) Subscribe() {

	s.rabbitMqService.Subscribe(func(msg []byte) error {
		r := new(SendMsgRequest)

		fmt.Printf("borker subscribe  ===>%s\n", msg)
		if err := json.Unmarshal(msg, r); err != nil {
			log.Printf("im service subscribe json.Unmarshal error [%+v]", err)
			return err
		}

		if _, err := s.SendMsg(r); err != nil {
			log.Printf("im service subscribe sendMsg error [%+v]", err)
			return err
		}
		return nil
	})

}

func (s *ImService) Run() {
	log.Println("websocket has listening at ....")
	http.HandleFunc("/ws", s.loginHandler)
	if err := http.ListenAndServe(s.Address, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *ImService) loginHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("messageHandler error : %+v", err)
		return
	}

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("read message error: %+v", err)
		return
	}

	if messageType != websocket.TextMessage {
		log.Printf("read message type error")
		return
	}
	loginRequest := new(LoginRequest)

	if err := json.Unmarshal(message, loginRequest); err != nil {
		log.Printf("read message json Unmarshal error: %+v", err)
		return
	}

	s.clients[loginRequest.Token] = append(s.clients[loginRequest.Token], conn)

	fmt.Println("客户上线==>", r)
}
