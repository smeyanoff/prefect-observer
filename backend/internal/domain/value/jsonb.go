package value

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONB - тип для хранения JSONB в GORM
type JSONB map[string]interface{}

// Value преобразует JSONB в []byte для хранения в базе
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan преобразует []byte в JSONB при чтении из базы
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("неверный тип для JSONB: %T", value)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	*j = data
	return nil
}
