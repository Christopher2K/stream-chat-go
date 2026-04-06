package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"net/http"
	"net/url"
	"time"

	"github.com/christopher2k/stream-chat-cli/auth"
	"github.com/gorilla/websocket"
)

var twitchAppClientId = "1nj2nn78541x4vzurqq2zl7avcu6ql"
var twitchAddres = "irc-ws.chat.twitch.tv"
var twitchPort = "443"

func main() {
	var credentials auth.AuthTokenResponse

	savedAuth, err := auth.GetSavedAuth()
	if err == nil {
		credentials = savedAuth
		fmt.Println("> restored authentication")
	} else {
		fmt.Println("> authentication with Twitch")
		client := http.Client{}
		data := url.Values{}
		data.Set("client_id", twitchAppClientId)
		data.Set("scope", "chat:read")
		resp, err := client.PostForm("https://id.twitch.tv/oauth2/device", data)
		if err != nil {
			panic(err)
		}

		var authDeviceResponse auth.AuthDeviceResponse
		json.NewDecoder(resp.Body).Decode(&authDeviceResponse)

		fmt.Printf("> go to %s to login\n", authDeviceResponse.VerificationURI)
		fmt.Println("> waiting for authentication...")

		checkAuth := func() (auth.AuthTokenResponse, error) {
			data := url.Values{}
			data.Set("client_id", twitchAppClientId)
			data.Set("device_code", authDeviceResponse.DeviceCode)
			data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
			data.Set("scope", "chat:read")

			var authTokenResponse auth.AuthTokenResponse
			resp, err := client.PostForm("https://id.twitch.tv/oauth2/token", data)
			if err != nil {
				panic(err)
			}
			if resp.StatusCode != 200 {
				return auth.AuthTokenResponse{}, errors.New("Unauthenticated")
			}
			json.NewDecoder(resp.Body).Decode(&authTokenResponse)
			return authTokenResponse, nil
		}

		for {
			tokenResponse, err := checkAuth()
			if err != nil {
				time.Sleep(10 * time.Second)
			} else {
				auth.SaveAuth(tokenResponse)
				credentials = tokenResponse
				fmt.Println("> authenticated.")
				break
			}
		}

	}

	fmt.Print("> channel: ")

	var channelName string
	fmt.Scan(&channelName)

	fmt.Println("> connecting to Twitch")
	url := url.URL{Scheme: "wss", Host: twitchAddres + ":" + twitchPort}

	wsConn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		panic(err)
	}

	done := make(chan bool)

	go func() {
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println(string(message))
			if strings.Contains(string(message), "376") {
				wsConn.WriteMessage(websocket.TextMessage, []byte("JOIN #"+channelName))
			}
		}
	}()

	wsConn.WriteMessage(websocket.TextMessage, []byte("PASS oauth:"+credentials.AccessToken))
	wsConn.WriteMessage(websocket.TextMessage, []byte("NICK llcoolchris_"))

	fmt.Println("> connected to Twitch")
	<-done
}
