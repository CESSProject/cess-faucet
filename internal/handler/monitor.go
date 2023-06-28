package handler

import (
	"cess-faucet/config"
	"time"
)

func (IpMap *IpLimit) IpMonitor() {
	for {
		select {
		case MonitorAccount := <-IpMap.MonitorQueue:
			IpMap.MapLock.Lock()
			if value, ok := IpMap.IpMap[MonitorAccount[0]]; ok {
				RegisterTime := value[MonitorAccount[1]]
				IpMap.MapLock.Unlock()
				for time.Now().Sub(RegisterTime) < config.AccountExistTime {
					time.Sleep(time.Second * 3)
				}
				IpMap.MapLock.Lock()
				accountmap, _ := IpMap.IpMap[MonitorAccount[0]]
				delete(accountmap, MonitorAccount[1])
				IpMap.MapLock.Unlock()
			} else {
				IpMap.MapLock.Unlock()
				panic("Elements of IpMap are missing")
			}
		default:
			continue
		}
	}

}
