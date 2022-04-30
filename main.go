package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	twitter "github.com/g8rswimmer/go-twitter/v2"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var client *twitter.Client

type authorize struct {
	Token string
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func isValidHandleName(handle string) bool {
	var isValid bool = true

	// an valid handle name only contains numbers, letters and undersocres.
	// also, must be more than 3 characters long and less or equal to 15.

	match, _ := regexp.MatchString("^[A-Za-z0-9_]*$", handle)
	if len(handle) < 3 || len(handle) > 15 || !match {
		isValid = false
	}

	return isValid
}

func translateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	// first, check if username is a valid handle name
	if !isValidHandleName(username) {
		fmt.Println("O username é inválido.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	opts := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{
			twitter.UserFieldID,
		},
	}

	user, err := client.UserNameLookup(r.Context(), []string{username}, opts)

	if err != nil {
		fmt.Println("Ocorreu um erro ao buscar o usuário.")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user.Raw.Users[0].ID)
}

func main() {
	token := flag.String("token", goDotEnvVariable("BEARER"), "twitter API token")

	client = &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/user/{username}", translateUserHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}