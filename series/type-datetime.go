package series

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

type dateTimeElement struct {
	e   time.Time
	nan bool
}

// force dateTimeElement struct to implement Element interface
var _ Element = (*dateTimeElement)(nil)

func (e *dateTimeElement) Set(value interface{}) {
	var err error
	e.nan = false
	switch val := value.(type) {
	case string:
		if val == "NaN" {
			e.nan = true
			return
		}
		e.e, err = time.Parse(time.RFC3339, val)
		if err != nil {
			e.nan = true
			return
		}
	case int:
		// e.e = time.UnixMilli(reflect.ValueOf(val).Int())
		e.e = time.UnixMilli(int64(val))
	case float64:
		f := val
		if math.IsNaN(f) ||
			math.IsInf(f, 0) ||
			math.IsInf(f, 1) {
			e.nan = true
			return
		}
		e.e = time.UnixMilli(int64(f))
	case bool:
		e.e = time.Time{}
	case Element:
		e.e = time.UnixMilli(reflect.ValueOf(val).Int())
		// if err != nil {
		// 	e.nan = true
		// 	return
		// }
	default:
		e.nan = true
		return
	}
}

func (e dateTimeElement) Copy() Element {
	if e.IsNA() {
		return &dateTimeElement{time.Time{}, true}
	}
	return &dateTimeElement{e.e, false}
}

func (e dateTimeElement) IsNA() bool {
	return e.nan
}

func (e dateTimeElement) Type() Type {
	return DateTime
}

func (e dateTimeElement) Val() ElementValue {
	if e.IsNA() {
		return nil
	}
	return e.e
}

func (e dateTimeElement) String() string {
	if e.IsNA() {
		return "NaN"
	}
	return fmt.Sprintf("%s", e.e.Format(time.RFC3339))
}

func (e dateTimeElement) Int() (int, error) {
	if e.IsNA() {
		return 0, fmt.Errorf("can't convert NaN to int")
	}
	return int(e.e.UnixMilli()), nil
}

func (e dateTimeElement) Float() float64 {
	if e.IsNA() {
		return math.NaN()
	}
	return float64(e.e.UnixMilli())
}

func (e dateTimeElement) Bool() (bool, error) {
	return false, fmt.Errorf("can't convert DateTime \"%v\" to bool", e.e)
}

func (e dateTimeElement) Eq(elem Element) bool {
	return e.e == elem.Val()
}

func (e dateTimeElement) Neq(elem Element) bool {
	return e.e != elem.Val()
}

func (e dateTimeElement) Less(elem Element) bool {
	if e.e.Before(elem.Val().(time.Time)) {
		return true
	}
	return false
}

func (e dateTimeElement) LessEq(elem Element) bool {
	return !e.e.After(elem.Val().(time.Time))
}

func (e dateTimeElement) Greater(elem Element) bool {
	return e.e.After(elem.Val().(time.Time))
}

func (e dateTimeElement) GreaterEq(elem Element) bool {
	return !e.e.Before(elem.Val().(time.Time))
}
