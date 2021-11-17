package gocbcore

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
)

func getMapValueString(dict map[string]interface{}, key string, def string) string {
	if dict != nil {
		if val, ok := dict[key]; ok {
			if valStr, ok := val.(string); ok {
				return valStr
			}
		}
	}
	return def
}

func getMapValueBool(dict map[string]interface{}, key string, def bool) bool {
	if dict != nil {
		if val, ok := dict[key]; ok {
			if valStr, ok := val.(bool); ok {
				return valStr
			}
		}
	}
	return def
}

func randomCbUID() []byte {
	out := make([]byte, 8)
	_, err := rand.Read(out)
	if err != nil {
		logWarnf("Crypto read failed: %s", err)
	}
	return out
}

func formatCbUID(data []byte) string {
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x%02x%02x",
		data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7])
}

func clientInfoString(connID, userAgent string) string {
	agentName := "gocbcore/" + goCbCoreVersionStr
	if userAgent != "" {
		agentName += " " + userAgent
	}

	clientInfo := struct {
		Agent  string `json:"a"`
		ConnID string `json:"i"`
	}{
		Agent:  agentName,
		ConnID: connID,
	}
	clientInfoBytes, err := json.Marshal(clientInfo)
	if err != nil {
		logDebugf("Failed to generate client info string: %s", err)
	}

	return string(clientInfoBytes)
}

func trimSchemePrefix(address string) string {
	idx := strings.Index(address, "://")
	if idx < 0 {
		return address
	}

	return address[idx+len("://"):]
}
