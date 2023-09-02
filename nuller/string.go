package nuller

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type NullString struct {
	sql.NullString
}

func (n NullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.String)
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	var strData string

	err := json.Unmarshal(data, &strData)
	if err != nil {
		return err
	}

	*n = NullString{
		sql.NullString{
			Valid:  !(strData == "null" || len(strData) == 0),
			String: strData,
		},
	}

	return nil
}

func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return string(n.String), nil
}

func (n *NullString) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not valid string")
	}

	n.String = str

	return nil
}
