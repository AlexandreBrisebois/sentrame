// migrate-client is an explicit offline operation, never called by service startup.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sentrame/internal/app"
)

func main() {
	db := flag.String("db", "", "explicit managed data-directory database path (required)")
	apply := flag.Bool("apply", false, "apply or resume the staged client-only migration")
	backup := flag.Bool("backup-confirmed", false, "confirm an independent current private coordinated backup")
	flag.Parse()
	if *db == "" {
		fmt.Fprintln(os.Stderr, "--db is required; default operation is offline preflight")
		os.Exit(2)
	}
	report, err := app.MigrateClientRegistry(*db, *apply, *backup)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	json.NewEncoder(os.Stdout).Encode(report)
}
