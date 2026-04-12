package main

type MessageType int

const (
	Unknown MessageType = iota

	RPL_WELCOME
	RPL_YOURHOST
	RPL_CREATED
	RPL_MYINFO
	RPL_MOTDSTART
	RPL_ENDOFMOTD
	RPL_NAMEREPLY
	RPL_ENDOFNAMES

	PRIVMSG
)

var messageTypes = map[MessageType]string{
	RPL_WELCOME:    "001",
	RPL_YOURHOST:   "002",
	RPL_CREATED:    "003",
	RPL_MYINFO:     "004",
	RPL_MOTDSTART:  "375",
	RPL_ENDOFMOTD:  "376",
	RPL_NAMEREPLY:  "353",
	RPL_ENDOFNAMES: "366",
	PRIVMSG:        "PRIVMSG",
}

func (m MessageType) String() string {
	return messageTypes[m]
}

type Message struct {
	Type       MessageType
	Host       *string
	User       *string
	Message    *string
	RawContent string
}
