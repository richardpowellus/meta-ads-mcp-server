// Package tz provides timezone helpers for MCP servers.
// All servers operate in America/Los_Angeles (Pacific) time.
package tz

import (
	"fmt"
	"time"
	_ "time/tzdata" // embed timezone database for Windows
)

// Pacific is the America/Los_Angeles timezone location.
var Pacific *time.Location

func init() {
	var err error
	Pacific, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		Pacific = time.FixedZone("PST", -8*60*60)
	}
}

// Now returns the current time in the Pacific timezone.
func Now() time.Time {
	return time.Now().In(Pacific)
}

// Suffix returns a string to append to MCP server instructions with the
// current local date/time and timezone directive.
func Suffix() string {
	now := Now()
	return fmt.Sprintf("\n\nCurrent local date/time: %s. "+
		"All dates and times should use the America/Los_Angeles (Pacific) timezone "+
		"unless explicitly stated otherwise. When creating transactions, always use "+
		"today's date in this timezone, not UTC.",
		now.Format("Monday, January 2, 2006 3:04 PM MST"))
}
