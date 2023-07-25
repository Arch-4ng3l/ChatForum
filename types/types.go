package types

import (
	"time"

	"github.com/Arch-4ng3l/ChatForum/util"
)

type GetMessagesRequest struct {
	Sender   string `json:"name"`
	Receiver string `json:"receiver"`
	Token    string `json:"token"`
}

type UpdateMessageRequest struct {
	Name       string `json:"name"`
	Token      string `json:"token"`
	NewContent string `json:"content"`
	MessageID  uint   `json:"id"`
}

type CreateMessageRequest struct {
	Name     string `json:"name"`
	Token    string `json:"token"`
	Receiver string `json:"receiver"`
	Message  string `json:"msg"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	Sender    string    `json:"name"`
	Receiver  string    `json:"receiver"`
	Message   string    `json:"msg"`
	CreatedAt time.Time `json:"createdAt"`
	ID        uint      `json:"id"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: util.CreateHash(password),
	}
}

func (req *SignUpRequest) NewUser() *User {
	return &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: util.CreateHash(req.Password),
	}
}
