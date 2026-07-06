package _type

import (
	"github.com/shatylos/ffmpeg-thumbnails/tools/apperrors"
	"strconv"
	"time"
)

var EmptyValueError = apperrors.New("empty value")
var CanNotBeConvertedError = apperrors.New("value can not be converted")

func ToInt64(value interface{}) (int64, error) {
	if value == nil {
		return 0, apperrors.Wrap(EmptyValueError, "the value \"%s\" can not be converted to int64", value)
	}
	switch value.(type) {
	case int:
		return int64(value.(int)), nil
	case int64:
		return value.(int64), nil
	case float64:
		return int64(value.(float64)), nil
	case string:
		iVal, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			err = apperrors.Wrap(err, "the value \"%s\" can not be converted to int64", value)
			return 0, err
		}
		return iVal, nil
	}
	return 0, apperrors.Wrap(CanNotBeConvertedError, "the value \"%s\" can not be converted to int64", value)
}

func ToFloat64(value interface{}) (newVal float64, err error) {
	if value == nil {
		return 0, apperrors.Wrap(EmptyValueError, "the value \"%s\" can not be converted to float64", value)
	}
	switch value.(type) {
	case int:
		newVal = float64(value.(int))
		return
	case int64:
		newVal = float64(value.(int64))
		return
	case float64:
		newVal = value.(float64)
		return
	case string:
		newVal, err = strconv.ParseFloat(value.(string), 64)
		if err != nil {
			err = apperrors.Wrap(CanNotBeConvertedError, "the value \"%s\" can not be converted to float64: %+v", value, err)
			return
		}
		return
	}
	err = apperrors.Wrap(CanNotBeConvertedError, "the value \"%s\" can not be converted to float64", value)
	return
}

func ToString(value interface{}) (string, error) {
	switch value.(type) {
	case nil:
		return "", apperrors.Wrap(EmptyValueError, "the value is empty")
	case string:
		return value.(string), nil
	case int:
		return strconv.Itoa(value.(int)), nil
	case int64:
		return strconv.FormatInt(value.(int64), 10), nil
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64), nil
	}

	return "", apperrors.Wrap(CanNotBeConvertedError, "the value can not be converted to string")
}

func ToTimeDuration(value interface{}) (time.Duration, error) {
	switch value.(type) {
	case int:
		return time.Duration(value.(int)), nil
	case int64:
		return time.Duration(value.(int64)), nil
	case float64:
		return time.Duration(value.(float64)), nil
	case string:
		return time.ParseDuration(value.(string))
	}
	return 0, apperrors.Wrap(CanNotBeConvertedError, "the value \"%s\" can not be converted to time.Duration", value)
}

func ToInt64Slice(value interface{}) ([]int64, error) {
	switch value.(type) {
	case []interface{}:
		var result []int64
		for _, item := range value.([]interface{}) {
			itemInt64, err := ToInt64(item)
			if err != nil {
				err = apperrors.Wrap(err, "error converting slice of interfaces to slice of int64")
				return nil, err
			}
			result = append(result, itemInt64)
		}
		return result, nil
	}
	return nil, apperrors.Wrap(CanNotBeConvertedError, "the value \"%s\" can not be converted to []int64", value)
}
