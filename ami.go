package main

import (
	"errors"
	"fmt"

	"github.com/ivahaev/amigo"
)

type amiClient struct {
	ami *amigo.Amigo
}

func amiRun() (client *amiClient, err error) {
	settings := &amigo.Settings{
		Username: cfg.AMI.UserName,
		Password: cfg.AMI.Password,
		Host:     cfg.AMI.Host,
		Port:     cfg.AMI.Port,
	}

	client.ami = amigo.New(settings)

	client.ami.On("connect", func(message string) {
		fmt.Println("Connected", message)
	})
	client.ami.On("error", func(message string) {
		fmt.Println("Connection error:", message)
		err = errors.New(message)
	})

	client.ami.Connect()

	return
}

func (c *amiClient) amiMonitor() {
	chans := map[string]chan map[string]string{}
	e := make(chan map[string]string, 100)
	c.ami.SetEventChannel(e)

	for msg := range e {
		uniqueID := msg["UniqueID"]
		if len(uniqueID) == 0 {
			uniqueID = msg["Uniqueid"]
		}

		switch msg["Event"] {
		case "Dial":
			{
				if ch, ok := chans[uniqueID]; ok {
					ch <- msg
					continue
				}

				ch := make(chan map[string]string, 3) // capacity is 3 because only 3 events will handled
				ch <- msg
				chans[uniqueID] = ch
				//go eventsProcessor(ch)
			}
		case "DialBegin": // asterisk 13
			{
				ch := make(chan map[string]string, 3)
				ch <- msg
				chans[uniqueID] = ch
				//go eventsProcessor(ch)
			}
		case "DialEnd": // asterisk 13
			if ch, ok := chans[uniqueID]; ok {
				ch <- msg
				continue
			}

			fmt.Printf("Unknown UniqueID on DialEnd event: %s\n", uniqueID)
		case "Hangup":
			if ch, ok := chans[uniqueID]; ok {
				ch <- msg
				close(ch)
				delete(chans, uniqueID)
				continue
			}

			fmt.Printf("Unknown UniqueID on Hangup event: %s\n", uniqueID)
		}
	}
	//return nil
}
