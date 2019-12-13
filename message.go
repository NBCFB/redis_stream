package internal

/*
Message has header and body, where header contains ID and stream channel name. Body is defined as an interface.
 */
type Message struct {
	ID			string  		`json:"id"`			// message ID
	Stream		string			`json:"stream"` 	// stream to subscribe
	Body		string			`json:"body"`		// message body in JSON-like text
}