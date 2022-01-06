package utils

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
	Id           string        `json:"id"`
	AppMetaData  *AppMetaData  `json:"app_metadata"`
	Email        string        `json:"email"`
	Exp          int           `json:"exp"`
	Sub          string        `json:"sub"`
	Role         string        `json:"role"`
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
