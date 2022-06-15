package constants

import (
	"github.com/uibricks/studio-engine/internal/pkg/env"
)

var (
	DbSchema string
)

const (
	TxSaveProject   string = "tx-save-project"
	TxSaveMapping          = "tx-save-mapping"
	TxCommitProject        = "tx-commit-project"

	TxRestoreProject   string = "tx-restore-project"
	TxRestoreMapping          = "tx-restore-mapping"
)

const (
	TypeRecentProjects   string = "recent"
	TypeTrashProjects          = "trash"
)

func init() {
	DbSchema = env.MustGet("db_schema")
}
