package model

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	phoneRegex = regexp.MustCompile("010-?[\\d]{4}-?[\\d]{4}")
)

var (
	ErrNotCorrectInput   = fmt.Errorf("Not Correct Input")
	ErrInputMustBeString = fmt.Errorf("Input Must Be String")
)

func isPhoneValid(p string) bool {
	if len(p) < 8 && len(p) > 12 {
		return false
	}
	return phoneRegex.MatchString(p)
}

type Phone string

func (p Phone) MarshalGQL(w io.Writer) {
	fmt.Println("Marshal: ", p)
	io.WriteString(w, strconv.Quote(string(p)))
}

func (p *Phone) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		if ok := isPhoneValid(v); !ok {
			return ErrNotCorrectInput
		}
		*p = Phone(strings.ReplaceAll(v, "-", ""))
		return nil
	case int:
		str := strconv.Itoa(v)
		if ok := isPhoneValid(str); !ok {
			return ErrNotCorrectInput
		}
		*p = Phone(str)
		return nil
	default:
		return ErrNotCorrectInput
	}
}

type ObjectID primitive.ObjectID

func (o ObjectID) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(primitive.ObjectID(o).Hex()))
}

func (o *ObjectID) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		result, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return ErrNotCorrectInput
		}
		*o = ObjectID(result)
		return nil
	default:
		return ErrNotCorrectInput
	}
}