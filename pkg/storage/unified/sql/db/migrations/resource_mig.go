package migrations

import (
	"fmt"

	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

func initResourceTables(mg *migrator.Migrator) string {
	marker := "Initialize resource tables"
	mg.AddMigration(marker, &migrator.RawSQLMigration{})

	resource_table := migrator.Table{
		Name: "resource",
		Columns: []*migrator.Column{
			// primary identifier
			{Name: "guid", Type: migrator.DB_NVarchar, Length: 36, Nullable: false, IsPrimaryKey: true},

			{Name: "resource_version", Type: migrator.DB_BigInt, Nullable: true},

			// K8s Identity group+(version)+namespace+resource+name
			{Name: "group", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "resource", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "namespace", Type: migrator.DB_NVarchar, Length: 63, Nullable: false},
			{Name: "name", Type: migrator.DB_NVarchar, Length: 253, Nullable: false},
			{Name: "value", Type: migrator.DB_LongText, Nullable: true},
			{Name: "action", Type: migrator.DB_Int, Nullable: false}, // 1: create, 2: update, 3: delete

			// Hashed label set
			{Name: "label_set", Type: migrator.DB_NVarchar, Length: 64, Nullable: true}, // null is no labels
		},
		Indices: []*migrator.Index{
			{Cols: []string{"namespace", "group", "resource", "name"}, Type: migrator.UniqueIndex},
		},
	}
	resource_history_table := migrator.Table{
		Name: "resource_history",
		Columns: []*migrator.Column{
			// primary identifier
			{Name: "guid", Type: migrator.DB_NVarchar, Length: 36, Nullable: false, IsPrimaryKey: true},
			{Name: "resource_version", Type: migrator.DB_BigInt, Nullable: true},

			// K8s Identity group+(version)+namespace+resource+name
			{Name: "group", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "resource", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "namespace", Type: migrator.DB_NVarchar, Length: 63, Nullable: false},
			{Name: "name", Type: migrator.DB_NVarchar, Length: 253, Nullable: false},
			{Name: "value", Type: migrator.DB_LongText, Nullable: true},
			{Name: "action", Type: migrator.DB_Int, Nullable: false}, // 1: create, 2: update, 3: delete

			// Hashed label set
			{Name: "label_set", Type: migrator.DB_NVarchar, Length: 64, Nullable: true}, // null is no labels
		},
		Indices: []*migrator.Index{
			{
				Cols: []string{"namespace", "group", "resource", "name", "resource_version"},
				Type: migrator.UniqueIndex,
				Name: "UQE_resource_history_namespace_group_name_version",
			},
			// index to support watch poller
			{Cols: []string{"resource_version"}, Type: migrator.IndexType},
		},
	}

	resource_version_table := migrator.Table{
		Name: "resource_version",
		Columns: []*migrator.Column{
			{Name: "group", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "resource", Type: migrator.DB_NVarchar, Length: 190, Nullable: false},
			{Name: "resource_version", Type: migrator.DB_BigInt, Nullable: false},
		},
		Indices: []*migrator.Index{
			{Cols: []string{"group", "resource"}, Type: migrator.UniqueIndex},
		},
	}

	tables := []migrator.Table{resource_table, resource_history_table, resource_version_table}

	// Initialize all tables
	for t := range tables {
		mg.AddMigration("drop table "+tables[t].Name, migrator.NewDropTableMigration(tables[t].Name))
		mg.AddMigration("create table "+tables[t].Name, migrator.NewAddTableMigration(tables[t]))
		for i := range tables[t].Indices {
			mg.AddMigration(fmt.Sprintf("create table %s, index: %d", tables[t].Name, i), migrator.NewAddIndexMigration(tables[t], tables[t].Indices[i]))
		}
	}

	mg.AddMigration("Add column previous_resource_version in resource_history", migrator.NewAddColumnMigration(resource_history_table, &migrator.Column{
		Name: "previous_resource_version", Type: migrator.DB_BigInt, Nullable: true,
	}))

	mg.AddMigration("Add column previous_resource_version in resource", migrator.NewAddColumnMigration(resource_table, &migrator.Column{
		Name: "previous_resource_version", Type: migrator.DB_BigInt, Nullable: true,
	}))

	mg.AddMigration("Add index to resource_history for polling", migrator.NewAddIndexMigration(resource_history_table, &migrator.Index{
		Cols: []string{"group", "resource", "resource_version"}, Type: migrator.IndexType,
	}))

	mg.AddMigration("Add index to resource for loading", migrator.NewAddIndexMigration(resource_table, &migrator.Index{
		Cols: []string{"group", "resource"}, Type: migrator.IndexType,
	}))

	mg.AddMigration("Add column shard in resource", migrator.NewAddColumnMigration(resource_version_table, &migrator.Column{
		Name: "shard", Type: migrator.DB_Int, Nullable: false, Default: "0",
	}))

	mg.AddMigration("Add index to resource_version", migrator.NewAddIndexMigration(resource_version_table, &migrator.Index{
		Cols: []string{"group", "resource", "shard"}, Type: migrator.IndexType,
	}))

	return marker
}
