package types

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
)

// MarshalAny is a convenience function for packing the provided value in an
// Any and then proto marshaling it to bytes
func MarshalJSONAny(m codec.JSONMarshaler, x interface{}) ([]byte, error) {
	msg, ok := x.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("can't proto marshal %T", x)
	}

	any := &types.Any{}
	err := any.Pack(msg)
	if err != nil {
		return nil, err
	}

	return m.MarshalJSON(any)
}

// UnmarshalAny is a convenience function for proto unmarshaling an Any from
// bz and then unpacking it to the interface pointer passed in as iface using
// the provided AnyUnpacker or returning an error
//
// Ex:
//		var x MyInterface
//		err := UnmarshalAny(unpacker, &x, bz)
func UnmarshalJSONAny(m codec.Marshaler, iface interface{}, bz []byte) error {
	any := &types.Any{}

	err := m.UnmarshalJSON(bz, any)
	if err != nil {
		return err
	}

	return m.UnpackAny(any, iface)
}
