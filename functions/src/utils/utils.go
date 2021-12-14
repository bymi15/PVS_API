package utils

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/joho/godotenv"
)

type IdentityResponse struct {
	Identity *Identity `json:"identity"`
	User     *User     `json:"user"`
	SiteUrl  string    `json:"site_url"`
	Alg      string    `json:"alg"`
}

type Identity struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type User struct {
	AppMetaData  *AppMetaData  `json:"app_metadata"`
	Email        string        `json:"email"`
	Exp          int           `json:"exp"`
	Sub          string        `json:"sub"`
	UserMetadata *UserMetadata `json:"user_metadata"`
}
type AppMetaData struct {
	Provider string `json:"provider"`
}
type UserMetadata struct {
	FullName string `json:"full_name"`
}

type Response struct {
	Msg              string `json:"msg"`
	IdentityResponse string `json:"identity_response"`
}

func SetDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
}

func CreateApiResponse(v interface{}) []byte {
	var response []byte
	if v != "" {
		jsonBody, err := json.Marshal(v)
		if err != nil {
			log.Fatalf("An error occurred in JSON marshal. Err: %s", err)
		}
		response = jsonBody
	}

	return response
}

func ParseRequestBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func ServeFunction(url string, handler func(http.ResponseWriter, *http.Request)) {
	port := flag.Int("port", -1, "specify a port")
	flag.Parse()
	listener := gateway.ListenAndServe
	addr := ""
	if *port != -1 {
		err := godotenv.Load()
		if err != nil {
			log.Print("Failed to load .env file")
		}
		addr = fmt.Sprintf(":%d", *port)
		listener = http.ListenAndServe
		http.Handle("/", http.FileServer(http.Dir("./public")))
	}
	http.HandleFunc(url, AuthMiddleware(handler))

	log.Printf("Function `%s` running on %s...", url, addr)
	log.Fatal(listener(addr, nil))
}

func AuthMiddleware(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lc, ok := lambdacontext.FromContext(r.Context())
		if !ok {
			log.Fatalf("error retrieving context %+v", r.Context())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		identityResponse := lc.ClientContext.Custom["netlify"]
		raw, _ := base64.StdEncoding.DecodeString(identityResponse)
		data := IdentityResponse{}
		_ = json.Unmarshal(raw, &data)
		if data.User == nil {
			log.Fatalf("forbidden access for request bearer %+v", identityResponse)
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		log.Printf("User '%s' has been successfully authenticated.", data.User.UserMetadata.FullName)
		h(w, r)
	}
}
