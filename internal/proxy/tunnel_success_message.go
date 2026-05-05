package proxy

import (
	"fmt"
	"strings"
)

func TunnelSuccessMessage(port int, instance string, user string, reconnect bool, usePrivateIP bool) string {
	mode := "Tunnel connection established."
	if reconnect {
		mode = "Tunnel connection re-established."
	}
	ipMode := "public"
	if usePrivateIP {
		ipMode = "private"
	}
	lines := []string{
		mode,
		fmt.Sprintf("Target instance: %s", instance),
		fmt.Sprintf("IP connectivity mode: %s", ipMode),
		fmt.Sprintf("Host: %s", connectionHost),
		fmt.Sprintf("Port: %d", port),
		fmt.Sprintf("User: %s", user),
		"Password: <db_password>",
		fmt.Sprintf("PostgreSQL client: %s", BuildPostgresClientTemplate(connectionHost, port, "<db_name>", "<db_user>")),
	}
	return wrapLines(lines, 120)
}

func wrapLines(lines []string, width int) string {
	if width <= 0 {
		width = 120
	}
	var out []string
	for _, line := range lines {
		if len(line) <= width {
			out = append(out, line)
			continue
		}
		for len(line) > width {
			out = append(out, line[:width])
			line = line[width:]
		}
		if line != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n") + "\n"
}
