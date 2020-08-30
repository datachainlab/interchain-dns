package types

import proto "github.com/gogo/protobuf/proto"

type PacketDataI interface {
	proto.Message
	ValidateBasic() error
	GetBytes() []byte
	GetTimeoutHeight() uint64
	GetTimeoutTimestamp() uint64
	Type() string
}

type PacketAcknowledgementI interface {
	proto.Message
	ValidateBasic() error
	GetBytes() []byte
	Type() string
}
