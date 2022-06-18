package types

import "time"

type SessionType struct {
	SessionUUID string    `json:"session_uuid"`
	Issue       string    `json:"issue"`
	IPAddress   string    `json:"ip_address"`
	LastSeen    time.Time `json:"last_seen"`
}
