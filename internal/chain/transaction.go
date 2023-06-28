package chain

import (
	"cess-faucet/logger"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"time"
)

type CessInfo struct {
	RpcAddr               string
	IdentifyAccountPhrase string
	TransactionName       string
	ChainModule           string
	ChainModuleMethod     string
}

var (
	CessPrefix      = []byte{0x50, 0xac}
	SubstratePrefix = []byte{0x2a}
	SSPrefix        = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
)

// etcd register
func (ci *CessInfo) TradeOnChain(Addr string) (bool, error) {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api := getSubstrateApi_safe()
	defer func() {
		releaseSubstrateApi()
		recover()
	}()
	keyring, err := signature.KeyringPairFromSecret(ci.IdentifyAccountPhrase, 0)
	if err != nil {
		return false, errors.Wrap(err, "KeyringPairFromSecret err")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return false, errors.Wrap(err, "GetMetadataLatest err")
	}
	addr, err := ParsingPublickey(Addr)
	if err != nil {
		return false, err
	}
	AccountID, err := types.NewAccountID(addr)
	if err != nil {
		return false, err
	}
	c, err := types.NewCall(meta, ci.TransactionName, AccountID)
	if err != nil {
		return false, errors.Wrap(err, "NewCall err")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return false, errors.Wrap(err, "NewExtrinsic err")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return false, errors.Wrap(err, "GetBlockHash err")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return false, errors.Wrap(err, "GetRuntimeVersionLatest err")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return false, errors.Wrap(err, "CreateStorageKey err")
	}
	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return false, errors.Wrap(err, "CreateStorageKey System Events err")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return false, errors.Wrap(err, "GetStorageLatest err")
	}
	if !ok {
		return false, errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return false, errors.Wrap(err, "Sign err")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return false, errors.Wrap(err, "SubmitAndWatchExtrinsic err")
	}
	defer sub.Unsubscribe()

	timeout := time.After(10 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				logger.InfoLogger.Sugar().Infof("[%v] tx blockhash: %#x", ci.TransactionName, status.AsInBlock)
				events := MyEventRecords{}
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return false, err
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					fmt.Println("Analyze event err: ", err)
				}

				if len(events.Sminer_DrawFaucetMoney) != 0 {
					return true, nil
				}
				return false, errors.New("Please wait for 24 hours to claim!")
			}
		case <-timeout:
			return false, errors.Errorf("[%v] tx timeout", ci.TransactionName)
		default:
			time.Sleep(time.Second)
		}
	}
}

func ParsingPublickey(address string) ([]byte, error) {
	err := VerityAddress(address, CessPrefix)
	if err != nil {
		err := VerityAddress(address, SubstratePrefix)
		if err != nil {
			return nil, errors.New("Invalid account")
		}
		data := base58.Decode(address)
		if len(data) != (34 + len(SubstratePrefix)) {
			return nil, errors.New("Public key decoding failed")
		}
		return data[len(SubstratePrefix) : len(data)-2], nil
	} else {
		data := base58.Decode(address)
		if len(data) != (34 + len(CessPrefix)) {
			return nil, errors.New("Public key decoding failed")
		}
		return data[len(CessPrefix) : len(data)-2], nil
	}
}

func VerityAddress(address string, prefix []byte) error {
	decodeBytes := base58.Decode(address)
	if len(decodeBytes) != (34 + len(prefix)) {
		return errors.New("Public key decoding failed")
	}
	if decodeBytes[0] != prefix[0] {
		return errors.New("Invalid account prefix")
	}
	pub := decodeBytes[len(prefix) : len(decodeBytes)-2]

	data := append(prefix, pub...)
	input := append(SSPrefix, data...)
	ck := blake2b.Sum512(input)
	checkSum := ck[:2]
	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[32+len(prefix)+i] {
			return errors.New("Invalid account")
		}
	}
	if len(pub) != 32 {
		return errors.New("Invalid account public key")
	}
	return nil
}
