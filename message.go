package main

import "encoding/json"

type MessageType int

const (
	MessageTypeQueryLatest MessageType = iota
	MessageTypeQueryAll    MessageType = iota
	MessageTypeBlocks      MessageType = iota
)

func (ms MessageType) name() string {
	switch ms {
	case MessageTypeQueryLatest:
		return "QUERY_LATEST"
	case MessageTypeQueryAll:
		return "QUERY_ALL"
	case MessageTypeBlocks:
		return "BLOCKS"
	default:
		return "UNKNOWN"
	}
}

type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

func newQueryLatestMessage() *Message {
	return &Message{
		Type: MessageTypeQueryLatest,
	}
}

func newQueryAllMessage() *Message {
	return &Message{
		Type: MessageTypeQueryAll,
	}
}

func newBlocksMessage(blocks Blocks) *Message {
	b, _ := json.Marshal(blocks)

	return &Message{
		Type: MessageTypeBlocks,
		Data: string(b),
	}
}
