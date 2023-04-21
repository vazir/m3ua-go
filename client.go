// Copyright 2018-2020 go-m3ua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package m3ua

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ishidawataru/sctp"
)

// Dial establishes a M3UA connection as a client.
//
// After successfully established the connection with peer, state-changing
// signals and heartbeats are automatically handled background in another goroutine.

const (
	//MAX_CLIENTS = 128
	MAX_STREAMS = 10
)

func Dial(ctx context.Context, net string, laddr *sctp.SCTPAddr, raddr *sctp.SCTPAddr, cfg *Config) (*Conn, error) {
	var err error
	conn := &Conn{
		mu:          new(sync.Mutex),
		mode:        modeClient,
		stateChan:   make(chan State),
		established: make(chan struct{}),
		sctpInfo:    &sctp.SndRcvInfo{PPID: 0x03000000, Stream: 0, Flags: 0x07},
		cfg:         cfg,
	}

	if conn.cfg.HeartbeatInfo.Interval == 0 {
		conn.cfg.HeartbeatInfo.Enabled = false
	}

	n, ok := netMap[net]
	if !ok {
		return nil, fmt.Errorf("invalid network: %s", net)
	}

	//conn.sctpConn, err = sctp.DialSCTP(n, laddr, raddr)
	conn.sctpConn, err = sctp.DialSCTPExt(
		n, laddr, raddr, sctp.InitMsg{NumOstreams: MAX_STREAMS, MaxInstreams: MAX_STREAMS})

	if err != nil {
		return nil, err
	}

	go func() {
		conn.stateChan <- StateAspDown
	}()

	go conn.monitor(ctx)
	select {
	case _, ok := <-conn.established:
		if !ok {
			return nil, ErrFailedToEstablish
		}
		return conn, nil
	case <-time.After(10 * time.Second):
		return nil, ErrTimeout
	}
}
