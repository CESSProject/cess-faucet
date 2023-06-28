package handler

import (
	"cess-faucet/config"
	"cess-faucet/internal/chain"
	"fmt"
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
	MapLock      sync.Mutex
	IpMap        map[string]map[string]time.Time
	MonitorQueue chan [2]string //[ip,account]
}

var IpLimitMap = &IpLimit{
	IpMap: make(map[string]map[string]time.Time),
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

	var trans chain.CessInfo
	trans.RpcAddr = config.Data.CoreData.CessRpcAddr
	trans.IdentifyAccountPhrase = config.Data.CoreData.IdAccountPhraseOrSeed
	trans.TransactionName = config.ChainTransaction_Register
	trans.ChainModule = config.ChainModule
	trans.ChainModuleMethod = config.ChainModule_Sminer_Search
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
				IpLimitMap.MapLock.Unlock()
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
	IpLimitMap.MapLock.Unlock()

	result.AsInBlock, err = trans.TradeOnChain(query.Address)
	if err != nil {
		result.Err = fmt.Sprintf("%s", err)
		result.AsInBlock = false
		ReplyFail(ctx, result)
		return
	}
	result.Err = ""
	ReplySuccess(ctx, result)
	return
}
