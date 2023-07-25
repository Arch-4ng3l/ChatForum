package storage

import "github.com/Arch-4ng3l/ChatForum/types"

type Storage interface {
	CreateNewUser(*types.SignUpRequest) error
	GetUserPassword(*types.LoginRequest) (string, error)

	CreateNewMessage(*types.CreateMessageRequest) error
	UpdateMessage(*types.UpdateMessageRequest) error
	GetMessages(*types.GetMessagesRequest) ([]*types.Message, error)
}
