package shared

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

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

	return st.Format(TimeLayout12Hour), nil
}
