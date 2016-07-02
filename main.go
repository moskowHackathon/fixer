package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"flag"
	"github.com/nlopes/slack"
	"github.com/moskowHackathon/fixer/service"
	"strconv"
)

func main() {
	token := flag.String("slack-token", "", "Token from slack")
	flag.Parse()

	if token == nil || *token == "" {
		flag.Usage()
		os.Exit(1)
	}

	api := slack.New(*token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	//api.SetDebug(true)

	// Find the user to post as.
	authTest, err := api.AuthTest()
	if err != nil {
		fmt.Printf("Error getting channels: %s\n", err)
		return
	}

	botID := authTest.UserID

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
			// Ignore hello

			case *slack.ConnectedEvent:
				//fmt.Println("Infos:", ev.Info)
				//fmt.Println("Connection counter:", ev.ConnectionCount)
				// Replace #general with your Channel ID
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "C1NBBSKEE"))

			case *slack.MessageEvent:
				fmt.Printf("Message: %+v\n", ev)
				if strings.Contains(ev.Text, "<@"+botID+">") == false {
					break

				}


				params := slack.NewPostMessageParameters()

				//attachment := slack.Attachment{
				//	Pretext: "some pretext",
				//	Text:    "some text",
				//	// Uncomment the following part to send a field too
				//	Fields: []slack.AttachmentField{
				//		slack.AttachmentField{
				//			Title: "a",
				//			Value: "no",
				//		},
				//	},
				//}

				attachment := service.GenerateMessageForSlack("qwe")

				params.Attachments = []slack.Attachment{attachment}

				responseChannel, responseTime, err := rtm.PostMessage("C1NBBSKEE", "Текст !!!", params)

				fmt.Println( "CHANNNNNELLL >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>" )
				fmt.Println(responseChannel)
				fmt.Println(responseTime)
				fmt.Println(err)
				fmt.Println( "CHANNNNNELLL <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<" )


//=======================================================================================================

				fmt.Println( " >>>>>>>>>>>>>>>>>>>>")
				request := service.GetRequest{
					ChanId: "1",
					UserMessage: "HI ALL!!!",
				}
				response, _ := service.SendMessage(request)
				fmt.Println( " <<<<<<<<<<<<<<<<<<<<")

				rtm.SendMessage(rtm.NewOutgoingMessage("Сам дурак. Ответ эксперта - " + strconv.Itoa(int(response.ID)), ev.Channel))

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("\033[0;31mError:\033[0m %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("\033[0;31mInvalid credentials\033[0m")
				os.Exit(2)

			default:

				// Ignore other events..
				 fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
