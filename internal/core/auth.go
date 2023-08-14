package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	tls "github.com/refraction-networking/utls"
)

const (
	AuthCookiesUrl     = "https://auth.riotgames.com/api/v1/authorization"
	AuthRequestUrl     = "https://auth.riotgames.com/api/v1/authorization"
	MultiFactorAuthUrl = "https://auth.riotgames.com/api/v1/authorization"
	CookieReAuthUrl    = "https://auth.riotgames.com/authorize?redirect_uri=https%3A%2F%2Fplayvalorant.com%2Fopt_in&client_id=play-valorant-web-prod&response_type=token%20id_token&nonce=1"
	EntitlementUrl     = "https://entitlements.auth.riotgames.com/api/token/v1"
	UserInfoUrl        = "https://auth.riotgames.com/userinfo"
)

type UserResponse struct {
	UserId string `json:"sub"`
}

type EntitlementsResponse struct {
	EntitlementsToken string `json:"entitlements_token"`
}

type Client struct {
	httpClient *http.Client
	AuthData   *AuthSaveData
}

type AuthSaveData struct {
	AuthTokens       UriTokens `json:"authTokens"`
	EntitlementToken string    `json:"entitlementToken"`
	UserId           string    `json:"userId"`
	SavedAt          time.Time `json:"savedAt"`
}

type UriTokens struct {
	AccessToken string `json:"accessToken"`
	IdToken     string `json:"idToken"`
	ExpiresIn   int    `json:"expiresIn"`
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

	return &Client{httpClient: &http.Client{Transport: transport, Jar: cookieJar},
		AuthData: &AuthSaveData{AuthTokens: UriTokens{}, EntitlementToken: "", UserId: "", SavedAt: time.Now()}}
}

func (c *Client) Authorize(username, password string) error {
	err := c.getPreAuth()
	if err != nil {
		return err
	}

	bodyMap := map[string]any{"username": username, "password": password, "type": "auth"}
	body, err := json.Marshal(bodyMap)

	if err != nil {
		return err
	}

	req, err := createNewRequest("PUT", AuthRequestUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	loginBody := new(LoginResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&loginBody); err != nil {
		return err
	}

	if loginBody.Type == "response" {
		tokens, err := parseUriTokens(loginBody.Response.Parameters.Uri)
		if err != nil {
			return err
		}

		c.AuthData.AuthTokens = *tokens
		c.AuthData.SavedAt = time.Now()
		c.SetUserId()

		return nil
	} else if loginBody.Type == "auth" {
		if _, ok := ResponseErrors[loginBody.Error]; ok {
			return ResponseErrors[loginBody.Error]
		}
		return ErrorRiotUnknownErrorType
	} else if loginBody.Type == "multifactor" {
		return ErrorRiotMultifactor
	}

	return ErrorRiotUnknownResponseType
}

func (c *Client) MultiFactorAuth(code string) error {
	bodyMap := map[string]any{"type": "multifactor", "code": code, "rememberDevice": true}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	req, err := createNewRequest("PUT", MultiFactorAuthUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	loginBody := new(LoginResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&loginBody); err != nil {
		return err
	}

	fmt.Println("loginbody: ", loginBody)

	if loginBody.Type == "response" {
		tokens, err := parseUriTokens(loginBody.Response.Parameters.Uri)
		if err != nil {
			return err
		}
		c.AuthData.AuthTokens = *tokens
		c.AuthData.SavedAt = time.Now()
		c.SetUserId()

		return nil
	} else if loginBody.Type == "auth" {
		if _, ok := ResponseErrors[loginBody.Error]; ok {
			return ResponseErrors[loginBody.Error]
		}
		return ErrorRiotUnknownErrorType
	} else if loginBody.Type == "multifactor" {
		return ErrorRiotMultifactor
	} else {
		return ErrorRiotUnknownResponseType
	}
}

func (c *Client) SetUserId() error {
	req, err := createNewRequest("GET", UserInfoUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthData.AuthTokens.AccessToken))
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body := new(UserResponse)

	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return err
	}

	c.AuthData.UserId = body.UserId
	return c.SetEntitlementToken()
}

func (c *Client) SetEntitlementToken() error {
	req, err := createNewRequest("POST", EntitlementUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthData.AuthTokens.AccessToken))
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body := new(EntitlementsResponse)

	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return err
	}

	entitlementsToken := body.EntitlementsToken

	c.AuthData.EntitlementToken = entitlementsToken
	return nil
}

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

	req, err := createNewRequest("POST", AuthCookiesUrl, bytes.NewBuffer(body))
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
