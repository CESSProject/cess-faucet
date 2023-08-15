package chain

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

type Event_PPBNoOnTimeSubmit struct {
	Phase     types.Phase
	Acc       types.AccountID
	SegmentId types.U64
	Topics    []types.Hash
}

type Event_PPDNoOnTimeSubmit struct {
	Phase     types.Phase
	Acc       types.AccountID
	SegmentId types.U64
	Topics    []types.Hash
}

type Event_ChallengeProof struct {
	Phase  types.Phase
	PeerId types.U64
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_VerifyProof struct {
	Phase  types.Phase
	PeerId types.U64
	Fileid types.Bytes
	Topics []types.Hash
}

// ------------------------Sminer---------------------------------
type Event_Registered struct {
	Phase      types.Phase
	Acc        types.AccountID
	StakingVal types.U128
	Topics     []types.Hash
}

type Event_TimedTask struct {
	Phase  types.Phase
	Topics []types.Hash
}

type Event_DrawFaucetMoney struct {
	Phase  types.Phase
	Topics []types.Hash
}

type Event_FaucetTopUpMoney struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}

type Event_LessThan24Hours struct {
	Phase  types.Phase
	Last   types.U32
	Now    types.U32
	Topics []types.Hash
}
type Event_AlreadyFrozen struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}

type Event_MinerExit struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}

type Event_MinerClaim struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}

type Event_IncreaseCollateral struct {
	Phase   types.Phase
	Acc     types.AccountID
	Balance types.U128
	Topics  []types.Hash
}

type Event_Deposit struct {
	Phase   types.Phase
	Balance types.U128
	Topics  []types.Hash
}

// ------------------------FileBank-------------------------------
type Event_DeleteFile struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_BuySpace struct {
	Phase  types.Phase
	Acc    types.AccountID
	Size   types.U128
	Fee    types.U128
	Topics []types.Hash
}

type Event_FileUpload struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}

type Event_FileUpdate struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_LeaseExpireIn24Hours struct {
	Phase  types.Phase
	Acc    types.AccountID
	Size   types.U128
	Topics []types.Hash
}

type Event_FileChangeState struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_BuyFile struct {
	Phase  types.Phase
	Acc    types.AccountID
	Money  types.U128
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_Purchased struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_InsertFileSlice struct {
	Phase  types.Phase
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_LeaseExpired struct {
	Phase  types.Phase
	Acc    types.AccountID
	Size   types.U128
	Topics []types.Hash
}

type Event_FillerUpload struct {
	Phase    types.Phase
	Acc      types.AccountID
	Filesize types.U64
	Topics   []types.Hash
}

type Event_ClearInvalidFile struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

type Event_RecoverFile struct {
	Phase  types.Phase
	Acc    types.AccountID
	Fileid types.Bytes
	Topics []types.Hash
}

// ------------------------FileMap--------------------------------
type Event_RegistrationScheduler struct {
	Phase  types.Phase
	Acc    types.AccountID
	Ip     types.Bytes
	Topics []types.Hash
}

// ------------------------other system---------------------------
type Event_UnsignedPhaseStarted struct {
	Phase  types.Phase
	Round  types.U32
	Topics []types.Hash
}

type Event_SignedPhaseStarted struct {
	Phase  types.Phase
	Round  types.U32
	Topics []types.Hash
}

type Event_SolutionStored struct {
	Phase            types.Phase
	Election_compute types.ElectionCompute
	Prev_ejected     types.Bool
	Topics           []types.Hash
}

type Event_Balances_Withdraw struct {
	Phase  types.Phase
	Who    types.AccountID
	Amount types.U128
	Topics []types.Hash
}
type Event_OutstandingChallenges struct {
	Phase  types.Phase
	PeerId types.U64
	Fileid types.Bytes
	Topics []types.Hash
}
type Event_ReceiveSpace struct {
	Phase  types.Phase
	Acc    types.AccountID
	Topics []types.Hash
}
