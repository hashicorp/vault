package memd

import "encoding/hex"

// CmdCode represents the specific command the packet is performing.
type CmdCode uint8

// These constants provide predefined values for all the operations
// which are supported by this library.
const (
	CmdGet                    = CmdCode(0x00)
	CmdSet                    = CmdCode(0x01)
	CmdAdd                    = CmdCode(0x02)
	CmdReplace                = CmdCode(0x03)
	CmdDelete                 = CmdCode(0x04)
	CmdIncrement              = CmdCode(0x05)
	CmdDecrement              = CmdCode(0x06)
	CmdNoop                   = CmdCode(0x0a)
	CmdAppend                 = CmdCode(0x0e)
	CmdPrepend                = CmdCode(0x0f)
	CmdStat                   = CmdCode(0x10)
	CmdTouch                  = CmdCode(0x1c)
	CmdGAT                    = CmdCode(0x1d)
	CmdHello                  = CmdCode(0x1f)
	CmdSASLListMechs          = CmdCode(0x20)
	CmdSASLAuth               = CmdCode(0x21)
	CmdSASLStep               = CmdCode(0x22)
	CmdGetAllVBSeqnos         = CmdCode(0x48)
	CmdDcpOpenConnection      = CmdCode(0x50)
	CmdDcpAddStream           = CmdCode(0x51)
	CmdDcpCloseStream         = CmdCode(0x52)
	CmdDcpStreamReq           = CmdCode(0x53)
	CmdDcpGetFailoverLog      = CmdCode(0x54)
	CmdDcpStreamEnd           = CmdCode(0x55)
	CmdDcpSnapshotMarker      = CmdCode(0x56)
	CmdDcpMutation            = CmdCode(0x57)
	CmdDcpDeletion            = CmdCode(0x58)
	CmdDcpExpiration          = CmdCode(0x59)
	CmdDcpSeqNoAdvanced       = CmdCode(0x64)
	CmdDcpOsoSnapshot         = CmdCode(0x65)
	CmdDcpFlush               = CmdCode(0x5a)
	CmdDcpSetVbucketState     = CmdCode(0x5b)
	CmdDcpNoop                = CmdCode(0x5c)
	CmdDcpBufferAck           = CmdCode(0x5d)
	CmdDcpControl             = CmdCode(0x5e)
	CmdDcpEvent               = CmdCode(0x5f)
	CmdGetReplica             = CmdCode(0x83)
	CmdSelectBucket           = CmdCode(0x89)
	CmdObserveSeqNo           = CmdCode(0x91)
	CmdObserve                = CmdCode(0x92)
	CmdGetLocked              = CmdCode(0x94)
	CmdUnlockKey              = CmdCode(0x95)
	CmdGetMeta                = CmdCode(0xa0)
	CmdSetMeta                = CmdCode(0xa2)
	CmdDelMeta                = CmdCode(0xa8)
	CmdGetClusterConfig       = CmdCode(0xb5)
	CmdGetRandom              = CmdCode(0xb6)
	CmdCollectionsGetManifest = CmdCode(0xba)
	CmdCollectionsGetID       = CmdCode(0xbb)
	CmdSubDocGet              = CmdCode(0xc5)
	CmdSubDocExists           = CmdCode(0xc6)
	CmdSubDocDictAdd          = CmdCode(0xc7)
	CmdSubDocDictSet          = CmdCode(0xc8)
	CmdSubDocDelete           = CmdCode(0xc9)
	CmdSubDocReplace          = CmdCode(0xca)
	CmdSubDocArrayPushLast    = CmdCode(0xcb)
	CmdSubDocArrayPushFirst   = CmdCode(0xcc)
	CmdSubDocArrayInsert      = CmdCode(0xcd)
	CmdSubDocArrayAddUnique   = CmdCode(0xce)
	CmdSubDocCounter          = CmdCode(0xcf)
	CmdSubDocMultiLookup      = CmdCode(0xd0)
	CmdSubDocMultiMutation    = CmdCode(0xd1)
	CmdSubDocGetCount         = CmdCode(0xd2)
	CmdGetErrorMap            = CmdCode(0xfe)
)

