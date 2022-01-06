package permissions

import (
	"log"

	"github.com/bymi15/PVS_API/db/models"
	"github.com/bymi15/PVS_API/db/utils"
)

func IsRoomOwner(authUserId string, room models.ShowcaseRoom) bool {
	return room.CreatedBy != nil && room.CreatedBy.UserId != nil && *(room.CreatedBy.UserId) != authUserId
}

func IsRoomPublicOrOwner(authUserId string, room models.ShowcaseRoom) bool {
	return room.IsListed != nil && !*(room.IsListed) && IsRoomOwner(authUserId, room)
}

func ParseRoleLevel(role string) int {
	if role == "member" {
		return 1
	} else if role == "staff" {
		return 2
	} else if role == "admin" {
		return 3
	}
	return 0
}

func CheckUserHasPermission(role string, user *utils.User) bool {
	if user == nil || ParseRoleLevel(role) < ParseRoleLevel(user.Role) {
		log.Fatalf("forbidden access for request bearer %+v", user)
		return false
	}
	log.Printf("User '%s' has access.", user.UserMetadata.FullName)
	return true
}
