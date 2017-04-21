package exporter

import (
	"io"
	"net"
	"strings"
	"time"

	"net/textproto"

	"github.com/pkg/errors"
)

// see https://github.com/keyurdg/driveshaft/blob/master/src/status-loop.cpp

// not concurrency safe, so caller should ensure is only in use by one goroutine
type driveshaft struct {
	addr string
	conn *textproto.Conn
}

func newDriveshaft(addr string) *driveshaft {
	return &driveshaft{
		addr: addr,
	}
}

func (d *driveshaft) connect() error {
	// XXX: configurable timeout?
	c, err := net.DialTimeout("tcp", d.addr, 10*time.Second)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to driveshaft")
	}
	d.conn = textproto.NewConn(c)
	return nil
}

func (d *driveshaft) close() {
	_ = d.conn.Close()
	d.conn = nil
}

func (d *driveshaft) getConnection() (*textproto.Conn, error) {
	if d.conn != nil {
		return d.conn, nil
	}

	if err := d.connect(); err != nil {
		return nil, err
	}
	return d.conn, nil
}

type threadStatus struct {
	function string
	state    string
	count    int
}

func (d *driveshaft) getThreads() ([]*threadStatus, error) {
	c, err := d.getConnection()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get driveshaft connection")
	}

	// driveshaft status closes the connection after ever status call
	defer d.close()

	id, err := c.Cmd("threads")
	if err != nil {
		return nil, errors.Wrap(err, "failed to send threads command")
	}
	c.StartResponse(id)
	defer c.EndResponse(id)

	metrics := make(map[string]map[string]int)
	for {
		data, err := c.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "failed to read threads response")
		}

		// 139832651147008\tfunctionName\t0\tWaiting for work
		parts := strings.SplitN(data, "\t", 4)
		if len(parts) != 4 {
			return nil, errors.Wrap(err, "invalid threads response")
		}

		name := parts[1]
		state := "unknown"
		status := parts[3]
		switch {
		case strings.HasPrefix(status, "Waiting"):
			state = "waiting"
		case strings.HasPrefix(status, "Starting"):
			state = "starting"
		case strings.HasPrefix(status, "job_"):
			state = "working"
		}

		function, ok := metrics[name]
		if !ok {
			metrics[name] = make(map[string]int)
		}
		count, ok := function[state]
		if !ok {
			count = 0
		}
		count++

		function[name] = count
	}

	var status []*threadStatus
	for k, v := range metrics {
		for state, count := range v {
			s := &threadStatus{
				function: k,
				state:    state,
				count:    count,
			}
			status = append(status, s)
		}
	}

	return status, nil
}
