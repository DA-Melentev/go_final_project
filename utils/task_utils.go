package utils

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	repeatRegexp = `^[dywm]( -?([0-9]+,)*-?[0-9]+)?( ([0-9]+,)*[0-9]+)?$`
)

var (
	weekdayMap = map[int]time.Weekday{
		1: time.Monday,
		2: time.Tuesday,
		3: time.Wednesday,
		4: time.Thursday,
		5: time.Friday,
		6: time.Saturday,
		7: time.Sunday,
	}
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	if len(repeat) == 0 {
		return "", errors.New("repeat parameter is empty")
	}

	exp, _ := regexp.Compile(repeatRegexp)
	if !exp.MatchString(repeat) {
		return "", errors.New("repeat parameter mismatch pattern")
	}

	r := repeat[0]
	switch r {
	case 'd':
		return nextDateD(now, parsedDate, repeat)
	case 'y':
		return nextDateY(now, parsedDate, repeat)
	case 'w':
		return nextDateW(now, parsedDate, repeat)
	case 'm':
		return nextDateM(now, parsedDate, repeat)
	default:
		return "", errors.New("wrong repeat modifier")
	}
}

func nextDateD(now time.Time, date time.Time, repeat string) (string, error) {
	tokens := strings.Split(repeat, ` `)
	if len(tokens) != 2 {
		return "", errors.New("mismatch format of `d` repeat modifier")
	}
	days, err := strconv.Atoi(tokens[1])
	if err != nil {
		return "", errors.New("wrong days count of `d` modifier")
	}
	if days < 1 || days > 400 {
		return "", errors.New("days count must be in range 1..400")
	}
	result := date
	for {
		result = result.AddDate(0, 0, days)
		if result.After(now) {
			break
		}
	}
	return result.Format("20060102"), nil
}

func nextDateY(now time.Time, date time.Time, repeat string) (string, error) {
	if len(repeat) != 1 {
		return "", errors.New("mismatch format of `y` repeat modifier")
	}

	result := date
	for {
		result = result.AddDate(1, 0, 0)
		if result.After(now) {
			break
		}
	}
	return result.Format("20060102"), nil
}

func nextDateW(now time.Time, date time.Time, repeat string) (string, error) {
	tokens := strings.Split(repeat, ` `)
	if len(tokens) != 2 {
		return "", errors.New("mismatch format of `w` repeat modifier")
	}
	weekdaysRaw := strings.Split(tokens[1], `,`)

	var requiredDays []time.Weekday
	for _, dayNumber := range weekdaysRaw {
		number, err := strconv.Atoi(dayNumber)
		if err != nil {
			return "", fmt.Errorf("error while handle `w` pattern `%s`: %v", repeat, err)
		}
		weekday, ok := weekdayMap[number]
		if !ok {
			return "", fmt.Errorf("mismatch format of `w` repeat modifier. wrong day of week - %d", number)
		}
		requiredDays = append(requiredDays, weekday)
	}
	result := date
	for {
		result = result.AddDate(0, 0, 1)
		if slices.Contains(requiredDays, result.Weekday()) && result.After(now) {
			break
		}
	}
	return result.Format("20060102"), nil
}

func nextDateM(now time.Time, date time.Time, repeat string) (string, error) {
	tokens := strings.Split(repeat, ` `)
	if len(tokens) < 2 || len(tokens) > 3 {
		return "", errors.New("mismatch format of `m` repeat modifier")
	}

	requiredDaysRaw := strings.Split(tokens[1], `,`)
	requiredDays, err := convertSliceToInt(requiredDaysRaw, -2, 31)
	if err != nil {
		return "", fmt.Errorf("error while handle `m` pattern `%s`: %v", repeat, err)
	}

	var requiredMonths []int
	if len(tokens) == 3 {
		requiredMonthsRaw := strings.Split(tokens[2], `,`)
		var err error
		requiredMonths, err = convertSliceToInt(requiredMonthsRaw, 1, 12)
		if err != nil {
			return "", fmt.Errorf("error while handle `m` pattern `%s`: %v", repeat, err)
		}
	}

	result := date
	for {
		result = result.AddDate(0, 0, 1)
		year, month, day := result.Date()
		validDay := false
		for _, requiredDay := range requiredDays {
			if requiredDay < 0 {
				var err error
				requiredDay, err = getPositiveDayOfMonth(year, month, requiredDay)
				if err != nil {
					return "", fmt.Errorf("error while handle `m` pattern `%s`: %v", repeat, err)
				}
			}
			if day == requiredDay {
				validDay = true
				break
			}
		}

		validMonth := false
		if len(tokens) == 3 {
			if slices.Contains(requiredMonths, int(month)) {
				validMonth = true
			}
		} else {
			validMonth = true
		}

		if validDay && validMonth && result.After(now) {
			break
		}
	}
	return result.Format("20060102"), nil
}

func convertSliceToInt(s []string, min, max int) ([]int, error) {
	var result []int
	for _, number := range s {
		day, err := strconv.Atoi(number)
		if err != nil {
			return nil, err
		}
		if day < min || day > max {
			return nil, fmt.Errorf("value %d is out of range %d..%d", day, min, max)
		}
		result = append(result, day)
	}
	return result, nil
}

func getPositiveDayOfMonth(year int, month time.Month, negativeDay int) (int, error) {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := 0
	for day := firstDay; day.Month() == firstDay.Month(); day = day.AddDate(0, 0, 1) {
		lastDay++
	}
	switch negativeDay {
	case -1:
		return lastDay, nil
	case -2:
		return lastDay - 1, nil
	default:
		return 0, errors.New("wrong negative day of month. -1, -2 only available")
	}
}
