package ServiceTemplate_Go

import (
	"github.com/stevenmahana/OrganizationServiceTemplate/src/controllers"
	"github.com/nats-io/go-nats"
	"runtime"
	"log"
	"os"
)

func main() {

	uri := os.Getenv("NATS_URI")

	// create NATS server connection
	natsConnection, err := nats.Connect(uri)
	if err != nil {
		log.Println("NATS Connection Error - ", err)
		panic(err)
	}

	// subscribe to subject
	log.Printf("Subscribing to subject 'organization'\n")
	natsConnection.Subscribe("organization", func(msg *nats.Msg) {

		// access New Controller
		ctlr := controllers.NewController()

		// response object with embedded error messages, receive as slice of bytes
		resp := ctlr.Controller(msg.Data)

		// publish response
		natsConnection.Publish(msg.Reply, []byte(resp));

		// handle the message
		log.Printf("Received message '%s\n", string(msg.Data) + "'")
	})

	// keep the connection alive
	runtime.Goexit()
}