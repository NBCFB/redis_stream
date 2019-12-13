package internal

import (
	"github.com/go-redis/redis"
	"time"
)

/*
Broker implements a broker for message publish and subscribe
 */
type Broker struct {
	Client *redis.Client
}

/*
XPub publishes a message to a channel. The channel is specified in the header of the message. It returns the result in
string.
 */
func (s *Broker) XPub(msg *Message, maxLen int) (res string, err error) {

	// values must be a map
	m := map[string]interface{}{
		"body": msg.Body,
	}

	// prepare parameters for publish message to the channel
	args := &redis.XAddArgs{
		Stream: msg.Stream,
		ID:     "*",
		Values:	m,
		MaxLen: int64(maxLen),
	}

	// publish
	res, err = s.Client.XAdd(args).Result()
	if err != nil {
		return
	}

	return res, nil
}

/*
XSub subscribe a stream channel. It expects to receive messages whose ID > lastID. It returns a receiver pointer.
 */
func (s *Broker) XSub(stream, lastID string, blockTime int) (recv *Receiver, err error){
	if lastID == "" {
		lastID = "0"
	}

	recv = &Receiver{
		Stream:   stream,
		LastID:   lastID,
		Messages: make(chan Message),
	}

	go func() {
	XREAD:
		for {
			xStr, err := s.Client.XRead(&redis.XReadArgs{
				Block: time.Duration(blockTime) * time.Millisecond,
				Streams: []string{
					recv.Stream,
					recv.LastID,
				},
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue XREAD
				}
				break XREAD
			}

			// wrap up messages
			for _, wad := range xStr {
				for _, evt := range wad.Messages {

					recv.Messages <- Message {
						ID: evt.ID,
						Stream: wad.Stream,
						Body: evt.Values["body"].(string),
					}

					recv.LastID = evt.ID
				}
			}
		}

		// close the message channel when all the messages are collected
		close(recv.Messages)
	}()

	return
}