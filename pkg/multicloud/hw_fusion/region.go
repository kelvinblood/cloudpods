package hw_fusion

import (
	"yunion.io/x/onecloud/pkg/multicloud"
)

// SRegion 结构体对应Fusion Compute的site 站点概念
type SRegion struct {
	multicloud.SRegion
	client *SHwFusionComputeClient

	Urn    string `json:"urn"` // site 唯一标识
	Name   string `json:"name"`
	URI    string `json:"uri"`
	IP     string `json:"ip"`
	IsDC   bool   `json:"isDC"`
	IsSelf bool   `json:"isSelf"`
	Status string `json:"status"`

	// 以下 region-show 项
	TimeZone    string      `json:"timeZone,omitempty"`
	Description interface{} `json:"description,omitempty"`
	NtpBridgeIP string      `json:"ntpBridgeIp,omitempty"`
	NtpIP       string      `json:"ntpIp,omitempty"`
	NtpCycle    int         `json:"ntpCycle,omitempty"`
	NtpIP2      string      `json:"ntpIp2,omitempty"`
	NtpIP3      interface{} `json:"ntpIp3,omitempty"`
}

func (this *SRegion) GetClient() *SHwFusionComputeClient {
	return this.client
}

// GetRegions region-list
func (self *SHwFusionComputeClient) GetRegions() ([]SRegion, error) {
	resp, err := self.invoke("GET", SITES_URI, nil)
	if err != nil {
		return nil, err
	}
	result := struct {
		Sites []SRegion `json:"sites"`
	}{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return result.Sites, nil
}

// GetRegionDetail region-show
func (self *SHwFusionComputeClient) GetRegionDetail(regionId string) (*SRegion, error) {
	var uri string
	for _, region := range self.regions {
		if region.Urn == regionId {
			uri = region.URI
		}
	}
	if len(uri) == 0 {
		uri = self.regions[self.defaultRegionIndex].URI
	}
	resp, err := self.invoke("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	result := SRegion{}
	err = resp.Unmarshal(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
