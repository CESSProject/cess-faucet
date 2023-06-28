package config

import "time"

// system exit code
const (
	Exit_Normal                   = 0
	Exit_LoginFailed              = -1
	Exit_RunningSystemError       = -2
	Exit_ExecutionPermissionError = -3
	Exit_InvalidIP                = -4
	Exit_CreateFolder             = -5
	Exit_CreateEmptyFile          = -6
	Exit_ConfFileNotExist         = -7
	Exit_ConfFileFormatError      = -8
	Exit_ConfFileTypeError        = -9
	Exit_CmdLineParaErr           = -10
)

var (
	LogfilePathPrefix         = "./log/"
	ChainModule               = "Sminer"
	ChainTransaction_Register = "Sminer.faucet"
	ChainModule_Sminer_Search = ""
	IpLimitAccountNum         = 12
	AccountExistTime          = 24 * 60 * 60 * time.Second
)
