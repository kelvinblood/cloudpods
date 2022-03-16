package hw_fusion

import "strings"

type SInstance struct {
	ObjectPrivs             []interface{} `json:"objectPrivs"`
	DrStatus                int           `json:"drStatus"`
	InitSyncStatus          int           `json:"initSyncStatus"`
	MinCompatibleimcSetting string        `json:"minCompatibleimcSetting"`
	RpoStatus               int           `json:"rpoStatus"`
	VMConfig                struct {
		CPU    interface{} `json:"cpu"`
		Disks  interface{} `json:"disks"`
		Gpu    interface{} `json:"gpu"`
		Memory interface{} `json:"memory"`
		Nics   []struct {
			IPList        string      `json:"ipList"`
			Ips6          []string    `json:"ips6"`
			NicType       interface{} `json:"nicType"`
			SequenceNum   int         `json:"sequenceNum"`
			VirtIo        interface{} `json:"virtIo"`
			IP            string      `json:"ip"`
			PortGroupUrn  interface{} `json:"portGroupUrn"`
			PortGroupName interface{} `json:"portGroupName"`
			Mac           string      `json:"mac"`
			Urn           interface{} `json:"urn"`
			URI           interface{} `json:"uri"`
			Name          interface{} `json:"name"`
		} `json:"nics"`
		Usb        interface{} `json:"usb"`
		Properties struct {
			GpuShareType      interface{} `json:"gpuShareType"`
			IsHpet            interface{} `json:"isHpet"`
			IsReserveResource interface{} `json:"isReserveResource"`
			SecureVMType      interface{} `json:"secureVmType"`
			AttachType        interface{} `json:"attachType"`
			BootOption        interface{} `json:"bootOption"`
			ClockMode         interface{} `json:"clockMode"`
			IsAutoUpgrade     interface{} `json:"isAutoUpgrade"`
			IsEnableFt        interface{} `json:"isEnableFt"`
			IsEnableHa        interface{} `json:"isEnableHa"`
			IsEnableMemVol    interface{} `json:"isEnableMemVol"`
			ReoverByHost      interface{} `json:"reoverByHost"`
			VMFaultProcess    interface{} `json:"vmFaultProcess"`
		} `json:"properties"`
	} `json:"vmConfig"`
	VMType            int         `json:"vmType"`
	Description       interface{} `json:"description"`
	HostUrn           string      `json:"hostUrn"`
	ClusterName       string      `json:"clusterName"`
	UUID              string      `json:"uuid"`
	ClusterUrn        string      `json:"clusterUrn"`
	ImcSetting        string      `json:"imcSetting"`
	HostName          string      `json:"hostName"`
	Urn               string      `json:"urn"`
	URI               string      `json:"uri"`
	CdRomStatus       string      `json:"cdRomStatus"`
	Idle              int         `json:"idle"`
	IsBindingHost     bool        `json:"isBindingHost"`
	IsLinkClone       bool        `json:"isLinkClone"`
	IsTemplate        bool        `json:"isTemplate"`
	LocationName      string      `json:"locationName"`
	PvDriverStatus    string      `json:"pvDriverStatus"`
	ToolInstallStatus string      `json:"toolInstallStatus"`
	ToolsVersion      string      `json:"toolsVersion"`
	CreateTime        string      `json:"createTime"`
	Group             interface{} `json:"group"`
	Location          string      `json:"location"`
	Status            string      `json:"status"`
	Name              string      `json:"name"`

	// 以下为 instance-show 项
	Params struct {
		ExternalUUID     string      `json:"externalUuid,omitempty"`
		VMSubStatus      string      `json:"vmSubStatus,omitempty"`
		CdromSequenceNum string      `json:"cdromSequenceNum,omitempty"`
		ParentObjUrn     interface{} `json:"parentObjUrn,omitempty"`
		PciIbCard        string      `json:"PCI_IB_CARD,omitempty"`
		Gpu              string      `json:"gpu,omitempty"`
	} `json:"params,omitempty"`
	AdditionalStatus   []interface{} `json:"additionalStatus,omitempty"`
	DataStoreUrns      []string      `json:"dataStoreUrns,omitempty"`
	DeleteTime         interface{}   `json:"deleteTime,omitempty"`
	DrDrillVMURI       interface{}   `json:"drDrillVmUri,omitempty"`
	DrDrillVMUrn       interface{}   `json:"drDrillVmUrn,omitempty"`
	IsMultiDiskSpeedup bool          `json:"isMultiDiskSpeedup,omitempty"`
	OsOptions          struct {
		Password    string      `json:"password,omitempty"`
		Hostname    string      `json:"hostname,omitempty"`
		OsType      string      `json:"osType,omitempty"`
		GuestOSName interface{} `json:"guestOSName,omitempty"`
		OsVersion   int         `json:"osVersion,omitempty"`
	} `json:"osOptions,omitempty"`
	VMRebootConfig struct {
		CPU struct {
			CoresPerSocket int `json:"coresPerSocket,omitempty"`
			Quantity       int `json:"quantity,omitempty"`
			Reservation    int `json:"reservation,omitempty"`
			Weight         int `json:"weight,omitempty"`
			Limit          int `json:"limit,omitempty"`
		} `json:"cpu,omitempty"`
		Memory struct {
			Reservation int `json:"reservation,omitempty"`
			Weight      int `json:"weight,omitempty"`
			QuantityMB  int `json:"quantityMB,omitempty"`
			Limit       int `json:"limit,omitempty"`
		} `json:"memory,omitempty"`
	} `json:"vmRebootConfig,omitempty"`
	VncAcessInfo struct {
		VncOldPassword interface{} `json:"vncOldPassword,omitempty"`
		HostIP         string      `json:"hostIp,omitempty"`
		VncPassword    string      `json:"vncPassword,omitempty"`
		VncPort        int         `json:"vncPort,omitempty"`
	} `json:"vncAcessInfo,omitempty"`
}

type AutoGenerated struct {
}

// GetInstances instance-list
func (self *SHwFusionComputeClient) GetInstances(regionId string) ([]SInstance, error) {
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	uri = strings.ReplaceAll(VMS_URI, PREFIX_SITE_URI, uri)
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := struct {
		Total int         `json:"total"`
		Vms   []SInstance `json:"vms"`
	}{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return result.Vms, nil
}

// GetInstanceDetail  instance-show
func (self *SHwFusionComputeClient) GetInstanceDetail(regionId, vmId string) (*SInstance, error) {
	// fixme 待完善返回的instance结构体
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	uri = strings.ReplaceAll(VMS_URI+`/`+vmId, PREFIX_SITE_URI, uri)
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := SInstance{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetInstanceStatistics instance-statistics
func (self *SHwFusionComputeClient) GetInstanceStatistics(regionId string) (interface{}, error) {
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	uri = strings.ReplaceAll(VMS_STATISTICS_URI, PREFIX_SITE_URI, uri)
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := struct {
		Running    int `json:"running"`
		Total      int `json:"total"`
		Stopped    int `json:"stopped"`
		Hibernated int `json:"hibernated"`
		Fault      int `json:"fault"`
	}{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
