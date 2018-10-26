package xlxstofirebase

import (
	"errors"
	"log"

	firebase "firebase.google.com/go"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

// XlsxToFirebase - XlsxToFirebase struct
type XlsxToFirebase struct {
	App     *firebase.App
	Context context.Context
}

// SetupFirebase - Setup firebase
func SetupFirebase(config *Config) (*XlsxToFirebase, error) {
	ctx := context.Background()

	sa := option.WithCredentialsFile(config.Credentials)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	var xtb XlsxToFirebase
	xtb.App = app
	xtb.Context = ctx

	return &xtb, nil
}

// SeedDatabase - Seed the firebase database with xlsx
func (xtb *XlsxToFirebase) SeedDatabase(config *Config) error {
	client, err := xtb.App.Firestore(xtb.Context)
	if err != nil {
		return err
	}

	xlsx, err := excelize.OpenFile(config.XlsxFile)
	if err != nil {
		log.Fatalf("excelize: error opening file: %v\n", err)
		return err
	}

	for _, sheet := range config.XlsxSheets {

		batch := client.Batch()

		for rowIndex, row := range xlsx.GetRows(sheet.Name) {
			// assuming first row is just headers
			if rowIndex == 0 {
				continue
			}

			recordID := createID()

			record, err := createRecord(recordID, sheet, row)
			if err != nil {
				log.Fatalf("xtb: error creating record: %v\n", err)
				return err
			}

			newRef := client.Collection(sheet.FirebaseDbRef).Doc(recordID)
			batch.Set(newRef, record)
		}

		_, err := batch.Commit(xtb.Context)
		if err != nil {
			log.Print("firebase: error commiting batch: \n", err)
			return err
		}
	}

	defer client.Close()

	return nil
}

func createRecord(recordID string, sheet XlsxSheet, row []string) (map[string]interface{}, error) {
	newRecordItem := make(map[string]interface{})

	for colIndex, colCell := range row {
		// check required cells from config
		if IntInArray(colIndex, sheet.RequiredCellIndexes) && colCell == "" {
			log.Print("excelize: required cells are missing: \n")
			return nil, errors.New("excelize: required cells are missing")
		}

		// check excluded cells from config
		if IntInArray(colIndex, sheet.ExcludedCellIndexes) || colCell == "" {
			continue
		}

		if colIndex == 0 {
			newRecordItem["id"] = recordID
		}

		switch sheet.Columns[colIndex]["type"] {
		case "string":
			if sheet.Columns[colIndex]["name"] != "" {
				newRecordItem[sheet.Columns[colIndex]["name"].(string)] = colCell
			} else {
				newRecordItem[sheet.Columns[colIndex]["name"].(string)] = ""
			}
		case "bool":
			if sheet.Columns[colIndex]["name"] == "True" || sheet.Columns[colIndex]["name"] == "true" || sheet.Columns[colIndex]["name"] == "Yes" || sheet.Columns[colIndex]["name"] == "yes" {
				newRecordItem[sheet.Columns[colIndex]["name"].(string)] = true
			} else {
				newRecordItem[sheet.Columns[colIndex]["name"].(string)] = false
			}
		}
	}

	return newRecordItem, nil
}

func createID() string {
	return xid.New().String()
}
