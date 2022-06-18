package types

type InitializedType struct {
	IPAddress string `json:"ip_address"`
	Secret    string `json:"secret"`
}

type CheckCredentialType struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

type AuthorizedResponseType struct {
	Status      string `json:"status"`
	Issue       string `json:"issue"`
	RedirectURL string `json:"redirect_url"`
}
