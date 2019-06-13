package socket
/*
	SetReadLimit - sets the limit in characters received. Probably good idea to limit it on the FE -
`{"type":"o","text":}
	SetReadDeadline - sets deadline for reading. If I just connected to WS and set this to 10 seconds - the
connection will die in 10 seconds if I don't run SetReadDeadline again and update the deadline further
	SetWriteDeadline - same as above but for writing operations
	WriteMessage(websocket.PingMessage, nil) - sends a ping message from the server to verify that the connection
is alive (this is on write)
	SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) - similar to
above, when ping is sent from the server, the socket receives pong on the read (this is on read)
*/
