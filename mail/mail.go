package mail

import (
	"log"
	"mime"
	"sync"

	"github.com/ruffrey/smtpd"
)

type Message struct {
	To      string
	From    string
	Subject string
	Parts   []MessagePart
}

type MessagePart struct {
	Type string
	Body []byte
}

type Reader interface {
	Messages() <-chan Message
}

type RawDogReader struct {
	ListenAddr string
	messages   chan Message
	initOnce   sync.Once
}

func (rdr *RawDogReader) Messages() <-chan Message {
	rdr.initOnce.Do(func() {
		rdr.messages = make(chan Message)
		go rdr.init()
	})
	return rdr.messages
}

func (rdr *RawDogReader) init() {
	server := smtpd.NewServer(rdr.messageHandler)
	server.MaxSize = 5 * 1024 * 1024

	server.Extend("PROXY", &proxyHandler{})
	log.Printf("Listening on %s", rdr.ListenAddr)
	err := server.ListenAndServe(rdr.ListenAddr)
	log.Fatalf("Server exited %v", err)
}

func (rdr *RawDogReader) messageHandler(msg *smtpd.Message) error {
	msgParts, err := msg.Parts()
	if err != nil {
		return err
	}
	parts := make([]MessagePart, 0, len(msgParts))
	for _, part := range msgParts {
		contentType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			log.Printf("Failed parsing media type %v", err)
			return err
		}
		parts = append(parts, MessagePart{
			Type: contentType,
			Body: part.Body,
		})
	}
	rdr.messages <- Message{
		To:      msg.To[0].String(),
		From:    msg.From.String(),
		Subject: msg.Subject,
		Parts:   parts,
	}
	return nil
}
