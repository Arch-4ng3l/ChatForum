package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arch-4ng3l/ChatForum/storage"
	"github.com/Arch-4ng3l/ChatForum/types"
	"github.com/Arch-4ng3l/ChatForum/util"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listeningAddr string
	store         storage.Storage
}

func NewAPIServer(addr string, store storage.Storage) *APIServer {
	return &APIServer{
		addr,
		store,
	}
}

func (s *APIServer) Run() error {

	router := mux.NewRouter()

	router.HandleFunc("/login", createHTTPHandlerFunc(s.handleLoginRequest))
	router.HandleFunc("/signup", createHTTPHandlerFunc(s.handleSignUpRequest))
	router.HandleFunc("/msg", createHTTPHandlerFunc(s.handleMessageRequest))

	return http.ListenAndServe(s.listeningAddr, router)

}
func (s *APIServer) handleMessageRequest(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetMessage(w, r)
	case "POST":
		return s.handleCreateMessage(w, r)
	case "PATCH":
		return s.handleUpdateMessage(w, r)
	default:
		return fmt.Errorf("Method %s Not Allowed", r.Method)
	}

}

func (s *APIServer) handleGetMessage(w http.ResponseWriter, r *http.Request) error {

	req := &types.GetMessagesRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if !util.AuthJWT(req.Sender, req.Token) {
		return WriteJSON(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	msgs, err := s.store.GetMessages(req)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, msgs)
}

func (s *APIServer) handleCreateMessage(w http.ResponseWriter, r *http.Request) error {

	req := &types.CreateMessageRequest{}

	if err := util.DecodeJson(r, req); err != nil {

		return err
	}

	if !util.AuthJWT(req.Name, req.Token) {
		return fmt.Errorf("Access Denied")
	}

	if err := s.store.CreateNewMessage(req); err != nil {

		return err
	}

	return WriteJSON(w, http.StatusOK, "")

}

func (s *APIServer) handleUpdateMessage(w http.ResponseWriter, r *http.Request) error {

	req := &types.UpdateMessageRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if !util.AuthJWT(req.Name, req.Token) {
		return WriteJSON(w, http.StatusForbidden, map[string]string{"error": "Access denied"})
	}
	s.store.UpdateMessage(req)

	return WriteJSON(w, http.StatusOK, nil)
}

func (s *APIServer) handleSignUpRequest(w http.ResponseWriter, r *http.Request) error {

	req := &types.SignUpRequest{}
	if err := util.DecodeJson(r, req); err != nil {

		return err
	}

	if err := s.store.CreateNewUser(req); err != nil {
		return err
	}

	token, err := util.CreateJWT(req.Name)
	if err != nil {

		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *APIServer) handleLoginRequest(w http.ResponseWriter, r *http.Request) error {

	req := &types.LoginRequest{}

	if err := util.DecodeJson(r, req); err != nil {

		return err
	}

	password, err := s.store.GetUserPassword(req)

	if err != nil {
		return err
	}

	fmt.Println(password)
	if password != util.CreateHash(req.Password) {

		return fmt.Errorf("Couldn't log in")
	}

	token, err := util.CreateJWT(req.Name)
	if err != nil {

		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func createHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
