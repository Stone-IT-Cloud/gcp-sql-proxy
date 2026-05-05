package proxy

import (
	"fmt"
	"io"
)

const (
	connectionHost     = "127.0.0.1"
	passwordLeaveEmpty = "[LEAVE EMPTY]"
)

// ConnectionInstructions formats the operator-ready connection details.
// Output format must remain stable for tests and operator automation.
func ConnectionInstructions(port int, email string) string {
	// Keep exact line content stable and minimal for clarity.
	return fmt.Sprintf(
		"Host: %s\nPort: %d\nUser: %s\nPassword: %s\n",
		connectionHost,
		port,
		email,
		passwordLeaveEmpty,
	)
}

func WriteConnectionInstructions(w io.Writer, port int, email string) {
	_, _ = io.WriteString(w, ConnectionInstructions(port, email))
}
