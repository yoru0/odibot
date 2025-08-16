package bot

import "strings"

func encodeCustomID(kind, lobbyID, actingID string) string {
	return kind + ":" + lobbyID + ":" + actingID
}

func decodeCustomID(id string) (kind, lobbyID, actingID string, ok bool) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}
