package database

import (
	"os"

	"github.com/Protarium-Network/nano-assault-neo-nex/globals"
)

func initPostgres() {
	mustExec := func(label, query string) {
		if _, err := Postgres.Exec(query); err != nil {
			globals.Logger.Criticalf("%s: %s", label, err.Error())
			os.Exit(1)
		}
	}

	mustExec("Failed to create rankings table", `CREATE TABLE IF NOT EXISTS rankings (
        owner_pid integer,
        unique_id bigint,
        category integer,
        score integer,
        order_by integer,
        update_mode integer,
        groups bytea,
        param bigint,
        updated_at bigint,
        PRIMARY KEY (unique_id, owner_pid, category)
    )`)

	mustExec("Failed to create common_datas table", `CREATE TABLE IF NOT EXISTS common_datas (
        unique_id bigserial,
        owner_pid integer,
        common_data bytea,
        updated_at bigint,
        PRIMARY KEY (unique_id, owner_pid)
    )`)

	globals.Logger.Success("Postgres tables created")
}
