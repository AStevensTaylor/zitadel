package initialise

import (
	"database/sql"
	_ "embed"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
)

const (
	eventstoreSchema       = "eventstore"
	eventsTable            = "events"
	uniqueConstraintsTable = "unique_constraints"
	projectionsSchema      = "projections"
	systemSchema           = "system"
	encryptionKeysTable    = "encryption_keys"
)

var (
	searchSchema         = "SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1"
	searchTable          = "SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1"
	searchSystemSequence = "SELECT sequence_name FROM [SHOW SEQUENCES] WHERE sequence_name = 'system_seq'"

	//go:embed sql/04_eventstore.sql
	createEventstoreStmt string
	//go:embed sql/05_projections.sql
	createProjectionsStmt string
	//go:embed sql/06_system.sql
	createSystemStmt string
	//go:embed sql/07_encryption_keys_table.sql
	createEncryptionKeysStmt string
	//go:embed sql/08_enable_hash_sharded_indexes.sql
	enableHashShardedIdx string
	//go:embed sql/09_events_table.sql
	createEventsStmt string
	//go:embed sql/10_system_sequence.sql
	createSystemSequenceStmt string
	//go:embed sql/11_unique_constraints_table.sql
	createUniqueConstraints string
)

func newZitadel() *cobra.Command {
	return &cobra.Command{
		Use:   "zitadel",
		Short: "initialize ZITADEL internals",
		Long: `initialize ZITADEL internals.

Prereqesits:
- cockroachdb with user and database
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(Config)
			if err := viper.Unmarshal(config); err != nil {
				return err
			}
			return verifyZitadel(config.Database)
		},
	}
}

func VerifyZitadel(db *sql.DB) error {
	if err := verify(db, exists(searchSchema, systemSchema), exec(createSystemStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchTable, encryptionKeysTable), createEncryptionKeys); err != nil {
		return err
	}

	if err := verify(db, exists(searchSchema, projectionsSchema), exec(createProjectionsStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchSchema, eventstoreSchema), exec(createEventstoreStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchTable, eventsTable), createEvents); err != nil {
		return err
	}

	if err := verify(db, exists(searchSystemSequence), exec(createSystemSequenceStmt)); err != nil {
		return err
	}

	if err := verify(db, exists(searchTable, uniqueConstraintsTable), exec(createUniqueConstraints)); err != nil {
		return err
	}
	return nil
}

func verifyZitadel(config database.Config) error {
	logging.WithFields("database", config.Database).Info("verify zitadel")
	db, err := database.Connect(config)
	if err != nil {
		return err
	}
	if err := VerifyZitadel(db); err != nil {
		return nil
	}

	return db.Close()
}

func createEncryptionKeys(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(createEncryptionKeysStmt); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func createEvents(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(enableHashShardedIdx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(createEventsStmt); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
