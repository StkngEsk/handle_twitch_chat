package websocket_twitch_connection

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"regexp"
	"strings"
	"time"

	"github.com/StkngEsk/handle_twitch_chat/common_types"
	"github.com/StkngEsk/handle_twitch_chat/handle_messages"
)

type TwitchProps struct {
	Channel    string
	conn       net.Conn // add this field
	MsgRate    time.Duration
	Name       string
	Port       string
	OAuthToken string
	Server     string
	startTime  time.Time // add this field
}

type TwitchBot interface {
	Connect()
	Disconnect()
	HandleChat() error
	JoinChannel()
	Say(msg string) error
	Start()
}

const PSTFormat = "Jan 2 15:04:05 PST"

// Regex for parsing PRIVMSG strings.
//
// First matched group is the user's name and the second matched group is the content of the
// user's message.
var msgRegex *regexp.Regexp = regexp.MustCompile(`^user-type=\s*\w*\s* :(\w+)!\w+@\w+\.tmi\.twitch\.tv (PRIVMSG) #\w+(?: :(.*))?$`)

// Regex for parsing user commands, from already parsed PRIVMSG strings.
//
// First matched group is the command name and the second matched group is the argument for the
// command.
var cmdRegex *regexp.Regexp = regexp.MustCompile(`^!(\w+)\s?(\w+)?`)

// for TwitchProps
func timeStamp() string {
	return TimeStamp(PSTFormat)
}

// the generic variation, for bots using the TwitchBot interface
func TimeStamp(format string) string {
	return time.Now().Format(format)
}

// Connects the bot to the Twitch IRC server. The bot will continue to try to connect until it
// succeeds or is manually shutdown.
func (bb *TwitchProps) Connect() {
	var err error
	fmt.Printf("[%s] Connecting to %s...\n", timeStamp(), bb.Server)

	// makes connection to Twitch IRC server
	bb.conn, err = net.Dial("tcp", bb.Server+":"+bb.Port)
	if nil != err {
		fmt.Printf("[%s] Cannot connect to %s, retrying.\n", timeStamp(), bb.Server)
		bb.Connect()
		return
	}
	fmt.Printf("[%s] Connected to %s!\n", timeStamp(), bb.Server)
	bb.startTime = time.Now()
}

// Makes the bot join its pre-specified channel.
func (bb *TwitchProps) JoinChannel() {
	fmt.Printf("[%s] Joining #%s...\n", timeStamp(), bb.Channel)
	bb.conn.Write([]byte("CAP REQ :twitch.tv/membership twitch.tv/tags twitch.tv/commands" + "\r\n"))
	bb.conn.Write([]byte("PASS " + bb.OAuthToken + "\r\n"))
	bb.conn.Write([]byte("NICK " + bb.Name + "\r\n"))
	bb.conn.Write([]byte("JOIN #" + bb.Channel + "\r\n"))

	fmt.Printf("[%s] Joined #%s as @%s!\n", timeStamp(), bb.Channel, bb.Name)
}

// Listens for and logs messages from chat. Responds to commands from the channel owner. The bot
// continues until it gets disconnected, told to shutdown, or forcefully shutdown.
func (bb *TwitchProps) HandleChat() error {
	fmt.Printf("[%s] Watching #%s...\n", timeStamp(), bb.Channel)

	// reads from connection
	tp := textproto.NewReader(bufio.NewReader(bb.conn))

	// listens for chat messages
	for {
		line, err := tp.ReadLine()
		if nil != err {

			// officially disconnects the bot from the server
			bb.Disconnect()

			return errors.New("bb.Bot.HandleChat: Failed to read line from channel. Disconnected.")
		}

		// logs the response from the IRC server
		fmt.Printf("[%s] %s\n", timeStamp(), line)

		if "PING :tmi.twitch.tv" == line {

			// respond to PING message with a PONG message, to maintain the connection
			bb.conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
			continue
		} else {
			var combinatePrivMsg string = line
			splitLine := strings.Split(line, ";")
			if len(splitLine) >= 16 {
				if strings.Split(splitLine[5], "=")[0] != "emotes" {
					combinatePrivMsg = splitLine[15]
				} else {
					combinatePrivMsg = splitLine[16]
				}
			}

			// handle a PRIVMSG message
			matches := msgRegex.FindStringSubmatch(combinatePrivMsg)
			if nil != matches {
				userName := matches[1]
				msgType := matches[2]

				switch msgType {
				case "PRIVMSG":
					msg := matches[3]
					fmt.Printf("[%s] %s: %s\n", timeStamp(), userName, msg)
					// Send Message to websocket client
					payload := getPayloadFromMessageTwitch(splitLine, msg)
					handle_messages.SendMessageToGameClient(payload)

					// parse commands from user message
					cmdMatches := cmdRegex.FindStringSubmatch(msg)
					if nil != cmdMatches {
						cmd := cmdMatches[1]
						//arg := cmdMatches[2]

						// channel-owner specific commands
						if userName == bb.Channel {
							switch cmd {
							case "tbdown":
								fmt.Printf(
									"[%s] Shutdown command received. Shutting down now...\n",
									timeStamp(),
								)

								bb.Disconnect()
								return nil
							default:
								// do nothing
							}
						}
					}
				default:
					// do nothing
				}
			}
		}
		time.Sleep(bb.MsgRate)
	}
}

// Handle message from twitch to get payload
func getPayloadFromMessageTwitch(splitLine []string, message string) common_types.PayloadFromMessageTwitch {

	var indexChange int8 = 0

	if len(splitLine) > 16 {
		indexChange = 1
	}

	payload := common_types.PayloadFromMessageTwitch{
		UserId:        strings.Split(splitLine[indexChange+14], "=")[1],
		IsBroadcaster: strings.Contains(splitLine[1], "broadcaster/1"),
		IsVip:         strings.Contains(splitLine[1], "vip/1"),
		IsMod:         strings.Contains(splitLine[1], "moderator/1"),
		IsSubscriber:  strings.Contains(splitLine[1], "subscriber/1"),
		Message:       message,
	}

	return payload
}

// Makes the bot send a message to the chat channel.
func (bb *TwitchProps) Say(msg string) error {
	if "" == msg {
		return errors.New("TwitchProps.Say: msg was empty.")
	}
	_, err := bb.conn.Write([]byte(fmt.Sprintf("PRIVMSG #%s %s\r\n", bb.Channel, msg)))
	if nil != err {
		return err
	}
	return nil
}

// Starts a loop where the bot will attempt to connect to the Twitch IRC server, then connect to the
// pre-specified channel, and then handle the chat. It will attempt to reconnect until it is told to
// shut down, or is forcefully shutdown.
func (bb *TwitchProps) Start() {
	var err error
	// Start Websocket to client
	go handle_messages.StartWebSocketToClient()

	for {

		bb.Connect()
		bb.JoinChannel()
		err = bb.HandleChat()
		if nil != err {

			// attempts to reconnect upon unexpected chat error
			time.Sleep(1000 * time.Millisecond)
			fmt.Println(err)
			fmt.Println("Starting bot again...")
		} else {
			return
		}
	}
}

// Officially disconnects the bot from the Twitch IRC server.
func (bb *TwitchProps) Disconnect() {
	bb.conn.Close()
	upTime := time.Now().Sub(bb.startTime).Seconds()
	fmt.Printf("[%s] Closed connection from %s! | Live for: %fs\n", timeStamp(), bb.Server, upTime)
}
