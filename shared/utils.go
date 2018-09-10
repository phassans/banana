package shared

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// NewNullString returns a null able sql string
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// GetTimeIn12HourFormat returns string time in 12 hour format
func GetTimeIn12HourFormat(lTime string) (string, error) {
	if lTime == "" {
		return "", nil
	}

	// determine startTime in format
	st, err := time.Parse(TimeLayout24Hour, strings.TrimSuffix(strings.Split(lTime, "T")[1], "Z"))
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	var rtime string
	var timePrefix string
	temp := st.Format(TimeLayout12Hour)
	if strings.Contains(temp, "pm") {
		timePrefix = "pm"
	} else {
		timePrefix = "am"
	}

	if strings.Contains(temp, ":00") {
		timeParts := strings.Split(temp, ":")
		rtime = fmt.Sprintf("%s%s", timeParts[0], timePrefix)
	} else {
		rtime = temp
	}

	return rtime, nil
}

func ConvertDBDate(dbDate string) (string, error) {
	listingDate, err := time.Parse(DateFormatSQL, strings.Split(dbDate, "T")[0])
	if err != nil {
		return "", err
	}

	return listingDate.Format(DateFormat), nil
}

func ConvertDBTime(dbTime string) (string, error) {
	listingTime, err := time.Parse(TimeLayout24Hour, strings.TrimRight(strings.Split(dbTime, "T")[1], "Z"))
	if err != nil {
		return "", err
	}

	return listingTime.Format(TimeLayout24Hour), nil
}
