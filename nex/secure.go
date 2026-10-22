package nex

import (
	"fmt"
	"os"
	"strconv"

	nexgo "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	commonRanking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	commonSecure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	utility "github.com/PretendoNetwork/nex-protocols-go/v2/utility"

	"github.com/Protarium-Network/nano-assault-neo-nex/database"
	"github.com/Protarium-Network/nano-assault-neo-nex/globals"
)

var SecureServer *nexgo.PRUDPServer
var SecureEndpoint *nexgo.PRUDPEndPoint

func StartSecureServer() {
	SecureServer = nexgo.NewPRUDPServer()

	SecureServer.PRUDPV1Settings.LegacyConnectionSignature = globals.LegacyConnectionSignature

	SecureServer.LibraryVersions.SetDefault(nexgo.NewLibraryVersion(globals.NEXMajor, globals.NEXMinor, globals.NEXPatch))
	SecureServer.AccessKey = globals.AccessKey

	SecureEndpoint = nexgo.NewPRUDPEndPoint(1)
	SecureEndpoint.IsSecureEndPoint = true
	SecureEndpoint.ServerAccount = globals.SecureServerAccount
	SecureEndpoint.AccountDetailsByPID = globals.AccountDetailsByPID
	SecureEndpoint.AccountDetailsByUsername = globals.AccountDetailsByUsername

	SecureServer.BindPRUDPEndPoint(SecureEndpoint)
	globals.SecureEndpoint = SecureEndpoint

	SecureEndpoint.OnData(func(packet nexgo.PacketInterface) {
		request := packet.RMCMessage()
		if request == nil {
			return
		}

		fmt.Printf("[NAN Secure] PID=%d protocol=0x%02X method=0x%02X\n", uint64(packet.Sender().PID()), request.ProtocolID, request.MethodID)
	})

	SecureEndpoint.OnConnectionEnded(func(connection *nexgo.PRUDPConnection) {
		fmt.Printf("[NAN Secure] PID=%d disconnected\n", uint64(connection.PID()))
	})

	SecureEndpoint.OnError(func(err *nexgo.Error) {
		globals.Logger.Errorf("[NAN Secure] %v", err)
	})

	secureProtocol := secure.NewProtocol()
	SecureEndpoint.RegisterServiceProtocol(secureProtocol)

	secureCommon := commonSecure.NewCommonProtocol(secureProtocol)
	secureCommon.EnableInsecureRegister()

	utilityProtocol := utility.NewProtocol()
	utilityProtocol.AcquireNexUniqueID = func(err error, packet nexgo.PacketInterface, callID uint32) (*nexgo.RMCMessage, *nexgo.Error) {
		if err != nil {
			globals.Logger.Error(err.Error())
			return nil, nexgo.NewError(nexgo.ResultCodes.Core.InvalidArgument, err.Error())
		}

		connection := packet.Sender().(*nexgo.PRUDPConnection)
		endpoint := connection.Endpoint().(*nexgo.PRUDPEndPoint)

		uniqueID := types.NewUInt64(uint64(connection.PID()))

		stream := nexgo.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())
		uniqueID.WriteTo(stream)

		rmcResponse := nexgo.NewRMCSuccess(endpoint, stream.Bytes())
		rmcResponse.ProtocolID = utility.ProtocolID
		rmcResponse.MethodID = utility.MethodAcquireNexUniqueID
		rmcResponse.CallID = callID

		return rmcResponse, nil
	}
	SecureEndpoint.RegisterServiceProtocol(utilityProtocol)

	rankingProtocol := ranking.NewProtocol()
	rankingCommon := commonRanking.NewCommonProtocol(rankingProtocol)

	rankingCommon.GetRankingsAndCountByCategoryAndRankingOrderParam = database.GetRankingsAndCountByCategoryAndRankingOrderParam
	rankingCommon.UploadCommonData = database.UploadCommonData
	rankingCommon.InsertRankingByPIDAndRankingScoreData = database.InsertRankingByPIDAndRankingScoreData

	SecureEndpoint.RegisterServiceProtocol(rankingProtocol)

	port, err := strconv.Atoi(os.Getenv("PN_NAN_SECURE_PORT"))
	if err != nil || port < 1 || port > 65535 {
		port = 60000
	}

	globals.Logger.Successf("[NAN] Secure server listening on UDP %d", port)
	SecureServer.Listen(port)
}
