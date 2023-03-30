package primitive

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"
)

type String sql.NullString

func (s *String) Scan(v interface{}) error {
	var ns sql.NullString
	if err := ns.Scan(v); err != nil {
		return err
	}

	if reflect.TypeOf(v) == nil {
		*s = String{String: ns.String, Valid: false}
	} else {
		*s = String{String: ns.String, Valid: true}
	}
	return nil
}

func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String, nil
}

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *String) UnmarshalJSON(by []byte) error {
	if str := string(by); str == `null` {
		return nil
	}
	err := json.Unmarshal(by, &s.String)
	s.Valid = (err == nil)
	return err
}

type Int64 sql.NullInt64

func (i *Int64) Scan(v interface{}) error {
	var ni sql.NullInt64
	if err := ni.Scan(v); err != nil {
		return err
	}

	if reflect.TypeOf(v) == nil {
		*i = Int64{Int64: ni.Int64, Valid: false}
	} else {
		*i = Int64{Int64: ni.Int64, Valid: true}
	}
	return nil
}

func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Int64)
}

func (i *Int64) UnmarshalJSON(by []byte) error {
	if str := string(by); str == `null` {
		return nil
	}
	err := json.Unmarshal(by, &i.Int64)
	i.Valid = (err == nil)
	return err
}

type Time sql.NullTime

func (t *Time) Scan(v interface{}) error {
	var nt sql.NullTime
	if err := nt.Scan(v); err != nil {
		return err
	}

	if reflect.TypeOf(v) == nil {
		*t = Time{Time: nt.Time, Valid: false}
	} else {
		*t = Time{Time: nt.Time, Valid: true}
	}
	return nil
}

func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

func (t *Time) UnmarshalJSON(by []byte) error {
	if str := string(by); str == `null` {
		return nil
	}
	err := json.Unmarshal(by, &t.Time)
	t.Valid = (err == nil)
	return err
}

type Date sql.NullString

func (d *Date) Scan(v interface{}) error {
	var nd sql.NullString
	if err := nd.Scan(v); err != nil {
		return err
	}

	if reflect.TypeOf(v) == nil {
		*d = Date{String: nd.String, Valid: false}
	} else {
		*d = Date{String: nd.String, Valid: true}
	}
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.String, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}
	dt, err := time.Parse("2006-01-02", d.String)
	if err != nil {
		return nil, err
	}

	return json.Marshal(dt.Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(by []byte) error {
	if str := string(by); str == `null` {
		return nil
	}
	err := json.Unmarshal(by, &d.String)
	d.Valid = (err == nil)
	return err
}

type Bool sql.NullBool

func (b *Bool) Scan(v interface{}) error {
	var nb sql.NullBool
	if err := nb.Scan(v); err != nil {
		return err
	}

	if reflect.TypeOf(v) == nil {
		*b = Bool{Bool: nb.Bool, Valid: false}
	} else {
		*b = Bool{Bool: nb.Bool, Valid: true}
	}
	return nil
}

type Float64 sql.NullFloat64

func (f *Float64) Scan(v interface{}) error {
	var nf sql.NullFloat64
	if err := nf.Scan(v); err != nil {
		return nil
	}

	if reflect.TypeOf(v) == nil {
		*f = Float64{Float64: nf.Float64, Valid: false}
	} else {
		*f = Float64{Float64: nf.Float64, Valid: true}
	}
	return nil
}

func (f Float64) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.Float64, nil
}

func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(f.Float64)
}

func (f *Float64) UnmarshalJSON(by []byte) error {
	if str := string(by); str == `null` {
		return nil
	}
	err := json.Unmarshal(by, &f.Float64)
	f.Valid = (err == nil)
	return err
}
