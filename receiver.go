package internal

/*
Receiver receives messages
 */
type Receiver struct {
	Stream   	string			// stream to receive from
	LastID   	string			// last ID we have checked
	Messages 	chan Message	// message queue for receiving
	Shutdown 	func() error	// function to shutdown the receiver
}

