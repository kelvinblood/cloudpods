package hw_fusion

import "strings"

type ParamsHost struct {
	NowHanaOptimizedStrategy  string `json:"nowHanaOptimizedStrategy"`
	NowEnableIORingAdaptation string `json:"nowEnableIORingAdaptation"`
	HanaOptimizedStrategy     string `json:"hanaOptimizedStrategy"`
	RealtimeUsedSizeMB        string `json:"realtimeUsedSizeMB"`
	EnableIORingAdaptation    string `json:"enableIORingAdaptation"`
}

type CPUResource struct {
	AllocatedSizeMHz int `json:"allocatedSizeMHz"`
	TotalSizeMHz     int `json:"totalSizeMHz"`
}

type MemResource struct {
	AllocatedSizeMB int `json:"allocatedSizeMB"`
	TotalSizeMB     int `json:"totalSizeMB"`
}

type SHost struct {
	Params                 ParamsHost    `json:"params"`
	HostRealName           string        `json:"hostRealName"`
	ClusterEnableIOTailor  bool          `json:"clusterEnableIOTailor"`
	HaRole                 interface{}   `json:"haRole"`
	HaState                interface{}   `json:"haState"`
	IsFailOverHost         bool          `json:"isFailOverHost"`
	Description            interface{}   `json:"description"`
	IP                     string        `json:"ip"`
	ClusterName            string        `json:"clusterName"`
	AttachedISOVMs         []interface{} `json:"attachedISOVMs"`
	BmcIP                  interface{}   `json:"bmcIp"`
	BmcUserName            interface{}   `json:"bmcUserName"`
	ClusterUrn             string        `json:"clusterUrn"`
	ComputeResourceStatics string        `json:"computeResourceStatics"`
	CPUMHz                 int           `json:"cpuMHz"`
	CPUQuantity            int           `json:"cpuQuantity"`
	CPUResource            CPUResource   `json:"cpuResource"`
	GdvmMemory             int           `json:"gdvmMemory"`
	GdvmMemoryReboot       int           `json:"gdvmMemoryReboot"`
	GpuCapacity            int           `json:"gpuCapacity"`
	GpuCapacityReboot      int           `json:"gpuCapacityReboot"`
	GsvmMemory             int           `json:"gsvmMemory"`
	GsvmMemoryReboot       int           `json:"gsvmMemoryReboot"`
	HostMultiPathMode      string        `json:"hostMultiPathMode"`
	ImcSetting             string        `json:"imcSetting"`
	IsMaintaining          bool          `json:"isMaintaining"`
	MaxImcSetting          string        `json:"maxImcSetting"`
	MemQuantityMB          int           `json:"memQuantityMB"`
	MemResource            MemResource   `json:"memResource"`
	MultiPathMode          string        `json:"multiPathMode"`
	NicQuantity            int           `json:"nicQuantity"`
	NtpCycle               int           `json:"ntpCycle"`
	NtpIP1                 string        `json:"ntpIp1"`
	NtpIP2                 string        `json:"ntpIp2"`
	NtpIP3                 interface{}   `json:"ntpIp3"`
	PhysicalCPUQuantity    int           `json:"physicalCpuQuantity"`
	Urn                    string        `json:"urn"`
	URI                    string        `json:"uri"`
	Status                 string        `json:"status"`
	Name                   string        `json:"name"`

	// 以下 host-show 项
	HostDNSCfg     interface{} `json:"hostDNSCfg,omitempty"`
	HostRoutetable interface{} `json:"hostRoutetable,omitempty"`
	Hypervisor     string      `json:"hypervisor,omitempty"`
	Vendor         string      `json:"vendor,omitempty"`
	Model          string      `json:"model,omitempty"`
	DefaultGateway interface{} `json:"defaultGateway,omitempty"`
}

// GetHosts host-list
func (self *SHwFusionComputeClient) GetHosts(regionId string) ([]SHost, error) {
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	uri = strings.ReplaceAll(HOSTS_URI, PREFIX_SITE_URI, uri)
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := struct {
		Total int     `json:"total"`
		Hosts []SHost `json:"hosts"`
	}{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return result.Hosts, nil
}

// GetHostDetail host-show
func (self *SHwFusionComputeClient) GetHostDetail(regionId, hostId string) (*SHost, error) {
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	uri = strings.ReplaceAll(HOSTS_URI+`/`+hostId, PREFIX_SITE_URI, uri)
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := SHost{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetHostsStatistics host-statistics
func (self *SHwFusionComputeClient) GetHostsStatistics(regionId string) (interface{}, error) {
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	uri = strings.ReplaceAll(HOSTS_STATISTICS_URI, PREFIX_SITE_URI, uri)
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := struct {
		Total       int `json:"total"`
		Rebooting   int `json:"rebooting"`
		Shutdowning int `json:"shutdowning"`
		Unknown     int `json:"unknown"`
		Booting     int `json:"booting"`
		Initial     int `json:"initial"`
		Normal      int `json:"normal"`
		Poweroff    int `json:"poweroff"`
		Fault       int `json:"fault"`
	}{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
