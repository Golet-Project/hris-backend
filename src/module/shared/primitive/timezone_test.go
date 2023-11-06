package primitive_test

import (
	"hris/module/shared/primitive"
	"testing"
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
