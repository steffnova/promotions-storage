package csv

import (
	"os"
	"testing"
)

func TestIntegrationFileReader(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"FileReader": func(t *testing.T) {
			fileName := "testFile"
			if _, err := os.Create(fileName); err != nil {
				t.Fatalf("Unexpecte error: %s", err)
			}
			defer os.Remove(fileName)

			if _, err := FileReader(fileName)(); err != nil {
				t.Fatalf("Unexpecte error: %s", err)
			}
		},
		"ErrorHandling": func(t *testing.T) {
			fileName := "random_file_name"
			if _, err := FileReader(fileName)(); err == nil {
				t.Fatalf("Expected an error. It seems that %s file exists", fileName)
			}
		},
	}

	for name, testCase := range testCases {
		t.Run(name, testCase)
	}

}