// Name returns the string representation of the CmdCode.
func (command CmdCode) Name() string {
	switch command {
	case CmdGet:
		return "CMD_GET"
	case CmdSet:
		return "CMD_SET"
	case CmdAdd:
		return "CMD_ADD"
	case CmdReplace:
		return "CMD_REPLACE"
	case CmdDelete:
		return "CMD_DELETE"
	case CmdIncrement:
		return "CMD_INCREMENT"
	case CmdDecrement:
		return "CMD_DECREMENT"
	case CmdNoop:
		return "CMD_NOOP"
	case CmdAppend:
		return "CMD_APPEND"
	case CmdPrepend:
		return "CMD_PREPEND"
	case CmdStat:
		return "CMD_STAT"
	case CmdTouch:
		return "CMD_TOUCH"
	case CmdGAT:
		return "CMD_GAT"
	case CmdHello:
		return "CMD_HELLO"
	case CmdSASLListMechs:
		return "CMD_SASLLISTMECHS"
	case CmdSASLAuth:
		return "CMD_SASLAUTH"
	case CmdSASLStep:
		return "CMD_SASLSTEP"
	case CmdGetAllVBSeqnos:
		return "CMD_GETALLVBSEQNOS"
	case CmdDcpOpenConnection:
		return "CMD_DCPOPENCONNECTION"
	case CmdDcpAddStream:
		return "CMD_DCPADDSTREAM"
	case CmdDcpCloseStream:
		return "CMD_DCPCLOSESTREAM"
	case CmdDcpStreamReq:
		return "CMD_DCPSTREAMREQ"
	case CmdDcpGetFailoverLog:
		return "CMD_DCPGETFAILOVERLOG"
	case CmdDcpStreamEnd:
		return "CMD_DCPSTREAMEND"
	case CmdDcpSnapshotMarker:
		return "CMD_DCPSNAPSHOTMARKER"
	case CmdDcpMutation:
		return "CMD_DCPMUTATION"
	case CmdDcpDeletion:
		return "CMD_DCPDELETION"
	case CmdDcpExpiration:
		return "CMD_DCPEXPIRATION"
	case CmdDcpFlush:
		return "CMD_DCPFLUSH"
	case CmdDcpSetVbucketState:
		return "CMD_DCPSETVBUCKETSTATE"
	case CmdDcpNoop:
		return "CMD_DCPNOOP"
	case CmdDcpBufferAck:
		return "CMD_DCPBUFFERACK"
	case CmdDcpControl:
		return "CMD_DCPCONTROL"
	case CmdGetReplica:
		return "CMD_GETREPLICA"
	case CmdSelectBucket:
		return "CMD_SELECTBUCKET"
	case CmdObserveSeqNo:
		return "CMD_OBSERVESEQNO"
	case CmdObserve:
		return "CMD_OBSERVE"
	case CmdGetLocked:
		return "CMD_GETLOCKED"
	case CmdUnlockKey:
		return "CMD_UNLOCKKEY"
	case CmdGetMeta:
		return "CMD_GETMETA"
	case CmdSetMeta:
		return "CMD_SETMETA"
	case CmdDelMeta:
		return "CMD_DELMETA"
	case CmdGetClusterConfig:
		return "CMD_GETCLUSTERCONFIG"
	case CmdGetRandom:
		return "CMD_GETRANDOM"
	case CmdSubDocGet:
		return "CMD_SUBDOCGET"
	case CmdSubDocExists:
		return "CMD_SUBDOCEXISTS"
	case CmdSubDocDictAdd:
		return "CMD_SUBDOCDICTADD"
	case CmdSubDocDictSet:
		return "CMD_SUBDOCDICTSET"
	case CmdSubDocDelete:
		return "CMD_SUBDOCDELETE"
	case CmdSubDocReplace:
		return "CMD_SUBDOCREPLACE"
	case CmdSubDocArrayPushLast:
		return "CMD_SUBDOCARRAYPUSHLAST"
	case CmdSubDocArrayPushFirst:
		return "CMD_SUBDOCARRAYPUSHFIRST"
	case CmdSubDocArrayInsert:
		return "CMD_SUBDOCARRAYINSERT"
	case CmdSubDocArrayAddUnique:
		return "CMD_SUBDOCARRAYADDUNIQUE"
	case CmdSubDocCounter:
		return "CMD_SUBDOCCOUNTER"
	case CmdSubDocMultiLookup:
		return "CMD_SUBDOCMULTILOOKUP"
	case CmdSubDocMultiMutation:
		return "CMD_SUBDOCMULTIMUTATION"
	case CmdSubDocGetCount:
		return "CMD_SUBDOCGETCOUNT"
	case CmdGetErrorMap:
		return "CMD_GETERRORMAP"
	default:
		return "CMD_x" + hex.EncodeToString([]byte{byte(command)})
	}
}
