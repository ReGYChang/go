package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// 處理不確定type的時間類型來源資料(string or float64),將資料統一轉為Timestamp的格式
func InterfaceToTimestamp(key interface{}) (timestamp int64, err error) {
	switch val := key.(type) {
	case string:
		t, err := ParseTimeString(val) // if there is not zone info, treat it as utc time.
		if err != nil {
			return 0, err
		}
		timestamp = t.UnixMilli()

		return timestamp, nil

	case float64:
		sec := int64(val)
		nsec := int64((val - float64(val)) * 1000000000)
		timestamp = time.Unix(sec, nsec).UnixMilli()

		return timestamp, nil

	default:
		return 0, fmt.Errorf("data type of Time field must be either string or float64")
	}
}

// 判斷 Interface{} type 類型
func CheckInterfaceType(key interface{}) (str string) {
	if key != nil {
		str = reflect.TypeOf(key).String()
	}

	return str
}

// 將轉型資料(key)從來源類型(stype)轉換成目的類型(dtype)
func TypeTransformations(key interface{}, stype string, dtype string) (newKey interface{}, err error) {

	if stype == dtype {
		return key, nil
	}

	switch {
	case stype == "int":
		_key := key.(int)

		switch {
		case dtype == "string":
			str := strconv.Itoa(_key)
			return str, nil

		case dtype == "float64":
			return float64(int64(_key)), nil

		}

	case stype == "string":
		_key := key.(string)
		var dTmp decimal.Decimal

		if dtype != "timestamp" {
			dTmp, err = decimal.NewFromString(_key)
			if err != nil {
				log.Err(err).Msg("")
				return nil, fmt.Errorf("type conversion failed")
			}
		}

		switch {
		case dtype == "int":
			return DecimalToInt(dTmp), nil

		case dtype == "float64":
			return DecimalToFloat64(dTmp), nil

		case dtype == "timestamp":
			t, err := InterfaceToTimestamp(key)
			if err != nil {
				log.Err(err).Msg("")
				return nil, fmt.Errorf("type conversion failed")
			}
			return t, nil

		}

	case stype == "float64":
		_key := key.(float64)
		d := decimal.NewFromFloat(_key)

		switch {
		case dtype == "int":
			return DecimalToInt(d), nil

		case dtype == "string":
			return d.String(), nil

		case dtype == "timestamp":
			t, err := InterfaceToTimestamp(key)
			if err != nil {
				log.Err(err).Msg("")
				return nil, fmt.Errorf("type conversion failed")
			}
			return t, nil

		}

	}

	return nil, fmt.Errorf("type conversion failed")
}

// Decimal to Float64
func DecimalToFloat64(d decimal.Decimal) float64 {
	tmpF, _ := d.Float64()
	return tmpF
}

// Decimal Float64 轉 Int
func DecimalToInt(d decimal.Decimal) int {
	tmpF, _ := d.Float64()
	// 四捨五入取整數
	i := decimal.NewFromFloat(tmpF).Round(0).IntPart()

	return int(i)
}

func Flatten(payload map[string]interface{}) map[string]interface{} {
	if info, ok := payload["info"]; ok {
		return info.(map[string]interface{})
	}
	return payload
}

// 只有轉時間格式
func Documentize(payload map[string]interface{}, timekey string) (map[string]interface{}, error) {
	ts, ok := payload[timekey]
	if !ok {
		return payload, nil
	}

	switch val := ts.(type) {
	case string:
		t, err := ParseTimeString(val) // if there is not zone info, treat it as utc time.
		if err != nil {
			return nil, err
		}
		payload[timekey] = t.Unix()

	case float64:
		sec := int64(val)
		nsec := int64((val - float64(val)) * 1000000000)
		payload[timekey] = time.Unix(sec, nsec).Unix()

	case int64:
		payload[timekey] = val

	default:
		return nil, fmt.Errorf("data type of Time field must be either string or float64")
	}

	return payload, nil
}

