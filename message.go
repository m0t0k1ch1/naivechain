package main

import "encoding/json"

type messageType int

const (
	messageTypeQueryLatest    messageType = iota
	messageTypeQueryAll       messageType = iota
	messageTypeResponseBlocks messageType = iota
)

func (ms messageType) name() string {
	switch ms {
	case messageTypeQueryLatest:
		return "QUERY_LATEST"
	case messageTypeQueryAll:
		return "QUERY_ALL"
	case messageTypeResponseBlocks:
		return "RESPONSE_BLOCKS"
	default:
		return "UNKNOWN"
	}
}

type Message struct {
	Type messageType `json:"type"`
	Data string      `json:"data"`
}

func newBlocksMessage(blocks []*Block) (*Message, error) {
	b, err := json.Marshal(blocks)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type: messageTypeResponseBlocks,
		Data: string(b),
	}, nil
}

func newQueryLatestMessage() *Message {
	return &Message{
		Type: messageTypeQueryLatest,
	}
}

func newQueryAllMessage() *Message {
	return &Message{
		Type: messageTypeQueryAll,
	}
}
