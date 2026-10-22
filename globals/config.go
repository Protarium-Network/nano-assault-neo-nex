package globals

// NEX configuration for "Nano Assault Neo" (Wii U, Comcept / Deep Silver, 2016).
//
// GameServerID and AccessKey come from kinnay.github.io's public Wii U NEX
// database (id 270,... = 0x101c5a00). The retail executable (MN9Game.rpx, an
// Unreal Engine 3 build using Nintendo's OnlineWiiUNexSystem) reads the game
// server ID from a field of its online-system object, not from a literal, so
// whether other regions use a different ID is NOT settled - watch the token
// service logs for other game_server_id values.
const (
	GameServerID = "10110100"
	AccessKey    = "7e484a8e"

	NEXMajor = 3
	NEXMinor = 0
	NEXPatch = 1
)

var LegacyConnectionSignature = true
