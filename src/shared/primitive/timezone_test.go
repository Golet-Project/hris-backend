package primitive_test

import (
	"hroost/shared/primitive"
	"testing"
	"time"
)

func TestTimezone_Value(t *testing.T) {
	t.Run("WIB", func(t *testing.T) {
		if primitive.WIB.Value() != 7 {
			t.Errorf("WIB should have value 7")
		}
	})

	t.Run("WITA", func(t *testing.T) {
		if primitive.WITA.Value() != 8 {
			t.Errorf("WITA should have value 8")
		}
	})

	t.Run("WIT", func(t *testing.T) {
		if primitive.WIT.Value() != 9 {
			t.Errorf("WIT should have value 9")
		}
	})
}

func TestTimezone_Valid(t *testing.T) {
	validValue := []int{7, 8, 9}

	for _, v := range validValue {
		if !primitive.Timezone(v).Valid() {
			t.Errorf("%d should be valid", v)
		}
	}

	invalidValue := []int{0, 1, 2, 3, 4, 5, 6, 10, 11, 12}
	for _, v := range invalidValue {
		if primitive.Timezone(v).Valid() {
			t.Errorf("%d should be invalid", v)
		}
	}
}

func TestTimezone_Now(t *testing.T) {
	t.Run("WIB", func(t *testing.T) {
		now, err := primitive.WIB.Now()

		loc := time.FixedZone("WIB", 7*60*60)
		now2 := time.Now().In(loc)

		if err != nil {
			t.Errorf("WIB should not return error")
		}

		if now.Location().String() != "WIB" {
			t.Errorf("WIB should have location WIB")
		}

		if now.Format("2006-01-02 15:04:05") != now2.Format("2006-01-02 15:04:05") {
			t.Errorf("WIB should have same time with now2")
		}
	})

	t.Run("WITA", func(t *testing.T) {
		now, err := primitive.WITA.Now()

		loc := time.FixedZone("WIT", 8*60*60)
		now2 := time.Now().In(loc)

		if err != nil {
			t.Errorf("WITA should not return error")
		}

		if now.Location().String() != "WITA" {
			t.Errorf("WITA should have location WITA")
		}

		if now.Format("2006-01-02 15:04:05") != now2.Format("2006-01-02 15:04:05") {
			t.Errorf("WIT should have same time with now2")
		}
	})

	t.Run("WIT", func(t *testing.T) {
		now, err := primitive.WIT.Now()

		loc := time.FixedZone("WITA", 9*60*60)
		now2 := time.Now().In(loc)

		if err != nil {
			t.Errorf("WIT should not return error")
		}

		if now.Location().String() != "WIT" {
			t.Errorf("WIT should have location WIT")
		}

		if now.Format("2006-01-02 15:04:05") != now2.Format("2006-01-02 15:04:05") {
			t.Errorf("WITA should have same time with now2")
		}
	})

	t.Run("invalid timezone", func(t *testing.T) {
		_, err := primitive.Timezone(0).Now()
		if err == nil {
			t.Errorf("invalid timezone should return error")
		}
	})
}
