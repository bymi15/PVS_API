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
	"github.com/bymi15/PVS_API/db/session"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	http.HandleFunc(url, handler)

	log.Printf("Function `%s` running on %s...", url, addr)
	log.Fatal(listener(addr, nil))
}

func GetAuthUser(w http.ResponseWriter, r *http.Request) *User {
	lc, ok := lambdacontext.FromContext(r.Context())
	if !ok {
		log.Printf("error retrieving context %+v", r.Context())
		return nil
	}
	identityResponse := lc.ClientContext.Custom["netlify"]
	raw, _ := base64.StdEncoding.DecodeString(identityResponse)
	data := IdentityResponse{}
	_ = json.Unmarshal(raw, &data)
	if data.User == nil {
		log.Printf("forbidden access for request bearer %+v", identityResponse)
	}
	return data.User
}

func CrudHandler(
	getHandler func(*mongo.Database, *User, http.ResponseWriter, *http.Request),
	createHandler func(*mongo.Database, *User, http.ResponseWriter, *http.Request),
	updateHandler func(*mongo.Database, *User, http.ResponseWriter, *http.Request),
	deleteHandler func(*mongo.Database, *User, http.ResponseWriter, *http.Request),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Crud Handler called...")
		db := session.InitDbSession()
		SetDefaultHeaders(w)
		authUser := GetAuthUser(w, r)

		switch r.Method {
		case "GET":
			getHandler(db, authUser, w, r)
		case "POST":
			createHandler(db, authUser, w, r)
		case "PUT":
			updateHandler(db, authUser, w, r)
		case "DELETE":
			deleteHandler(db, authUser, w, r)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
