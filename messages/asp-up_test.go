// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package messages

import (
	"testing"

	"github.com/vazir/m3ua-go/messages/params"
)

func TestAspUp(t *testing.T) {
	cases := []testCase{
		{
			"has-all",
			NewAspUp(
				params.NewAspIdentifier(1),
				params.NewInfoString("deadbeef"),
			),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x1c,
				// AspIdentifier
				0x00, 0x11, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
				// InfoString
				0x00, 0x04, 0x00, 0x0c, 0x64, 0x65, 0x61, 0x64,
				0x62, 0x65, 0x65, 0x66,
			},
		},
		{
			"has-asp",
			NewAspUp(params.NewAspIdentifier(1), nil),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x10,
				// AspIdentifier
				0x00, 0x11, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01,
			},
		},
		{
			"has-info",
			NewAspUp(nil, params.NewInfoString("deadbeef")),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x14,
				// InfoString
				0x00, 0x04, 0x00, 0x0c, 0x64, 0x65, 0x61, 0x64,
				0x62, 0x65, 0x65, 0x66,
			},
		},
		{
			"has-none",
			NewAspUp(nil, nil),
			[]byte{
				// Header
				0x01, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x08,
			},
		},
	}

	runTests(t, cases, func(b []byte) (serializeable, error) {
		v, err := ParseAspUp(b)
		if err != nil {
			return nil, err
		}
		v.Payload = nil
		return v, nil
	})
}
