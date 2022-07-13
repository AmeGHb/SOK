package validator

import (
	"strconv"
)

func ValidateUser(request map[string]string) (int, string) {

	var err error

	_id := request["_id"]
	if _id == "" {
		return 422, "Id must be field"
	}

	_name := request["name"]
	if _name == "" {
		return 422, "Name must be field"
	}

	_, err = strconv.ParseFloat(request["balance"], 64)

	if err != nil {
		return 422, "Value should be a number"
	}

	return 200, ""

}

func ValidateId(request map[string]string) (int, string) {

	_id := request["_id"]
	if _id == "" {
		return 422, "Id must be field"
	}
	return 200, ""

}

func ValidateValues(request map[string]string) (int, string) {

	var err error

	_id := request["_id"]
	if _id == "" {
		return 422, "Id must be field"
	}

	_, err = strconv.ParseFloat(request["value"], 64)
	if err != nil {
		return 422, "Value should be a number"
	}

	sign := request["sign"]
	if (sign != "plus") || (sign != "minus") {
		return 422, "Sign should be plus or minus"
	}

	return 200, ""

}
