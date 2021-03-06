package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bymi15/PVS_API/db/models"
	"github.com/bymi15/PVS_API/db/permissions"
	"github.com/bymi15/PVS_API/db/services"
	"github.com/bymi15/PVS_API/db/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func getShowcaseRoomService(db *mongo.Database) services.ShowcaseRoomService {
	return services.NewShowcaseRoomService(db)
}

func getHandler(db *mongo.Database, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	service := getShowcaseRoomService(db)
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
		showcaseRoom, err := service.GetShowcaseRoomById(id, userId)
		if showcaseRoom == nil || (err != nil && err.Error() == mongo.ErrNoDocuments.Error()) {
			log.Printf("GetShowcaseRoomById not found, %s", err.Error())
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("GetShowcaseRoomById err: %s", err.Error())
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
			showcaseRooms, err := service.GetShowcaseRoomsByUser(userId)
			if showcaseRooms == nil || (err != nil && err.Error() == mongo.ErrNoDocuments.Error()) {
				log.Printf("GetShowcaseRoomsByUser not found, %s", err.Error())
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("GetShowcaseRoomsByUser err: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			response = utils.CreateApiResponse(showcaseRooms)
			w.Write(response)
		} else {
			log.Print("Forbidden userId empty")
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
	} else {
		showOnlyListed := true
		if showAll == "true" {
			if permissions.CheckUserHasPermission("staff", authUser) {
				showOnlyListed = false
			} else {
				log.Printf("Forbidden authuser")
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		} else {
			// Get all public rooms
			showcaseRooms, err := service.GetShowcaseRooms(showOnlyListed)
			if err != nil {
				log.Printf("GetShowcaseRooms err: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			response = utils.CreateApiResponse(showcaseRooms)
			w.Write(response)
		}
	}
}

func createHandler(db *mongo.Database, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	service := getShowcaseRoomService(db)
	if authUser == nil || !permissions.CheckUserHasPermission("member", authUser) {
		log.Print("Forbidden authuser")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
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
		log.Printf("ParseRequestBody err: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = service.CreateShowcaseRoom(showcaseRoom)
	if err != nil {
		log.Printf("CreateShowcaseRoom err: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response = utils.CreateApiResponse(showcaseRoom)
	w.Write(response)
}

func updateHandler(db *mongo.Database, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	service := getShowcaseRoomService(db)
	if authUser == nil || !permissions.CheckUserHasPermission("member", authUser) {
		log.Print("Forbidden authuser")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	var response []byte
	var showcaseRoom models.ShowcaseRoom
	err := utils.ParseRequestBody(r, &showcaseRoom)
	if err != nil {
		log.Printf("ParseRequestBody err: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId := ""
	// members can only update their own showcase room (staff, admin can update any)
	if !permissions.CheckUserHasPermission("staff", authUser) {
		userId = authUser.Id
	}
	err = service.UpdateShowcaseRoom(id, userId, showcaseRoom)
	if err != nil {
		log.Printf("UpdateShowcaseRoom err: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response = utils.CreateApiResponse(showcaseRoom)
	w.Write(response)
}

func deleteHandler(db *mongo.Database, authUser *utils.User, w http.ResponseWriter, r *http.Request) {
	service := getShowcaseRoomService(db)
	if authUser == nil || !permissions.CheckUserHasPermission("member", authUser) {
		log.Print("Forbidden authuser")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	id := r.URL.Query().Get("id")
	var response []byte

	userId := ""
	// members can only delete their own showcase room ( admin can delete any)
	if !permissions.CheckUserHasPermission("admin", authUser) {
		userId = authUser.Id
	}

	err := service.DeleteShowcaseRoom(id, userId)
	if err != nil {
		log.Printf("DeleteShowcaseRoom err: %s", err.Error())
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
