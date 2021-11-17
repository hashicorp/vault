package memd

var cidSupportedOps = []CmdCode{
	CmdGet,
	CmdSet,
	CmdAdd,
	CmdReplace,
	CmdDelete,
	CmdIncrement,
	CmdDecrement,
	CmdAppend,
	CmdPrepend,
	CmdTouch,
	CmdGAT,
	CmdGetReplica,
	CmdGetLocked,
	CmdUnlockKey,
	CmdGetMeta,
	CmdSetMeta,
	CmdDelMeta,
	CmdSubDocGet,
	CmdSubDocExists,
	CmdSubDocDictAdd,
	CmdSubDocDictSet,
	CmdSubDocDelete,
	CmdSubDocReplace,
	CmdSubDocArrayPushLast,
	CmdSubDocArrayPushFirst,
	CmdSubDocArrayInsert,
	CmdSubDocArrayAddUnique,
	CmdSubDocCounter,
	CmdSubDocMultiLookup,
	CmdSubDocMultiMutation,
	CmdSubDocGetCount,
	CmdDcpMutation,
	CmdDcpExpiration,
	CmdDcpDeletion,
}

func makeCidSupportedTable() []bool {
	var cidTableLen uint32
	for _, cmd := range cidSupportedOps {
		if uint32(cmd) >= cidTableLen {
			cidTableLen = uint32(cmd) + 1
		}
	}
	cidTable := make([]bool, cidTableLen)
	for _, cmd := range cidSupportedOps {
		cidTable[cmd] = true
	}
	return cidTable
}

var cidSupportedTable = makeCidSupportedTable()

// IsCommandCollectionEncoded returns whether a particular command code
// should have its key collection encoded when collections support is
// enabled for a particular connection
func IsCommandCollectionEncoded(cmd CmdCode) bool {
	cmdIdx := int(cmd)
	if cmdIdx < 0 || cmdIdx >= len(cidSupportedTable) {
		return false
	}
	return cidSupportedTable[cmdIdx]
}
