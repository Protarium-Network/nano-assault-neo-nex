package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

func parseRankingDataList(rows *sql.Rows) (types.List[*ranking_types.RankingRankData], error) {
	var results types.List[*ranking_types.RankingRankData]

	for rows.Next() {
		result := ranking_types.NewRankingRankData()
		var updateDate int64
		var userPid uint64
		var updateMode uint32

		err := rows.Scan(
			&userPid,
			&result.UniqueID,
			&result.Category,
			&result.Score,
			&result.Order,
			&updateMode,
			&result.Groups,
			&result.Param,
			&updateDate,
		)
		if err != nil {
			fmt.Printf("Error scanning the row: %v\n", err)
			return nil, err
		}

		result.PrincipalID = types.NewPID(userPid)
		result.UpdateTime.FromTimestamp(time.Unix(updateDate, 0))

		var commonData []byte
		err = Postgres.QueryRow(`SELECT "common_data" FROM "common_datas" WHERE "owner_pid" = $1`, userPid).Scan(&commonData)
		if err != nil && err != sql.ErrNoRows {
			fmt.Printf("Couldn't find common data: %v\n", err)
			return nil, err
		}

		if len(commonData) > 0 {
			result.CommonData = types.NewBuffer(commonData)
		}

		results = append(results, &result)
	}

	return results, nil
}
