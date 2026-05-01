package bible

import (
	"fmt"
	"time"

	"maxwelljensen/bibel/pkg"
)

// EasterProgression handles calculating time until Easter
type EasterProgression struct {
	easterType string // "orthodox" or "latin"
}

// NewEasterProgression creates a new Easter progression calculator
func NewEasterProgression(easterType string) *EasterProgression {
	if easterType != "latin" && easterType != "orthodox" {
		easterType = "orthodox" // Default to orthodox
	}
	return &EasterProgression{easterType: easterType}
}

// GetNextEasterDate calculates the date of the next Easter (Orthodox or Catholic)
func (ep *EasterProgression) GetNextEasterDate() (time.Time, error) {
	currentYear := time.Now().Year()

	// Try current year
	var easterDate time.Time
	var err error

	if ep.easterType == "latin" {
		easterDate, err = eastertime.CatholicByYear(currentYear)
	} else {
		easterDate, err = eastertime.OrthodoxByYear(currentYear)
	}

	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()

	// If Easter has already passed this year, try next year
	if easterDate.Before(now) {
		if ep.easterType == "latin" {
			easterDate, err = eastertime.CatholicByYear(currentYear + 1)
		} else {
			easterDate, err = eastertime.OrthodoxByYear(currentYear + 1)
		}

		if err != nil {
			return time.Time{}, err
		}
	}

	return easterDate, nil
}

// GetTimeUntilEaster calculates the duration until the next Easter
func (ep *EasterProgression) GetTimeUntilEaster() (time.Duration, error) {
	easterDate, err := ep.GetNextEasterDate()
	if err != nil {
		return 0, err
	}

	return time.Until(easterDate), nil
}

// FormatEasterProgress formats the time until Easter as a human-readable string
func (ep *EasterProgression) FormatEasterProgress() (string, error) {
	duration, err := ep.GetTimeUntilEaster()
	if err != nil {
		return "", err
	}

	// Calculate days, hours, minutes
	totalSeconds := int(duration.Seconds())
	days := totalSeconds / (24 * 60 * 60)
	remainingSeconds := totalSeconds % (24 * 60 * 60)
	hours := remainingSeconds / (60 * 60)
	remainingSeconds %= (60 * 60)
	minutes := remainingSeconds / 60

	var label string
	if ep.easterType == "latin" {
		label = "Roman Catholic Easter"
	} else {
		label = "Easter"
	}

	switch {
	case days > 0:
		return fmt.Sprintf("%s in %d days, %d hours, %d minutes", label, days, hours, minutes), nil
	case hours > 0:
		return fmt.Sprintf("%s in %d hours, %d minutes", label, hours, minutes), nil
	default:
		return fmt.Sprintf("%s in %d minutes", label, minutes), nil
	}
}

// GetEasterProgressPercentage returns progress as a percentage (0.0 to 1.0)
// Since Easter is an annual event, we calculate progress through the year
func (ep *EasterProgression) GetEasterProgressPercentage() (float64, error) {
	currentYear := time.Now().Year()

	// Get Easter date for current year
	var currentEaster time.Time
	var err error

	if ep.easterType == "latin" {
		currentEaster, err = eastertime.CatholicByYear(currentYear)
	} else {
		currentEaster, err = eastertime.OrthodoxByYear(currentYear)
	}

	if err != nil {
		return 0.0, err
	}

	// Get Easter date for previous year to calculate yearly cycle
	var previousEaster time.Time
	if ep.easterType == "latin" {
		previousEaster, err = eastertime.CatholicByYear(currentYear - 1)
	} else {
		previousEaster, err = eastertime.OrthodoxByYear(currentYear - 1)
	}

	if err != nil {
		return 0.0, err
	}

	now := time.Now()

	// If current Easter has passed, we're between current Easter and next Easter
	if currentEaster.Before(now) {
		// Get next Easter
		var nextEaster time.Time
		if ep.easterType == "latin" {
			nextEaster, err = eastertime.CatholicByYear(currentYear + 1)
		} else {
			nextEaster, err = eastertime.OrthodoxByYear(currentYear + 1)
		}

		if err != nil {
			return 0.0, err
		}

		// Progress from current Easter to next Easter
		totalDuration := nextEaster.Sub(currentEaster)
		elapsedDuration := now.Sub(currentEaster)

		return elapsedDuration.Seconds() / totalDuration.Seconds(), nil
	} else {
		// We're between previous Easter and current Easter
		totalDuration := currentEaster.Sub(previousEaster)
		elapsedDuration := now.Sub(previousEaster)

		return elapsedDuration.Seconds() / totalDuration.Seconds(), nil
	}
}

