package harvest

import "time"

// CalculateNextRun exposes calculateNextRun for testing.
func CalculateNextRun(schedule *Schedule) *time.Time {
	return calculateNextRun(schedule)
}
