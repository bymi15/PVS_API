package main

import (
	"log"
	"net/http"

	"github.com/bymi15/PVS/PVS_API/db"
	"github.com/bymi15/PVS/PVS_API/db/models"
	"github.com/bymi15/PVS/PVS_API/functions/src/utils"
)

func handler(w http.ResponseWriter, r *http.Request) {
	client := db.InitMongoClient()
	id := r.URL.Query().Get("id")

	utils.SetDefaultHeaders(w)
	var response []byte

	switch r.Method {
	case "GET":
		if id != "" {
			// Get by id
			showcaseRoom, err := client.ShowcaseRoomService.GetShowcaseRoomById(id)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			response = utils.CreateApiResponse(showcaseRoom)
		} else {
			// Get all
			showcaseRooms, err := client.ShowcaseRoomService.GetShowcaseRooms()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			response = utils.CreateApiResponse(showcaseRooms)
		}
	case "POST":
		showcaseRoom := models.NewShowcaseRoom()
		err := utils.ParseRequestBody(r, &showcaseRoom)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = client.ShowcaseRoomService.CreateShowcaseRoom(showcaseRoom)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		response = utils.CreateApiResponse(showcaseRoom)
	case "PUT":
		var showcaseRoom models.ShowcaseRoom
		err := utils.ParseRequestBody(r, &showcaseRoom)
		if err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = client.ShowcaseRoomService.UpdateShowcaseRoom(id, showcaseRoom)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response = utils.CreateApiResponse(showcaseRoom)
	case "DELETE":
		err := client.ShowcaseRoomService.DeleteShowcaseRoom(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response = utils.CreateApiResponse("")
	}

	w.Write(response)
}

func main() {
	utils.ServeFunction("/api/showcaseRooms", handler)
}
