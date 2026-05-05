package proxy

import (
	"fmt"
	"io"
	"sync"
)

var reconnectMsgMu sync.Mutex

func emitTunnelReadyMessage(w io.Writer, port int, instance string, user string, reconnect bool, usePrivateIP bool) error {
	reconnectMsgMu.Lock()
	defer reconnectMsgMu.Unlock()
	_, err := fmt.Fprint(w, TunnelSuccessMessage(port, instance, user, reconnect, usePrivateIP))
	return err
}
