package database

import (
	"fmt"
	"time"

	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
	"github.com/Protarium-Network/nano-assault-neo-nex/globals"
)

func InsertRankingByPIDAndRankingScoreData(pid types.PID, rankingScoreData ranking_types.RankingScoreData, uniqueID types.UInt64) error {
	globals.Logger.Info(rankingScoreData.FormatToString(1))
	now := time.Now().Unix()

	var exists bool
	err := Postgres.QueryRow(`SELECT EXISTS(SELECT 1 FROM rankings WHERE owner_pid = $1 AND category = $2)`, uint32(pid), uint32(rankingScoreData.Category)).Scan(&exists)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	if !exists {
		_, err := Postgres.Exec(`
			INSERT INTO rankings (owner_pid, unique_id, category, score, order_by, update_mode, groups, param, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
			uint32(pid),
			uint64(uniqueID),
			uint32(rankingScoreData.Category),
			uint32(rankingScoreData.Score),
			uint32(rankingScoreData.OrderBy),
			uint32(rankingScoreData.UpdateMode),
			[]byte(rankingScoreData.Groups),
			int64(rankingScoreData.Param),
			now,
		)
		return err
	}

	_, err = Postgres.Exec(`
		UPDATE rankings SET score = $1, groups = $2, param = $3, updated_at = $4, order_by = $5, update_mode = $6
		WHERE category = $7 AND owner_pid = $8 AND unique_id = $9
	`,
		uint32(rankingScoreData.Score),
		[]byte(rankingScoreData.Groups),
		int64(rankingScoreData.Param),
		now,
		uint32(rankingScoreData.OrderBy),
		uint32(rankingScoreData.UpdateMode),
		uint32(rankingScoreData.Category),
		uint32(pid),
		uint64(uniqueID),
	)

	return err
}
