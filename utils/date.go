package utils

import (
	"database/sql"
	"encoding/json"
	"time"
)

type JsonDate struct {
	sql.NullTime
}

func (v JsonDate) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time.Format("2006-01-02"))
	} else {
		return json.Marshal("")
	}
}

func (v *JsonDate) UnmarshalJSON(data []byte) error {
	var x *string
	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	if x != nil && *x != "" {
		v.Time, err = time.Parse("2006-01-02", *x)
		if err != nil {
			return err
		}
		v.Valid = true
	} else {
		v.Valid = false
	}
	return nil
}
