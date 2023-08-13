// package core

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"
// )

// func Login(username string, password string) {
// 	payload := strings.NewReader("{\"client_id\":\"play-valorant-web-prod\",\"nonce\":\"1\",\"redirect_uri\":\"https://playvalorant.com/opt_in\",\"response_type\":\"token id_token\", \"scope\": \"account openid\"}")
// 	req, _ := http.NewRequest("POST", AuthCookiesUrl, payload)

// 	req.Header.Add("Content-Type", "application/json")

// 	res, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		fmt.Printf("Error fetching cookies from '%s'. %s", AuthCookiesUrl, err)
// 	}

// 	defer res.Body.Close()
// 	cookies := res.Cookies()

// 	if res.StatusCode == 403 {
// 		fmt.Printf("%s Responded with 403, retrying... \n", AuthCookiesUrl)

// 		req, _ := http.NewRequest("POST", AuthCookiesUrl, payload)
// 		req.Header.Add("Content-Type", "application/json")

// 		for _, cookie := range cookies {
// 			req.AddCookie(cookie)
// 		}

// 		fmt.Println("req: ", req)
// 		res, err = http.DefaultClient.Do(req)
// 		if err != nil {
// 			fmt.Printf("Error fetching cookies from '%s'. %s", AuthCookiesUrl, err)
// 		}

// 		cookies = res.Cookies()
// 		fmt.Println("new cookies: ", cookies)
// 		defer res.Body.Close()

// 		fmt.Println(res)
// 	}

// }
package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	tls "github.com/refraction-networking/utls"
)

const (
	AuthCookiesUrl     = "https://auth.riotgames.com/api/v1/authorization"
	AuthRequestUrl     = "https://auth.riotgames.com/api/v1/authorization"
	MultiFactorAuthUrl = "https://auth.riotgames.com/api/v1/authorization"
	CookieReAuthUrl    = "https://auth.riotgames.com/authorize?redirect_uri=https%3A%2F%2Fplayvalorant.com%2Fopt_in&client_id=play-valorant-web-prod&response_type=token%20id_token&nonce=1"
	EntitlementUrl     = "https://entitlements.auth.riotgames.com/api/token/v1"
	PlayerInfo         = "https://auth.riotgames.com/userinfo"
)

type UserResponse struct {
	UserId string `json:"sub"`
}

type EntitlementsResponse struct {
	EntitlementsToken string `json:"entitlements_token"`
}

type Client struct {
	httpClient *http.Client
}

type UriTokens struct {
	AccessToken string
	IdToken     string
	ExpiresIn   int
}

type LoginResponseBody struct {
	Type     string `json:"type"`
	Error    string `json:"error"`
	Response struct {
		Mode       string `json:"mode"`
		Parameters struct {
			Uri string `json:"uri"`
		} `json:"parameters"`
	} `json:"response"`
	Country string `json:"country"`
}

var (
	RiotUserAgent = "RiotClient/63.0.9.4909983.4789131 rso-auth (Windows;10;;Professional, x64)"
	tlsConfig     = &tls.Config{
		MaxVersion: tls.VersionTLS13,
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
)

func New(proxy *url.URL) *Client {
	transport := &http.Transport{DialTLS: dialTls}
	cookieJar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	if proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	return &Client{httpClient: &http.Client{Transport: transport, Jar: cookieJar}}
}

func (c *Client) Authorize(username, password string) (*UriTokens, error) {
	err := c.getPreAuth()
	if err != nil {
		return nil, err
	}

	bodyMap := map[string]any{"password": password, "type": "auth", "username": username}
	body, err := json.Marshal(bodyMap)

	if err != nil {
		return nil, err
	}

	req, err := createNewRequest("PUT", "https://auth.riotgames.com/api/v1/authorization", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	loginBody := new(LoginResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&loginBody); err != nil {
		return nil, err
	}

	if loginBody.Type == "response" {
		return parseUriTokens(loginBody.Response.Parameters.Uri)
	} else if loginBody.Type == "auth" {
		if _, ok := ResponseErrors[loginBody.Error]; ok {
			return nil, ResponseErrors[loginBody.Error]
		}
		return nil, ErrorRiotUnknownErrorType
	} else if loginBody.Type == "multifactor" {
		return nil, ErrorRiotMultifactor
	} else {
		return nil, ErrorRiotUnknownResponseType
	}
}

func (c *Client) SubmitTwoFactor(code string) (*UriTokens, error) {
	bodyMap := map[string]any{"type": "multifactor", "code": code, "rememberDevice": true}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return nil, err
	}

	req, err := createNewRequest("PUT", "https://auth.riotgames.com/api/v1/authorization", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	loginBody := new(LoginResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&loginBody); err != nil {
		return nil, err
	}

	if loginBody.Type == "response" {
		return parseUriTokens(loginBody.Response.Parameters.Uri)
	} else if loginBody.Type == "auth" {
		if _, ok := ResponseErrors[loginBody.Error]; ok {
			return nil, ResponseErrors[loginBody.Error]
		}
		return nil, ErrorRiotUnknownErrorType
	} else if loginBody.Type == "multifactor" {
		return nil, ErrorRiotMultifactor
	} else {
		return nil, ErrorRiotUnknownResponseType
	}
}

func (c *Client) GetUserId() {}

func (c *Client) getPreAuth() error {
	nonce, err := GenerateNonce()
	if err != nil {
		return err
	}

	bodyMap := map[string]any{
		"acr_values": "", "claims": "",
		"client_id": "riot-client", "code_challenge": "",
		"code_challenge_method": "", "nonce": nonce,
		"redirect_uri":  "http://localhost/redirect",
		"scope":         "openid link ban lol_region account",
		"response_type": "token id_token",
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	req, err := createNewRequest("POST", "https://auth.riotgames.com/api/v1/authorization", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
