package server

import (
	"errors"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/terinkov_HW2/models"

	// "github.com/terinkov_HW2/storage"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/terinkov_HW2/docs"
)

// var sessionRepository storage.SessionRepository

type taskStorage interface {
	GetTaskById(key string) (*models.Task, error)
	PostTaskById(value models.Task) error
	DeleteTaskById(key string) error
	UpdateTaskById(value models.Task) error
}

// Интерфейс репозитория сессий
type SessionStorage interface {
	PostSession(session models.Session) error
	GetSession(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
}

type UserStorage interface {
	GetUserByUserLogin(userLogin string) (*models.User, error)
	PostUserByNameAndPassword(user models.User) error
	DeleteUserByUserLogin(userLogin string) error
	UpdateUser(newUserInfo models.User) error
	LoginUser(login string, password string) (*models.User, error)
}

type Server struct {
	taskStorage taskStorage
	sessionsStorage SessionStorage
	userStorage UserStorage
}

func newServer(storage taskStorage, sessionsStorage SessionStorage, userStorage UserStorage) *Server {
	return &Server{taskStorage: storage, sessionsStorage: sessionsStorage, userStorage: userStorage}
}

// @Summary Create and run in progress task
// @Description Create and run in progress task from user, by "image" and "filter"
// @Param key query string true "Value"
// @Success 201 {string} string "Task created"
// @Failure 400 {string} string "Missing key"
// @Failure 405 {string} string "Failed to create task"
// @Router /task [post]
func (s *Server) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	// {"image":"IMAGE", "filter":"FILTER_NAME"}
	sessionId, err := getResultSessionTokenCheck(r)
	if err!=nil {
		http.Error(w, "You haven't login yet!", http.StatusUnauthorized)
		return 
	}

	_, err = s.sessionsStorage.GetSession(sessionId)
	if err!=nil {
		http.Error(w, "There's no user with that session!", http.StatusUnauthorized)
		return
	}

	// if userSession.getSessionExpired() {
	// 	s.sessionsStorage.DeleteSession(userSession.SessionId)
	// 	http.Error(w,"Expired session",http.StatusUnauthorized)
	// 	return
	// }

	taskUuid := uuid.New().String()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	task, err := getCreatingStructPostTaskStatus(body, taskUuid)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = s.taskStorage.PostTaskById(task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		task.DoTask() // Вызов метода DoTask на копии задачи

		// После выполнения задачи вы можете обновить статус в хранилище
		s.taskStorage.UpdateTaskById(task)
	}()

	err = getStatusResponse(w, map[string]string{"uuid_value": task.UUID}, http.StatusCreated)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

}

func getResultSessionTokenCheck(r *http.Request) (string, error) {
	for _, c := range r.Cookies() {
		fmt.Println(c)
   	}
	
	sessionToken, err := r.Cookie("session-token") 
	
	log.Println(sessionToken)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			log.Println("No cookie!")
			return "", errors.New("there is no session token")
		}
		// For any other type of error, return a bad request status
		log.Println("Error")
		return "", errors.New("an error occured on server")
	}
	return sessionToken.Value, nil
}


func getCreatingStructPostTaskStatus(body []byte, taskUuid string) (models.Task, error) {
	var task models.Task
	err := json.Unmarshal(body, &task)
	if err != nil {
		return task, err
	}
	task.UUID=taskUuid
	task.Status="in_progress"
	return task, err
}

func (s *Server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getResultSessionTokenCheck(r)
	if err!=nil {
		http.Error(w, "You haven't login yet!", http.StatusUnauthorized)
		return 
	}

	_, err = s.sessionsStorage.GetSession(sessionId)
	if err!=nil {
		http.Error(w, "There's no user with that session!", http.StatusUnauthorized)
		return
	}

	key := getStatusRequest(r)
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	value, err := s.taskStorage.GetTaskById(key)
	if err != nil || value == nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	err = getStatusResponse(w, map[string]string{"status": value.Status}, http.StatusOK)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := getResultSessionTokenCheck(r)
	if err!=nil {
		http.Error(w, "You haven't login yet!", http.StatusUnauthorized)
		return 
	}

	_, err = s.sessionsStorage.GetSession(sessionId)
	if err!=nil {
		http.Error(w, "There's no user with that session!", http.StatusUnauthorized)
		return
	}

	key := getStatusRequest(r)
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	value, err := s.taskStorage.GetTaskById(key)
	if err != nil || value == nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		log.Println("Key not found")
		return
	}

	err = getStatusResponse(w, map[string]string{"status": value.Status, "result": value.Content}, http.StatusOK)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
}

func getStatusRequest(r *http.Request) string {
	key := chi.URLParam(r, "task_id")
	return key
}
func getStatusResponse(w http.ResponseWriter, jsonContent map[string]string, serverStatus int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(serverStatus)
	err := json.NewEncoder(w).Encode(jsonContent)
	return err
}


// User
func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Registration:")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	newUser, err := getCreatingStructUserStatus(body)
	if err!=nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = s.userStorage.PostUserByNameAndPassword(newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusMethodNotAllowed)
		return
	}

	err = getStatusResponse(w, map[string]string{"user_login": newUser.Login}, http.StatusCreated)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	
}

func getCreatingStructUserStatus(body []byte) (models.User, error) {
	var user models.User
	err := json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}
	// user.UserId=userId
	log.Println(user)
	return user, err
}

func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	// var loginInfo models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return 
	}

	loginInfo, err := getLogingStructUserStatus(body)
	if err!=nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return 
	}

	loginUser, err := s.userStorage.LoginUser(loginInfo.Login, loginInfo.Password)
	if err!=nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return 
	}

	sessionToken := uuid.NewString()
	expires := time.Now().Add(120*time.Second)

	err = s.sessionsStorage.PostSession(models.Session{
		UserLogin: loginUser.Login,
		SessionId: sessionToken,
		ExpirityTime: expires,
	})

	if err!=nil{
		http.Error(w, "Can't create Token", http.StatusInternalServerError)
		return
	}
	
	cookie:=http.Cookie{
		Name: "session-token",
		Value: sessionToken,
		Expires: expires,
	}
	http.SetCookie(w, &cookie)

	json.NewEncoder(w).Encode(map[string]string{"token":sessionToken})
}

func getLogingStructUserStatus(body []byte) (models.User, error) {
	var user models.User
	err := json.Unmarshal(body, &user)
	if err != nil {
		return user, err
	}
	return user, err
}

// func (s *Server) getLoginUserInfo(body []byte) (models.User, error) {
// 	var userLoginInfo models.User
// 	err := json.Unmarshal(body, &userLoginInfo)
// 	user = s.userStorage.LoginUser(userLoginInfo.Login)
// 	if err != nil {
// 		return user, err
// 	}
// 	return user, err
// }


func CreateAndRunServer(storage taskStorage, sesStor SessionStorage, usrStor UserStorage, addr string) error {
	server := newServer(storage, sesStor, usrStor)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/task", server.postTaskHandler)
	r.Get("/status/{task_id}", server.getStatusHandler)
	r.Get("/result/{task_id}", server.getResultHandler)
	r.Post("/register", server.registerUser)
	r.Post("/login", server.loginUser)
	

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return httpServer.ListenAndServe()
}

