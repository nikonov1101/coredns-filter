package piholedb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// PiHoleDatabaseReader reads blocked domains from the gravity.db,
// used in the Pi Hole project https://docs.pi-hole.net/database/gravity/
type PiHoleDatabaseReader struct {
	db *sql.DB
}

func New(path string) (*PiHoleDatabaseReader, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	// open in readonly mode, with a large cache and no journal: for maximum reading performance.
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=ro&_cache_size=1000000&immutable=true&_journal_mode=OFF", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("gravity: failed to open database at %q: %w", path, err)
	}

	return &PiHoleDatabaseReader{db: db}, nil
}

func (ftl *PiHoleDatabaseReader) IsBlocked(name string) bool {
	started := time.Now()
	defer func() {
		gravityLookupDuration.Set(time.Since(started).Seconds())
	}()

	name = normalizeDomain(name)

	row := ftl.db.QueryRow("select 1 from gravity where domain = ?", name)
	var count int64 = -1
	if err := row.Scan(&count); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("gravity: lookup failure: unexpected error: name=%s, err=%v\n", name, err)
		}
		return false
	}

	return true
}

func normalizeDomain(s string) string {
	// get rid of trailing dot, if any
	if s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}