// 解析時間字符串; useLocalTime 用來設定是否使用本機時區解析
func ParseTimeString(ts string, useLocalTime ...bool) (time.Time, error) {
	layouts := []string{
		// 未帶時區相關參數直接當預設時區(UTC +0 time)處理
		"2006-01-02T15:04:05.999999999",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02",
		"2006/01/02T15:04:05.999999999",
		"2006/01/02 15:04:05.999999999",
		"2006/01/02 15:04:05",
		"2006/01/02",

		"2006-01-02T15:04:05.999999999Z07:00", // RFC3339
		"2006-01-02T15:04:05.999999999Z0700",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999999999Z0700",
		"2006-01-02T15:04:05.999999999 -07:00",
		"2006-01-02T15:04:05.999999999 -0700",
		"2006-01-02 15:04:05.999999999 -07:00",
		"2006-01-02 15:04:05.999999999 -0700",
		"2006/01/02T15:04:05.999999999Z07:00",
		"2006/01/02T15:04:05.999999999Z0700",
		"2006/01/02 15:04:05.999999999Z07:00",
		"2006/01/02 15:04:05.999999999Z0700",
		"2006/01/02T15:04:05.999999999 -07:00",
		"2006/01/02T15:04:05.999999999 -0700",
		"2006/01/02 15:04:05.999999999 -07:00",
		"2006/01/02 15:04:05.999999999 -0700",

		"2006-01-02Z07:00",
		"2006-01-02Z0700",
		"2006-01-02 -07:00",
		"2006-01-02 -0700",
		"2006/01/02Z07:00",
		"2006/01/02Z0700",
		"2006/01/02 -07:00",
		"2006/01/02 -0700",
	}

	var t time.Time
	for idx, layout := range layouts {
		t, err := time.Parse(layout, ts)
		if err == nil {
			if idx < 6 {
				var loc *time.Location
				if len(useLocalTime) > 0 && useLocalTime[0] {
					loc = time.Now().Location()
				} else {
					loc, _ = time.LoadLocation("") // UTC +0 time
				}
				t, _ = time.ParseInLocation(layout, ts, loc)
			}
			return t, nil
		}
	}

	return t, fmt.Errorf("cannot parse time string %s", ts)
}

// 判斷收到的時間格式表示類行為何,並轉成Unix時間戳類型(int64)
//func ParseTimeToInt64(payload map[string]interface{}) (map[string]interface{}, error) {
//
//	switch val := payload["Time"].(type) {
//	case string:
//		t, err := ParseTimeString(val) // if there is not zone info, treat it as utc time.
//		if err != nil {
//			return nil, err
//		}
//		payload["Time"] = t.Unix()
//
//	case float64:
//		sec := int64(val)
//		nsec := int64((val - float64(val)) * 1000000000)
//		payload["Time"] = time.Unix(sec, nsec).Unix()
//
//	default:
//		return nil, fmt.Errorf("data type of Time field must be either string or float64")
//	}
//
//	return payload, nil
//}

// ParseString convert other type data to string
func ParseString(v interface{}) (string, error) {
	switch v.(type) {
	case int:
		return strconv.FormatInt(int64(v.(int)), 10), nil
	case float64:
		return strconv.FormatFloat(v.(float64), 'g', -1, 64), nil
	default:
		return "", fmt.Errorf("unable to convert %v to string", v)
	}
}

// ParseFloat64 convert other type data to float64
func ParseFloat64(v interface{}) (float64, error) {
	switch v.(type) {
	case int:
		return float64(int64(v.(int))), nil
	case string:
		_v, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return 0, err
		}
		return _v, nil
	default:
		return 0, fmt.Errorf("unable to convert value %v to float64", v)
	}
}

// ParseInt convert other type data to int
func ParseInt(v interface{}) (int, error) {
	switch v.(type) {
	case string:
		_v, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		return int(_v), nil
	case float64:
		return int(v.(float64)), nil
	default:
		return 0, fmt.Errorf("unable to convert value %v to int", v)
	}
}

// ParseTimestamp convert other type data to timestamp(int64)
func ParseTimestamp(v interface{}) (int64, error) {
	switch v.(type) {
	case string:
		t, err := ParseTimeString(v.(string)) // if there is not zone info, treat it as utc time.
		if err != nil {
			return 0, err
		}
		return int64(t.Unix()), nil
	default:
		return 0, fmt.Errorf("unable to convert value %v to timestamp", v)
	}
}
