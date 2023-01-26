package handle_messages

import (
	"fmt"
	"log"
	"net/http"

	"github.com/StkngEsk/handle_twitch_chat/common_types"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Greeting string `json:"greeting"`
}

var (
	wsUpgrader = websocket.Upgrader{}

	wsConn *websocket.Conn
)

func SendMessageToGameClient(payload common_types.PayloadFromMessageTwitch) {
	fmt.Print("\n[PAYLOAD]: \n")
	fmt.Print(payload)
	fmt.Printf("\n[CURRENT MESSAGE TO CLIENT]: %s\n", payload.Message)
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
