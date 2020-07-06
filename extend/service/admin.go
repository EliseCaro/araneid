package service

import (
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/shirou/gopsutil/net"
)

type DefaultAdminService struct{}

/** 数据总包**/
type Dashboard struct {
	Network *net.IOCountersStat
}

/** 网络监听控制面板 **/
func (service *DefaultAdminService) NetworkDashboard() *net.IOCountersStat {
	result := &net.IOCountersStat{}
	if stat, err := net.IOCounters(false); err == nil {
		item := &net.IOCountersStat{}
		for _, v := range stat {
			item.BytesRecv += v.BytesRecv
			item.BytesSent += v.BytesSent
			item.PacketsSent += v.PacketsSent
			item.PacketsRecv += v.PacketsRecv
		}
		network := _func.GetCache("network_cache")
		if network != "" {
			cache := network.(*net.IOCountersStat)
			result = &net.IOCountersStat{
				BytesSent:   item.BytesSent - cache.BytesSent,
				BytesRecv:   item.BytesRecv - cache.BytesRecv,
				PacketsRecv: item.PacketsRecv - cache.PacketsRecv,
				PacketsSent: item.PacketsSent - cache.PacketsSent,
			}
		}
		_ = _func.SetCache("network_cache", item)
	}
	return result
}

/** 返回数据结构 **/
func (service *DefaultAdminService) DashboardInitialized() *Dashboard {
	return &Dashboard{
		Network: service.NetworkDashboard(),
	}
}
