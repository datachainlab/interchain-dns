package types

import (
	"github.com/golang/protobuf/proto"
)

type PacketDataI interface {
	proto.Message

	ValidateBasic() error
	GetBytes() []byte
	Type() string
}

type PacketAcknowledgementI interface {
	proto.Message

	ValidateBasic() error
	GetBytes() []byte
	Type() string
}
