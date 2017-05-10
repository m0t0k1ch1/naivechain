package main

import "encoding/json"

type MessageType int

const (
	MessageTypeQueryLatest    MessageType = iota
	MessageTypeQueryAll       MessageType = iota
	MessageTypeResponseBlocks MessageType = iota
)

func (ms MessageType) name() string {
	switch ms {
	case MessageTypeQueryLatest:
		return "QUERY_LATEST"
	case MessageTypeQueryAll:
		return "QUERY_ALL"
	case MessageTypeResponseBlocks:
		return "RESPONSE_BLOCKS"
	default:
		return "UNKNOWN"
	}
}

type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

func newBlocksMessage(blocks Blocks) (*Message, error) {
	b, err := json.Marshal(blocks)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type: MessageTypeResponseBlocks,
		Data: string(b),
	}, nil
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
