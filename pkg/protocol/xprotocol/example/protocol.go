package example

import (
	"context"
	"errors"
	"fmt"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/types"
)

/**
 * Request command
 * 0     1     2           4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |magic| type| dir |      requestId        |     payloadLength     |     payload bytes ...       |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 *
 * Response command
 * 0     1     2     3     4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |magic| type| dir |      requestId        |   status  |      payloadLength    | payload bytes ..|
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 */

func init() {
	xprotocol.RegisterProtocol(ProtocolName, &proto{})
}

type proto struct{}

// types.Protocol
func (proto *proto) Name() types.ProtocolName {
	return ProtocolName
}

func (proto *proto) Encode(ctx context.Context, model interface{}) (types.IoBuffer, error) {
	switch frame := model.(type) {
	case *Request:
		return encodeRequest(ctx, frame)
	case *Response:
		return encodeResponse(ctx, frame)
	default:
		log.Proxy.Errorf(ctx, "[protocol][mychain] encode with unknown command : %+v", model)
		return nil, errors.New("unknown command type")
	}
}

func (proto *proto) Decode(ctx context.Context, data types.IoBuffer) (interface{}, error) {
	if data.Len() >= MinimalDecodeLen {
		magic := data.Bytes()[0]
		dir := data.Bytes()[2]

		// 1. magic assert
		if magic != Magic {
			return nil, fmt.Errorf("[protocol][mychain] decode failed, magic error = %d", magic)
		}

		// 2. decode
		switch dir {
		case DirRequest:
			return decodeRequest(ctx, data)
		case DirResponse:
			return decodeResponse(ctx, data)
		default:
			// unknown cmd type
			return nil, fmt.Errorf("[protocol][mychain] decode failed, direction error = %d", dir)
		}
	}

	return nil, nil
}

func NewCodec() types.Protocol {
	return &proto{}
}

// TODOs

// Heartbeater
func (proto *proto) Trigger(requestId uint64) xprotocol.XFrame {
	// not supported for poc demo
	return nil
}

func (proto *proto) Reply(requestId uint64) xprotocol.XRespFrame {
	// not supported for poc demo
	return nil
}

// Hijacker
func (proto *proto) Hijack(statusCode uint32) xprotocol.XRespFrame {
	// not supported for poc demo
	return nil
}

func (proto *proto) Mapping(httpStatusCode uint32) uint32 {
	// not supported for poc demo
	return 0
}
