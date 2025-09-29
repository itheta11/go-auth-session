package dto

type LoginPayload struct {
	Username    string
	Password    string
	RedirectUrl string
	AppCode     string
}

type LoggedInDto struct {
	IsLoggedIn bool
	Username   string
	Email      string
	SessionId  string
	TimeLeft   int
}

type LoginTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignOut struct {
	IsSignedOut bool
	Message     string
}
