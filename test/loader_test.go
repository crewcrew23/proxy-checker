package test

import (
	"os"
	"testing"

	"github.com/crewcrew23/proxy-checker/internal/loader"
)

func TestLoader(t *testing.T) {
	tests := []struct {
		Name            string
		FilePath        string
		Create          bool
		ExpectedReed    bool
		TestData        string
		ExpectedProblem bool
	}{
		{
			Name:            "CREATE AND LOADE FILE",
			FilePath:        "./test.txt",
			Create:          true,
			ExpectedReed:    false,
			ExpectedProblem: false,
		},
		{
			Name:            "NOT CREATED FILE AND LOADE",
			FilePath:        "./test.txt",
			Create:          false,
			ExpectedReed:    false,
			ExpectedProblem: true,
		},
		{
			Name:            "CREATE AND REED FILE",
			FilePath:        "./test.txt",
			Create:          true,
			ExpectedReed:    true,
			TestData:        "test data",
			ExpectedProblem: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_ = os.Remove(tt.FilePath)

			if tt.Create {
				file, err := os.Create(tt.FilePath)
				if err != nil {
					t.Fatalf("error creating file: %v", err)
				}

				if tt.ExpectedReed {
					_, err := file.WriteString(tt.TestData)
					if err != nil {
						file.Close()
						t.Errorf("error of write to the file")
					}

				}
				file.Close()
			}

			t.Cleanup(func() {
				_ = os.Remove(tt.FilePath)
			})

			val, err := loader.LoadProxies(tt.FilePath)
			if err == nil && tt.ExpectedReed && val[0] != tt.TestData {
				t.Error("correct data into the file")
			}

			if tt.ExpectedProblem && err == nil {
				t.Errorf("expected error, got none")
			} else if !tt.ExpectedProblem && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
