package conversion

import (
	"time"
)

// func GetTimeDifferenceInSeconds(t1, t2 string) (int, error) {
// 	var err error
// 	var t1Object, t2Object *time.Time
// 	if t1 == "" || t2 == "" {
// 		return -1, err
// 	}
// 	t1Object, err = GetTimeObjFromString(t1)
// 	t2Object, err = GetTimeObjFromString(t2)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return int(t2Object.Sub(*t1Object).Seconds()), nil
// }

func GetTimeStringInUTC(t time.Time) string {
	loc, _ := time.LoadLocation("UTC")
	locBasedTime := t.In(loc)
	return GetTimeStringInRFC3339(locBasedTime)
}

func GetTimeStringInRFC3339(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z07:00")
}

// func GetTimeStringInRFC3339WithoutOffset(t time.Time) string {
// 	return t.Format("20060102T150405Z")
// }

// func GetTimeObjFromString(timeString string) (*time.Time, error) {
// 	t, err := time.Parse("2006-01-02T15:04:05Z07:00", timeString)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &t, nil
// }

// func GetTimeInFormat(timeToFormat time.Time, format string) string {
// 	return timeToFormat.Format(format)
// }

// func GetRelativeSecondsToCurrentTime(sourceTime string) float64 {

// 	if sourceTime, _ := GetTimeObjFromString(sourceTime); sourceTime != nil {
// 		timeDifference := time.Since(*sourceTime)
// 		return timeDifference.Seconds()
// 	}
// 	return 0
// }

// func GetElapsedDays(sourceTime string) int {
// 	if sourceTime, _ := GetTimeObjFromString(sourceTime); sourceTime != nil {
// 		timeDifference := time.Since(*sourceTime)
// 		return int(timeDifference.Hours() / 24)
// 	}
// 	return 0
// }

// func AddDays(sourceTime string, days int) string {
// 	if sourceTime, _ := GetTimeObjFromString(sourceTime); sourceTime != nil {
// 		newDate := sourceTime.AddDate(0, 0, days)
// 		return GetTimeStringInUTC(newDate)
// 	}
// 	return ""
// }

// func GetTimeInDayMMYYTT(timetoFormat *time.Time) string {
// 	var timeInString string
// 	timeInString = timetoFormat.Weekday().String() + ", " + strconv.Itoa(timetoFormat.Day()) + " " + timetoFormat.Month().String() + " " + strconv.Itoa(timetoFormat.Year()) + " " + strconv.Itoa(timetoFormat.Hour()) + "." + strconv.Itoa(timetoFormat.Minute())
// 	// if timetoFormat.Hour() >= 12 {
// 	// 	timeInString = timeInString + " PM"
// 	// } else {
// 	// 	timeInString = timeInString + " AM"
// 	// }
// 	return timeInString
// }

// func ConvertRFCToDayDDMMYYYYHHMM(timeToConvert time.Time) string {
// 	return timeToConvert.Format("Mon, 02 Jan 2006 15:04")
// }

// func ConvertISOStringToRFC3339(timeInString string) (string, error) {
// 	if mytime, err := time.Parse("2006-01-02T15:04:05.000Z", timeInString); err != nil {
// 		return "", err
// 	} else {
// 		return GetTimeStringInUTC(mytime), nil
// 	}
// }
