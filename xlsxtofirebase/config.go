package xlxstofirebase

import (
	"encoding/json"
	"io/ioutil"
)

// Config - config struct
type Config struct {
	Credentials       string      `json:"credentials"`
	FirebaseProjectID string      `json:"firebase_project_id"`
	XlsxFile          string      `json:"xlsx_file"`
	XlsxSheets        []XlsxSheet `json:"xlsx_sheets"`
}

// XlsxSheet - sheet config struct
type XlsxSheet struct {
	Name                string                   `json:"name"`
	RequiredCellIndexes []int                    `json:"required_cell_indexes"`
	ExcludedCellIndexes []int                    `json:"excluded_cell_indexes"`
	Columns             []map[string]interface{} `json:"columns"`
	FirebaseDbRef       string                   `json:"firebase_db_ref"`
}

// SetupConfig - Setup config
func SetupConfig(file string) (*Config, error) {
	config, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var c *Config
	err = json.Unmarshal(config, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
