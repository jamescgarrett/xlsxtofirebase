# XLSX To Firebase Store

## Usage

1) Generate and download your Firebase private key file [Instructions Here](https://firebase.google.com/docs/admin/setup)

2) Create config JSON file
```json
{
  "credentials": "./path/to/your/private/key-file", // File downloaded in step 1
  "firebase_project_id": "FIREBASE_PROJECT_ID", // can be found in file downloaded in step 1
  "xlsx_file": "./path/to/your/xlsxfile.xlsx",
  "xlsx_sheets": [ // create an object for each sheet
    {
      "name": "SheetName1",
      "required_cell_indexes": [0], // Index of required row cells
      "excluded_cell_indexes": [12], // Index of row cells to skip
      "columns": [ // map your column names and define the value type
        { "name": "propName1", "type": "string" },
        { "name": "propName2", "type": "string" },
        { "name": "propName3", "type": "bool" },
        { "name": "propName4", "type": "bool" },
        { "name": "propName5", "type": "string" },
      ],
      "firebase_db_ref": "CollectionName1"
    },
    {
      "name": "SheetName2",
      "required_cell_indexes": [0],
      "excluded_cell_indexes": [],
      "columns": [
        { "name": "propName1", "type": "string" },
        { "name": "propName2", "type": "bool" }
      ],
      "firebase_db_ref": "Locations"
    }
  ]
}
```
3) Setup your go script: `main.go`
```go
package main

import (
	"flag"
	"log"
	xtb "xlsxtofirebase/xlsxtofirebase" // not published yet,this is using local copy
)

func main() {
	flagConfig := flag.String("config", "../config.yaml", "Path to config") // flag - path to your config file
	flag.Parse()

	config, err := xtb.SetupConfig(*flagConfig)
	if err != nil {
		log.Fatalf("config: %v\n", err)
	}

	firebase, err := xtb.SetupFirebase(config)
	if err != nil {
		log.Fatalf("firebase: %v\n", err)
	}

	err = firebase.SeedDatabase(config)
	if err != nil {
		log.Fatalf("firebase: %v\n", err)
	}
}
```

4) Run the script:
 ```
 go run main.go -config=./config/config.json
 ```

## Caveats
- Currently need to trim all empty columns and rows from the spreadsheet
- First row should be used as headers (will be skipped)
- Bool values are checked for one of the following to determine value: `True`, `true`, `Yes` or `yes`
- Currently only available for Firestore

## To Do
- Add Relationships
- Add more value types
- Configure for Firebase Realtime Database
- Handle empty columns/rows gracefully
