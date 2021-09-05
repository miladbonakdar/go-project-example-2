package parser

import (
	"log"
	"strconv"
)

func ParseNumber(strValue string) (int, error) {
	number, err := strconv.Atoi(strValue)
	if err != nil {
		log.Printf("error in parsing number with string value of \"%s\" ,Error : %s", strValue, err.Error())
		return 0, err
	}
	return number, nil
}
