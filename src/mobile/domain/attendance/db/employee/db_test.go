package employee_test

import (
	"hroost/infrastructure/store/postgres"
	"hroost/mobile/domain/attendance/db/employee"
	"testing"
)

func TestNew(t *testing.T) {
	pgResolver := postgres.Resolver{}

	t.Run("empty config", func(t *testing.T) {
		db, err := employee.New(nil)
		if err == nil {
			t.Error("expct error got nil")
		}

		errMsg := "config required"
		if err.Error() != errMsg {
			t.Errorf("expect error message %s but got %s", errMsg, err.Error())
		}

		if db != nil {
			t.Errorf("expect db nil got %T", db)
		}
	})

	t.Run("empty pgResolver", func(t *testing.T) {
		db, err := employee.New(&employee.Config{})

		if err == nil {
			t.Error("expect error got nil")
		}
		errorMsg := "pgResolver required"
		if err.Error() != errorMsg {
			t.Errorf("expect errror message '%s' but got '%s'", errorMsg, err.Error())
		}
		if db != nil {
			t.Errorf("expect db nil but got %T", db)
		}
	})

	t.Run("valid config", func(t *testing.T) {
		db, err := employee.New(&employee.Config{
			PgResolver: &pgResolver,
		})
		if err != nil {
			t.Errorf("expect error nil but got %s", err)
		}
		var expectedDb *employee.Db
		if db == nil {
			t.Errorf("expect %T but got nil", expectedDb)
		}
	})
}
