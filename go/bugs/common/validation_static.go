package common

import (
	"errors"
	"github.com/saichler/l8types/go/ifs"
	"time"
)

func SafeCast[T any](element interface{}) (*T, error) {
	if element == nil {
		return nil, errors.New("entity not found")
	}
	result, ok := element.(*T)
	if !ok {
		return nil, errors.New("unexpected response type")
	}
	return result, nil
}

func ValidateRequired(value, fieldName string) error {
	if value == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}

func ValidateRequiredInt64(value int64, fieldName string) error {
	if value == 0 {
		return errors.New(fieldName + " is required")
	}
	return nil
}

func ValidateEnum[T ~int32](value T, nameMap map[int32]string, enumName string) error {
	if value == 0 {
		return errors.New(enumName + " must be specified")
	}
	if _, ok := nameMap[int32(value)]; !ok {
		return errors.New("invalid " + enumName + " value")
	}
	return nil
}

func ValidateDateInPast(timestamp int64, fieldName string) error {
	if timestamp > time.Now().Unix() {
		return errors.New(fieldName + " must be in the past")
	}
	return nil
}

func ValidateDateNotZero(timestamp int64, fieldName string) error {
	if timestamp == 0 {
		return errors.New(fieldName + " is required")
	}
	return nil
}

func ValidateDateAfter(date1, date2 int64, field1Name, field2Name string) error {
	if date1 <= date2 {
		return errors.New(field1Name + " must be after " + field2Name)
	}
	return nil
}

func ValidateConditionalRequired(condition bool, value, conditionDesc, fieldName string) error {
	if condition && value == "" {
		return errors.New(fieldName + " is required when " + conditionDesc)
	}
	return nil
}

func GenerateID(id *string) {
	if *id == "" {
		*id = ifs.NewUuid()
	}
}
