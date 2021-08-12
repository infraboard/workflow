package http

import (
	"bufio"
	"context"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

func NewProxy() *Proxy {
	return &Proxy{
		maxRespBodyBufferBytes: 64 * 1024,
		log:                    zap.L().Named("Websocket Proxy"),
		pingInterval:           20 * time.Second,
		pingWait:               3 * time.Second,
		pongWait:               60 * time.Second,
	}
}

type Proxy struct {
	maxRespBodyBufferBytes int
	log                    logger.Logger
	pingInterval           time.Duration
	pingWait               time.Duration
	pongWait               time.Duration
}

func (p *Proxy) Proxy(ctx context.Context, conn *websocket.Conn, stream io.ReadWriteCloser) {
	go p.handleRead(ctx, conn, stream)
	go p.handleWrite(ctx, conn, stream)
	p.handlePing(ctx, conn)
}

// read loop -- take messages from websocket and write to http request
func (p *Proxy) handleWrite(ctx context.Context, conn *websocket.Conn, writer io.Writer) {
	if p.pingInterval > 0 && p.pingWait > 0 && p.pongWait > 0 {
		conn.SetReadDeadline(time.Now().Add(p.pongWait))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(p.pongWait)); return nil })
	}

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("read loop done")
			return
		default:
		}
		p.log.Debug("[read] reading from socket.")
		_, payload, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err) {
				p.log.Debugf("[read] websocket closed:", err)
				return
			}
			p.log.Warnf("error reading websocket message:", err)
			return
		}
		p.log.Debugf("[read] read payload:", string(payload))
		p.log.Debug("[read] writing to requestBody:")
		n, err := writer.Write(payload)
		writer.Write([]byte("\n"))
		p.log.Debugf("[read] wrote to requestBody", n)
		if err != nil {
			p.log.Warnf("[read] error writing message to http server:", err)
			return
		}
	}
}

func (p *Proxy) handleRead(ctx context.Context, conn *websocket.Conn, reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	var scannerBuf []byte
	if p.maxRespBodyBufferBytes > 0 {
		scannerBuf = make([]byte, 0, 64*1024)
		scanner.Buffer(scannerBuf, p.maxRespBodyBufferBytes)
	}

	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			p.log.Warnf("[write] empty scan", scanner.Err())
			continue
		}
		p.log.Debugf("[write] scanned", scanner.Text())
		if err := conn.WriteMessage(websocket.TextMessage, scanner.Bytes()); err != nil {
			p.log.Errorf("[write] error writing websocket message:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		p.log.Errorf("scanner err:", err)
	}
}

func (p *Proxy) handlePing(ctx context.Context, conn *websocket.Conn) {
	if !(p.pingInterval > 0 && p.pingWait > 0 && p.pongWait > 0) {
		return
	}

	ticker := time.NewTicker(p.pingInterval)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			p.log.Debug("ping loop done")
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(p.pingWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
