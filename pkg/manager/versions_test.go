package manager

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestWriteCurrentVersion(t *testing.T) {
	targetTag := "boom"
	targetBranch := "ch666"

	WriteCurrentVersion(targetTag, targetBranch)
	defer os.Remove(".env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	writtenTag := os.Getenv("CURRENT_TAG")
	writtenBranch := os.Getenv("CURRENT_BRANCH")

	if writtenTag != targetTag {
		t.Errorf("got tag %q want %q", writtenTag, targetTag)
	}

	if writtenBranch != targetBranch {
		t.Errorf("got branch %q want %q", writtenBranch, targetBranch)
	}
}
