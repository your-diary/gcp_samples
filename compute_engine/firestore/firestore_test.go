package firestore

import (
	"compute_engine/config"
	"fmt"
	"math/rand"
	"testing"
)

const configFile = "../config.json"

func Test_01(t *testing.T) {
	config, err := config.New(configFile)
	if err != nil {
		t.Fatal(err)
	}

	db, err := New(config.Firestore)
	if err != nil {
		t.Fatal(err)
	}

	content := fmt.Sprintf("%v", rand.Int())

	rows, err := db.selectByContent(content)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Insert(content)
	if err != nil {
		t.Fatal(err)
	}

	rowsAfter, err := db.selectByContent(content)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(rows)
	fmt.Println(rowsAfter)
	if len(rows)+1 != len(rowsAfter) {
		t.Fatal()

	}

}
