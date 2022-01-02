package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bymi15/PVS_API/db"
	"github.com/bymi15/PVS_API/db/models"
	"github.com/bymi15/PVS_API/functions/src/utils"
)

func getHandler(client db.MongoDbClient, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	userIdParam := r.URL.Query().Get("userId")
	showAll := r.URL.Query().Get("showAll")
	var response []byte

	userId := ""
	if authUser != nil {
		userId = authUser.Id
	}

	if id != "" {
		// Get room by id (public listed room or created by user requested)
		showcaseRoom, err := client.ShowcaseRoomService.GetShowcaseRoomById(id, userId)
		if err != nil || showcaseRoom == nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		response = utils.CreateApiResponse(showcaseRoom)
		w.Write(response)
		return
	} else if userIdParam != "" {
		// Get rooms by auth user
		if userId != "" {
			showcaseRooms, err := client.ShowcaseRoomService.GetShowcaseRoomsByUser(userId)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			response = utils.CreateApiResponse(showcaseRooms)
			w.Write(response)
		} else {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
	} else {
		showOnlyListed := true
		if showAll == "true" {
			if utils.CheckUserHasPermission("staff", authUser) {
				showOnlyListed = false
			} else {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			}
		} else {
			// Get all public rooms
			showcaseRooms, err := client.ShowcaseRoomService.GetShowcaseRooms(showOnlyListed)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			response = utils.CreateApiResponse(showcaseRooms)
			w.Write(response)
		}
	}
}

func createHandler(client db.MongoDbClient, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	if authUser == nil || !utils.CheckUserHasPermission("member", authUser) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}

	var response []byte
	showcaseRoom := models.NewShowcaseRoom()
	now := time.Now().Format("2006-01-02")
	showcaseRoom.CreatedBy = &models.User{
		UserId:      &authUser.Id,
		FullName:    &authUser.UserMetadata.FullName,
		Email:       &authUser.Email,
		Role:        &authUser.Role,
		DateCreated: &now,
	}
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
	w.Write(response)
}

func updateHandler(client db.MongoDbClient, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	if authUser == nil || !utils.CheckUserHasPermission("member", authUser) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}

	id := r.URL.Query().Get("id")
	var response []byte
	var showcaseRoom models.ShowcaseRoom
	err := utils.ParseRequestBody(r, &showcaseRoom)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId := ""
	// members can only update their own showcase room (staff, admin can update any)
	if !utils.CheckUserHasPermission("staff", authUser) {
		userId = authUser.Id
	}
	err = client.ShowcaseRoomService.UpdateShowcaseRoom(id, userId, showcaseRoom)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response = utils.CreateApiResponse(showcaseRoom)
	w.Write(response)
}

func deleteHandler(client db.MongoDbClient, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	if authUser == nil || !utils.CheckUserHasPermission("member", authUser) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}

	id := r.URL.Query().Get("id")
	var response []byte

	userId := ""
	// members can only delete their own showcase room ( admin can delete any)
	if !utils.CheckUserHasPermission("admin", authUser) {
		userId = authUser.Id
	}

	err := client.ShowcaseRoomService.DeleteShowcaseRoom(id, userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response = utils.CreateApiResponse("")
	w.Write(response)
}

func main() {
	utils.ServeFunction("/api/showcase-rooms", utils.CrudHandler(getHandler, createHandler, updateHandler, deleteHandler))
}
