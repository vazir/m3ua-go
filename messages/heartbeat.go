// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/vazir/m3ua-go/messages/params"
)

// Heartbeat is a Heartbeat type of M3UA message.
//
// Spec: 3.5.5, RFC4666.
type Heartbeat struct {
	*Header
	AspIdentifier *params.Param
	HeartbeatData *params.Param
}

// NewHeartbeat creates a new Heartbeat.
func NewHeartbeat(hbData *params.Param) *Heartbeat {
	h := &Heartbeat{
		Header: &Header{
			Version:  1,
			Reserved: 0,
			Class:    MsgClassASPSM,
			Type:     MsgTypeHeartbeat,
		},
		HeartbeatData: hbData,
	}
	h.SetLength()

	return h
}

// MarshalBinary returns the byte sequence generated from a Heartbeat.
func (h *Heartbeat) MarshalBinary() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize Heartbeat")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *Heartbeat) MarshalTo(b []byte) error {
	if len(b) < h.MarshalLen() {
		return ErrTooShortToMarshalBinary
	}

	h.Header.Payload = make([]byte, h.MarshalLen()-8)

	if param := h.HeartbeatData; param != nil {
		if err := param.MarshalTo(h.Header.Payload); err != nil {
			return err
		}
	}

	return h.Header.MarshalTo(b)
}

// ParseHeartbeat decodes given byte sequence as a Heartbeat.
func ParseHeartbeat(b []byte) (*Heartbeat, error) {
	h := &Heartbeat{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a M3UA common header.
func (h *Heartbeat) UnmarshalBinary(b []byte) error {
	var err error
	h.Header, err = ParseHeader(b)
	if err != nil {
		return errors.Wrap(err, "failed to decode Header")
	}

	prs, err := params.ParseMultiParams(h.Header.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode Params")
	}
	for _, pr := range prs {
		switch pr.Tag {
		case params.HeartbeatData:
			h.HeartbeatData = pr
		default:
			return ErrInvalidParameter
		}
	}
	return nil
}

// SetLength sets the length in Length field.
func (h *Heartbeat) SetLength() {
	if param := h.HeartbeatData; param != nil {
		param.SetLength()
	}

	h.Header.Length = uint32(h.MarshalLen())
}

// MarshalLen returns the serial length of Heartbeat.
func (h *Heartbeat) MarshalLen() int {
	l := 8
	if param := h.HeartbeatData; param != nil {
		l += param.MarshalLen()
	}
	return l
}

// String returns the Heartbeat values in human readable format.
func (h *Heartbeat) String() string {
	return fmt.Sprintf("{Header: %s, HeartbeatData: %s}",
		h.Header.String(),
		h.HeartbeatData.String(),
	)
}

// Version returns the version of M3UA in int.
func (h *Heartbeat) Version() uint8 {
	return h.Header.Version
}

// MessageType returns the message type in int.
func (h *Heartbeat) MessageType() uint8 {
	return MsgTypeHeartbeat
}

// MessageClass returns the message class in int.
func (h *Heartbeat) MessageClass() uint8 {
	return MsgClassASPSM
}

// MessageClassName returns the name of message class.
func (h *Heartbeat) MessageClassName() string {
	return "ASPSM"
}

// MessageTypeName returns the name of message type.
func (h *Heartbeat) MessageTypeName() string {
	return "Heartbeat"
}

// Serialize returns the byte sequence generated from a Heartbeat.
//
// DEPRECATED: use MarshalBinary instead.
func (h *Heartbeat) Serialize() ([]byte, error) {
	log.Println("DEPRECATED: MarshalBinary instead")
	return h.MarshalBinary()
}

// SerializeTo puts the byte sequence in the byte array given as b.
//
// DEPRECATED: use MarshalTo instead.
func (h *Heartbeat) SerializeTo(b []byte) error {
	log.Println("DEPRECATED: MarshalTo instead")
	return h.MarshalTo(b)
}

// DecodeHeartbeat decodes given byte sequence as a Heartbeat.
//
// DEPRECATED: use ParseHeartbeat instead.
func DecodeHeartbeat(b []byte) (*Heartbeat, error) {
	log.Println("DEPRECATED: use ParseHeartbeat instead")
	return ParseHeartbeat(b)
}

// DecodeFromBytes sets the values retrieved from byte sequence in a M3UA common header.
//
// DEPRECATED: use UnmarshalBinary instead.
func (h *Heartbeat) DecodeFromBytes(b []byte) error {
	log.Println("DEPRECATED: use UnmarshalBinary instead")
	return h.UnmarshalBinary(b)
}

// Len returns the serial length of Heartbeat.
//
// DEPRECATED: use MarshalLen instead.
func (h *Heartbeat) Len() int {
	log.Println("DEPRECATED: use MarshalLen instead")
	return h.MarshalLen()
}
