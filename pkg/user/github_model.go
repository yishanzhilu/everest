package user

// GithubToken ..
type GithubToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
	ErrorURI     string `json:"error_uri,omitempty"`
}

// NewGithubTokenResp ..
func NewGithubTokenResp(accessToken, tokenType, scope string) *GithubToken {
	return &GithubToken{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		Scope:        scope,
		ErrorMessage: "",
		ErrorDesc:    "",
		ErrorURI:     "",
	}
}

type githubTokenRespError struct {
	ErrorMessage string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
	ErrorURI     string `json:"error_uri,omitempty"`
}

func (e *githubTokenRespError) Error() string {
	return e.ErrorDesc
}

// NewGithubTokenErr ..
func NewGithubTokenErr() *GithubToken {
	return &GithubToken{
		AccessToken:  "",
		TokenType:    "",
		Scope:        "",
		ErrorMessage: "bad_verification_code",
		ErrorDesc:    "The code passed is incorrect or expired.",
		ErrorURI:     "https://developer.github.com/apps/managing-oauth-apps/troubleshooting-oauth-app-access-token-request-errors/#bad-verification-code",
	}
}

// GithubUser ..
type GithubUser struct {
	Name string `json:"name"`
	ID   uint64 `json:"id"`
}
