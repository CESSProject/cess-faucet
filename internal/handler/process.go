package handler

import (
	"cess-faucet/config"
	"cess-faucet/internal/chain"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strings"
	"sync"
	"time"
)

type TransRq struct {
	Address string `json:"Address"`
}
type TransRs struct {
	Err       string `json:"Err"`
	AsInBlock bool   `json:"AsInBlock"`
}

type IpLimit struct {
	MapLock sync.Mutex
	IpMap   map[string]map[string]time.Time //[ip,account list]
}

var IpLimitMap = &IpLimit{
	IpMap: make(map[string]map[string]time.Time),
}

type RequestLimit struct {
	RequestMap sync.Map //[account,time]
}

var RequestLimitMap = &RequestLimit{
	RequestMap: sync.Map{},
}

func ReplySuccess(ctx *gin.Context, r interface{}) {
	ctx.JSON(200, gin.H{
		"Result": r,
		"Status": "Success",
	})
}

func ReplyFail(ctx *gin.Context, r interface{}) {
	ctx.JSON(400, gin.H{
		"Result": r,
		"Status": "Fail",
	})
}

func Transfer(ctx *gin.Context) {
	var err error
	var query TransRq
	var result TransRs
	if err = ctx.ShouldBindBodyWith(&query, binding.JSON); err != nil {
		result.Err = fmt.Sprintf("%s", err)
		result.AsInBlock = false
		ReplyFail(ctx, result)
		return
	}
	userAgent := ctx.GetHeader("User-Agent")
	if strings.Contains(userAgent, "Go-http-client") {
		result.Err = "Portal faucet temporarily disable"
		result.AsInBlock = false
		ReplyFail(ctx, result)
		return
	}
	fmt.Println("ClientIP:", ctx.ClientIP(), "User-Agent:", userAgent)
	IpLimitMap.MapLock.Lock()
	defer IpLimitMap.MapLock.Unlock()
	AccountSlice, ok := IpLimitMap.IpMap[ctx.ClientIP()]
	if ok {
		if len(AccountSlice) >= config.IpLimitAccountNum {
			var ResponseFail = true
			for i, v := range AccountSlice {
				if time.Now().Sub(v) > config.AccountExistTime {
					delete(AccountSlice, i)
					AccountSlice[query.Address] = time.Now()
					ResponseFail = false
					break
				}
			}
			if ResponseFail {
				result.Err = "Too many requests from this ip, please try again 24H later"
				result.AsInBlock = false
				ReplyFail(ctx, result)
				return
			}
		} else {
			AccountSlice[query.Address] = time.Now()
		}
	} else {
		IpLimitMap.IpMap[ctx.ClientIP()] = make(map[string]time.Time, config.IpLimitAccountNum)
		IpLimitMap.IpMap[ctx.ClientIP()][query.Address] = time.Now()
	}

	lastRequest, ok := RequestLimitMap.RequestMap.Load(query.Address)
	if ok {
		last := lastRequest.(time.Time)
		if time.Now().Sub(last) < config.AccountExistTime {
			result.Err = "Too many claims, please come back after 24 hours"
			result.AsInBlock = false
			ReplyFail(ctx, result)
			return
		}
	}

	trans, accountId, money, err := WithDirectTransfer(query.Address, 100000000000000)
	if err != nil {
		result.Err = fmt.Sprintf("%s", err)
		result.AsInBlock = false
		ReplyFail(ctx, result)
		return
	}

	result.AsInBlock, err = trans.TradeOnChainByDirectTransfer(accountId, money)
	if err != nil {
		result.Err = fmt.Sprintf("%s", err)
		result.AsInBlock = false
		ReplyFail(ctx, result)
		return
	}
	RequestLimitMap.RequestMap.Store(query.Address, time.Now())
	result.Err = ""
	ReplySuccess(ctx, result)
	return
}

func WithSminerFaucet(Address string) (chain.CessInfo, *types.AccountID, error) {
	var trans chain.CessInfo
	trans.RpcAddr = config.Data.CoreData.CessRpcAddr
	trans.IdentifyAccountPhrase = config.Data.CoreData.IdAccountPhraseOrSeed
	trans.TransactionName = config.ChainTransaction_Faucet

	addr, err := chain.ParsingPublickey(Address)
	if err != nil {
		return chain.CessInfo{}, nil, err
	}
	AccountID, err := types.NewAccountID(addr)
	if err != nil {
		return chain.CessInfo{}, nil, err
	}
	return trans, AccountID, nil
}

func WithDirectTransfer(Address string, money uint64) (chain.CessInfo, types.MultiAddress, types.UCompact, error) {
	var trans chain.CessInfo
	trans.RpcAddr = config.Data.CoreData.CessRpcAddr
	trans.IdentifyAccountPhrase = config.Data.CoreData.IdAccountPhraseOrSeed
	trans.TransactionName = config.ChainTransaction_Transfer

	addr, err := chain.ParsingPublickey(Address)
	if err != nil {
		return chain.CessInfo{}, types.MultiAddress{}, types.UCompact{}, err
	}

	AccountID, err := types.NewMultiAddressFromAccountID(addr)
	if err != nil {
		return chain.CessInfo{}, types.MultiAddress{}, types.UCompact{}, err
	}
	types.NewUCompactFromUInt(money)

	return trans, AccountID, types.NewUCompactFromUInt(money), nil
}
