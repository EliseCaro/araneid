package service

import (
	"fmt"
	_func "github.com/beatrice950201/araneid/extend/func"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"strconv"
	"time"
)

type DefaultAdminService struct{}

/** 数据总包**/
type Dashboard struct {
	Network *net.IOCountersStat
	Memory  *mem.VirtualMemoryStat
	CPU     *CpuDashboard
	Load    *load.AvgStat
}

/** cpu 返回数据 **/
type CpuDashboard struct {
	UsedPercent float64 `json:"used_percent"` // cpu使用百分比
	ModelName   string  `json:"modelName"`    // CPU名字
	Cores       int32   `json:"cores"`        // CPU 总核心数
}

/** 磁盘数据 **/
type DiskStatus struct {
	Name        string  `json:"name"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	Total       uint64  `json:"total"`
	FsType      string  `json:"fs_type"`
	Detail      string  `json:"detail"`
	UsedPercent float64 `json:"used_percent"` // 使用百分比
}

/** 浮点保留两位小数 **/
func (service *DefaultAdminService) decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

/** 组装详情数据 **/
func (service *DefaultAdminService) diskDetail(item *DiskStatus) string {
	html := fmt.Sprintf(`挂载点 : %s <br>`, item.Name)
	html += fmt.Sprintf(`磁盘类型 : %s <br>`, item.FsType)
	html += fmt.Sprintf(`已用容量 : %d GB <br>`, int(item.Used))
	html += fmt.Sprintf(`可用容量 : %d GB <br>`, int(item.Free))
	html += fmt.Sprintf(`总共容量 : %d GB <br>`, int(item.Total))
	return html
}

/** 网络监听控制面板 **/
func (service *DefaultAdminService) networkDashboard() *net.IOCountersStat {
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

/** 内存监控 **/
func (service *DefaultAdminService) memoryDashboard() *mem.VirtualMemoryStat {
	v, _ := mem.VirtualMemory()
	return v
}

/** CPU监控 **/
func (service *DefaultAdminService) cpuDashboard() *CpuDashboard {
	info, _ := cpu.Info()
	percent, _ := cpu.Percent(time.Second, false)
	var res CpuDashboard
	for _, v := range info {
		res.ModelName += v.ModelName
		res.Cores += v.Cores
	}
	for _, v := range percent {
		res.UsedPercent += v
	}
	return &res
}

/** CPU负载监控 **/
func (service *DefaultAdminService) loadDashboard() *load.AvgStat {
	info, _ := load.Avg()
	return info
}

/** 硬盘监控;从控制器直接提交 **/
func (service *DefaultAdminService) DiskDashboard() []*DiskStatus {
	var result []*DiskStatus
	parts, _ := disk.Partitions(true)
	for _, v := range parts {
		diskInfo, _ := disk.Usage(v.Mountpoint)
		item := &DiskStatus{
			Name:        diskInfo.Path,
			FsType:      diskInfo.Fstype,
			Total:       diskInfo.Total / 1024 / 1024 / 1024,
			Free:        diskInfo.Free / 1024 / 1024 / 1024,
			Used:        diskInfo.Used / 1024 / 1024 / 1024,
			UsedPercent: service.decimal(diskInfo.UsedPercent),
		}
		if item.Total > 0 {
			item.Detail = service.diskDetail(item)
			result = append(result, item)
		}
	}
	return result
}

/** 返回数据结构 **/
func (service *DefaultAdminService) DashboardInitialized() *Dashboard {
	return &Dashboard{
		Network: service.networkDashboard(),
		Memory:  service.memoryDashboard(),
		CPU:     service.cpuDashboard(),
		Load:    service.loadDashboard(),
	}
}

/** 获取监控面板 远行程序数量 **/
func (service *DefaultAdminService) DashboardProcessing() []map[string]interface{} {
	var result []map[string]interface{}
	result = append(result, map[string]interface{}{"title": "远行蜘蛛池", "count": new(DefaultArachnidService).aliveNum()})
	result = append(result, map[string]interface{}{"title": "远行采集器", "count": new(DefaultCollectService).aliveNum()})
	result = append(result, map[string]interface{}{"title": "远行发布器", "count": new(DefaultCollectService).alivePushNum()})
	result = append(result, map[string]interface{}{"title": "云盘资料数", "count": new(DefaultAdjunctService).aliveNum()})
	return result
}
