package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"github.com/Protarium-Network/nano-assault-neo-nex/globals"
	"github.com/lib/pq"
)

func GetRankingsAndCountByCategoryAndRankingOrderParam(category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error) {
	var pqErr *pq.Error
	globals.Logger.Info(rankingOrderParam.FormatToString(1))

	rows, err := Postgres.Query(`
		SELECT owner_pid, unique_id, category, score, order_by, update_mode, groups, param, updated_at 
		FROM "rankings" 
		WHERE category = $1
	`, uint32(category))

	if errors.Is(err, sql.ErrNoRows) || (errors.As(err, &pqErr) && pqErr.SQLState() == "42P01") {
		fmt.Printf("No rows")
		return types.List[ranking_types.RankingRankData]{}, 0, err
	}
	if err != nil {
		fmt.Printf("Error selecting rankings data.")
		return types.List[ranking_types.RankingRankData]{}, 0, err
	}
	defer rows.Close()

	resultsPtr, err := parseRankingDataList(rows)
	if err != nil {
		fmt.Printf("Error parsing ranking data")
		return types.List[ranking_types.RankingRankData]{}, 0, err
	}

	var results types.List[ranking_types.RankingRankData]
	for _, item := range resultsPtr {
		if item != nil {
			results = append(results, *item)
		}
	}

	globals.Logger.Info(results.String())

	return results, uint32(len(results)), nil
}
