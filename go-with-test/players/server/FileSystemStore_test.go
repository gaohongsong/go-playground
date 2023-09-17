package server

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	fmt.Printf("create tmpfile %s\n", tmpfile.Name())

	if err != nil {
		t.Fatalf("create tempfile failed: %v", err)
	}

	tmpfile.Write([]byte(initialData))

	return tmpfile, func() {
		fmt.Printf("delete tmpfile %s\n", tmpfile.Name())
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}
}

func TestFileSystemStore(t *testing.T) {
	var recordsJson = `[
			{"name": "Cleo", "Wins": 10},
			{"name": "Chris", "Wins": 33}]`

	t.Run("/league from a reader", func(t *testing.T) {

		//database := strings.NewReader(recordsJson)

		database, cleanDatabase := createTempFile(t, recordsJson)
		defer cleanDatabase()

		store := NewFileSystemStore(database)
		//store := FileSystemStore{database}

		got := store.GetLeague()
		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		// read again
		got = store.GetLeague()
		assertLeague(t, want, got)
	})

	t.Run("get player score", func(t *testing.T) {
		//database := strings.NewReader(recordsJson)
		database, cleanDatabase := createTempFile(t, recordsJson)
		defer cleanDatabase()

		//store := FileSystemStore{database}
		store := NewFileSystemStore(database)

		got := store.GetPlayerScore("Chris")
		want := 33
		assertScoreEquals(t, want, got)
	})

	t.Run("store player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, recordsJson)
		defer cleanDatabase()

		//store := FileSystemStore{database}
		store := NewFileSystemStore(database)
		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, want, got)
	})

	t.Run("store new player", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, recordsJson)
		defer cleanDatabase()

		store := NewFileSystemStore(database)
		//store := FileSystemStore{database}
		store.RecordWin("Petter")

		got := store.GetPlayerScore("Petter")
		want := 1
		assertScoreEquals(t, want, got)
	})

}
