package nosleepok

import (
	"time"
)

// time.Sleep with no coordination context is not our concern.
func Pause() {
	time.Sleep(10 * time.Millisecond)
}

// Coordination without sleep is the ideal.
func Coord(done <-chan struct{}) {
	<-done
}

// Sleep in a function that only deals with strings is fine.
func Slow(msg string) string {
	time.Sleep(time.Millisecond)
	return msg
}
