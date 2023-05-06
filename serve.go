package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"manager/debugger"
	"manager/maria"
	"manager/mongo"
	"manager/smtp_server"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// @Summary registerPerson
// @ID registerPerson
// @Tags person
// @Param request body []maria.Person true "body json"
// @Produce plain
// @Success 200 {string} string "Created"
// @Failure 400,404 {string} string "error"
// @Router /register [post]
func registerPerson(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	debugger.CheckError("Read Body", err)
	var registration []maria.Person
	debugger.CheckError("Unmarshal", json.Unmarshal(b, &registration))
	for _, e := range registration {
		id := mongo.CreatePerson(e.Email)
		smtp_server.SendMail(id, e.Email)
		w.Write(maria.CreatePlayer(e))
	}
}

// @Summary updatePlayer
// @ID updatePlayer
// @Tags update
// @Param id request body maria.Player true "json"
// @Produce plain
// @Success 200 {string} string "success"
// @Failure 400,404 {string} string "error"
// @Router /update/player [post]
func updatePlayer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	debugger.CheckError("Read Body", err)

	var player maria.Player
	debugger.CheckError("Unmarshal", json.Unmarshal(body, &player))
	maria.UpdatePlayer(player)
}

// @Summary updateTrainer
// @ID updateTrainer
// @Tags update
// @Param id request body maria.Person true "json"
// @Produce plain
// @Success 200 {string} string "success"
// @Failure 400,404 {string} string "error"
// @Router /update/trainer [post]
func updateTrainer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	debugger.CheckError("Read Body", err)

	var trainer maria.Person
	debugger.CheckError("Unmarshal", json.Unmarshal(body, &trainer))
	maria.UpdateTrainer(trainer)
}

// @Summary signPlayer
// @ID signPlayer
// @Tags sign
// @Param id request body maria.Credentials true "json"
// @Produce plain
// @Success 200 {string} string "success"
// @Failure 400,404 {string} string "error"
// @Router /sign [post]
func signPlayer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	debugger.CheckError("Read Body", err)

	var signPlayer maria.Credentials
	debugger.CheckError("Unmarshal", json.Unmarshal(body, &signPlayer))

	maria.SignPlayer(signPlayer)
}

// @Summary signPlayer
// @ID signPlayer
// @Tags sign
// @Param hash request body string true "json"
// @Produce plain
// @Success 200 {string} string "success"
// @Failure 400,404 {string} string "error"
// @Router /check [post]
func checkHash(w http.ResponseWriter, r *http.Request) {
	hash, err := io.ReadAll(r.Body)
	debugger.CheckError("Read Body", err)
	w.Write([]byte(mongo.CheckHash(string(hash))))
}

// @Summary jwtGeneration
// @ID jwtGeneration
// @Tags sign
// @Produce plain
// @Success 200 {object} Token
// @Failure 401,404 {string} string "error"
// @Router /auth [get]
func jwtGeneration(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok {
		credentials := maria.GetCredentials(username)
		debugger.CheckError("Compare", bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)))
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(6 * time.Hour)
		claims["authorized"] = true
		claims["user"] = credentials.Username
		tokenString, err := token.SignedString([]byte(credentials.Post))
		debugger.CheckError("Signed String", err)
		data := Token{Token: tokenString, Id: credentials.Person_id}
		response, err := json.Marshal(data)
		debugger.CheckError("Marshal", err)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	} else {
		w.WriteHeader(401)
	}
}

type Token struct {
	Token string `json:"token"`
	Id    int    `json:"id"`
}

func HandleRequests() {
	r := mux.NewRouter()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/auth", jwtGeneration).Methods("GET")
	r.HandleFunc("/register", registerPerson).Methods("POST")
	r.HandleFunc("/update/player", updatePlayer).Methods("POST")
	r.HandleFunc("/update/trainer", updateTrainer).Methods("POST")
	r.HandleFunc("/sign", signPlayer).Methods("POST")
	r.HandleFunc("/check", checkHash).Methods("POST")

	debugger.CheckError("Listen and Serve", http.ListenAndServe(":4000", handlers.CORS(header, methods, origins)(r)))
}
