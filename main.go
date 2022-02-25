package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/slack-go/slack"
)

var token, userID string

func init() {
	t, ok := os.LookupEnv("SLACK_TOKEN")
	if !ok {
		log.Fatal("No slack token env found can't continue, please set SLACK_TOKEN and try again")
	}
	token = t
}

func main() {
	client := slack.New(token)

	uid := promptui.Prompt{
		Label: "What is your slack member ID? (i.e WWWAA22AA)",
	}

	result, err := uid.Run()
	if err != nil {
		log.Fatal(err)
	}
	userID = result

	prompt := promptui.Select{
		Label: "Do you want to include private channels [yes/no]",
		Items: []string{"yes", "no"},
	}
	_, result, err = prompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	types := []string{"public_channel"}

	if result == "yes" {
		prompt := promptui.Prompt{
			Label:     "Are you sure, you can not rejoin private channels unless invited",
			IsConfirm: true,
		}

		_, err := prompt.Run()

		if err == nil {
			types = append(types, "private_channel")
		}
	}

	allChannels := []slack.Channel{}
	channels, next, err := client.GetConversationsForUser(&slack.GetConversationsForUserParameters{
		ExcludeArchived: true,
		UserID:          userID,
		Types:           types,
	})
	if err != nil {
		log.Fatal(err)
	}
	allChannels = append(allChannels, channels...)
	if next != "" {
		for {
			log.Println(next)
			channels, next, err = client.GetConversationsForUser(&slack.GetConversationsForUserParameters{
				Cursor:          next,
				ExcludeArchived: true,
				Types:           types,
				UserID:          userID,
			})
			if err != nil {
				log.Fatal(err)
			}
			allChannels = append(allChannels, channels...)
			if next == "" {
				break
			}
		}
	}

	chanMap := make(map[string]slack.Channel)
	chanSlice := []string{}

	for _, c := range allChannels {
		chanMap[c.Name] = c
		chanSlice = append(chanSlice, c.Name)
	}
	sort.Strings(chanSlice)
	//HACK set index 0 as exit option
	chanSlice = append([]string{"select to exit"}, chanSlice...)

	log.Println(len(chanSlice) - 1)
	for {
		prompt := promptui.Select{
			Label: "Select Channels to stay in",
			Items: chanSlice,
		}

		index, result, err := prompt.Run()

		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		//HACK to remove exit selection
		if index == 0 {
			chanSlice = chanSlice[1:]
			break
		}

		//Slow but one MUST MAINTAIN ORDER!
		copy(chanSlice[index:], chanSlice[index+1:])
		chanSlice = chanSlice[:len(chanSlice)-1]

		log.Printf("Removed %q from leave list\n", result)
	}

	finalPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Will leave %d channels, are you sure?", len(chanSlice)),
		IsConfirm: true,
	}

	_, err = finalPrompt.Run()

	if err != nil {
		log.Fatal("Exiting")
	}

	for _, c := range chanSlice {
		channel, _ := chanMap[c]
		log.Printf("Attempting to leave %s\n", channel.ID)
		_, err := client.LeaveConversation(channel.ID)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}
