package handle_messages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/StkngEsk/handle_twitch_chat/handle_messages/db"
	"github.com/StkngEsk/handle_twitch_chat/handle_messages/models"
	"github.com/StkngEsk/handle_twitch_chat/types"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type Message struct {
	Greeting string `json:"greeting"`
}

var (
	wsUpgrader = websocket.Upgrader{}

	wsConn *websocket.Conn

	oauth2Config *clientcredentials.Config
)

const (
	cpGlobe string = "#61BF00"

	csGlobe string = "#FFFFFF"

	opacityGlobe float64 = 0.2

	opacityMulti float64 = 0.8

	pricePrimaryColorGlobe int8 = 99

	priceMultiColorGlobe int16 = 199

	paletteColorUrl string = "https://coolors.co/palettes/trending"
)

func SendMessageToGameClient(payload types.PayloadFromMessageTwitch, wsTwitch net.Conn) {

	command := strings.Split(payload.Message, " ")[0]

	switch command {
	case "!drop":
		user := db.GetUsers(payload.UserId)
		messageToDrop := handleMessageToDrop(user, payload)
		messageToClient, _ := json.Marshal(messageToDrop)
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(messageToClient))
		if err != nil {
			fmt.Print(err)
		}
	case "!color":
		message := "Holi @" + payload.UserName + ", escoge el color de tu globo en esta página " + paletteColorUrl + " y copia el codigo de tu color, luego enviarlo asi: !gcolor e9c46a ó !gmulticolor 80ed99 e9c46a"
		fmt.Print("\n[MESSAGE]: \n")
		fmt.Print(message)
		say(message, wsTwitch)

	}

	/*err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Print(err)
	}*/

	//EXAMPLE SAVE POINTS
	/*var user = models.User{
		Name:                "test",
		Score:               math.Floor(305.78545454*100) / 100,
		ImageUrl:            "image",
		IDUserTwitch:        "IDUserTwitch",
		ColorPrimaryGlobe:   "ColorPrimaryGlobe",
		ColorSecondaryGlobe: "ColorSecondaryGlobe",
		OpacityGlobe:        math.Floor(0.8800000*10) / 10,
	}

	_, status, err := db.SavePoints(user)
	if err != nil {
		log.Fatal("An error occurred while trying to save points " + err.Error())
		return
	}

	if status == false {
		log.Fatal("Record could not be saved.", 400)
		return
	}*/

	//EXAMPLE SAVE ALL POINTS

	/*var users = []models.User{
		{
			Name:                "test",
			Score:               math.Floor(305.78545454*100) / 100,
			ImageUrl:            "image",
			IDUserTwitch:        "5125347094",
			ColorPrimaryGlobe:   "ColorPrimaryGlobe",
			ColorSecondaryGlobe: "ColorSecondaryGlobe",
			OpacityGlobe:        math.Floor(0.8800000*10) / 10,
		},
		{
			Name:                "rickEsk91",
			Score:               math.Floor(305.78545454*100) / 100,
			ImageUrl:            "image",
			IDUserTwitch:        "512534709",
			ColorPrimaryGlobe:   "ColorPrimaryGlobe",
			ColorSecondaryGlobe: "ColorSecondaryGlobe",
			OpacityGlobe:        math.Floor(0.8800000*10) / 10,
		},
	}

	err := db.SaveAllUserPoints(users)
	if err != nil {
		log.Fatal("An error occurred while trying to save points " + err.Error())
		return
	}*/

}

func handleMessageToDrop(user models.User, payload types.PayloadFromMessageTwitch) types.PayloadToClient {
	payloadToClient := types.PayloadToClient{
		IsBroadcaster: payload.IsBroadcaster,
		UserId:        payload.UserId,
		Username:      payload.UserName,
		DisplayName:   payload.UserName,
		Emotes:        payload.Emotes,
		IsMod:         payload.IsMod,
		Message:       payload.Message,
	}

	if strings.EqualFold(user.IDUserTwitch, "") {

		client := &http.Client{}

		r, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users?id="+payload.UserId, nil)

		r.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
		r.Header.Add("Authorization", "Bearer "+getTwitchAccessToken())
		r.Header.Add("Client-id", os.Getenv("TWITCH_CLIENT_ID"))

		resp, _ := client.Do(r)

		// read response
		body, _ := ioutil.ReadAll(resp.Body)

		// decode json
		var userTwitch = new(types.GetUserTwitchResponse)
		err := json.Unmarshal(body, &userTwitch)

		if err != nil {
			fmt.Print(err)
		}

		splitUrl := strings.Split(userTwitch.Data[0].ProfileImageUrl, "300x300")
		payloadToClient.UrlUserImage = strings.Join([]string{splitUrl[0], "70x70", splitUrl[1]}, "")
		payloadToClient.CPGlobe = cpGlobe
		payloadToClient.CSGlobe = csGlobe
		payloadToClient.OpacityGlobe = opacityGlobe

	} else {

		payloadToClient.UrlUserImage = user.ImageUrl
		payloadToClient.CPGlobe = user.ColorPrimaryGlobe
		payloadToClient.CSGlobe = user.ColorSecondaryGlobe
		payloadToClient.OpacityGlobe = user.OpacityGlobe
	}

	return payloadToClient

}

func getTwitchAccessToken() string {
	oauth2Config = &clientcredentials.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return token.AccessToken
}

// Makes the bot send a message to the chat channel.
func say(msg string, wsTwitch net.Conn) error {
	if "" == msg {
		return errors.New("say: msg was empty.")
	}
	_, err := wsTwitch.Write([]byte(fmt.Sprintf("PRIVMSG #%s %s\r\n", os.Getenv("CHANNEL_NAME"), msg)))
	if nil != err {
		return err
	}
	return nil
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		// check the http.Request
		// make sure it's OK to access
		return true
	}
	var err error
	wsConn, err = wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("could not upgrade: %s\n", err.Error())
		return
	}

	defer wsConn.Close()

	// event loop
	for {
		var msg Message

		err := wsConn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error reading JSON: %s\n", err.Error())
			break
		}

		fmt.Printf("Message Received: %s\n", msg.Greeting)
		SendMessage(msg.Greeting)
	}
}

func SendMessage(msg string) {
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Printf("error sending message: %s\n", err.Error())
	}
}

func StartWebSocketToClient() {
	router := mux.NewRouter()

	router.HandleFunc("/socket", WsEndpoint)

	log.Fatal(http.ListenAndServe(":9100", router))
}
