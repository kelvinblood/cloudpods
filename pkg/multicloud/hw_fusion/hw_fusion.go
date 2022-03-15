package hw_fusion

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/onecloud/pkg/cloudprovider"
	"yunion.io/x/onecloud/pkg/util/httputils"
	"yunion.io/x/pkg/errors"
)

type HwFusionClientConfig struct {
	cpcfg    cloudprovider.ProviderConfig
	endpoint string
	user     string
	password string
	debug    bool
}

func NewHwFusionClientConfig(endpoint, user, password string) *HwFusionClientConfig {
	cfg := &HwFusionClientConfig{
		endpoint: endpoint,
		user:     user,
		password: password,
	}
	return cfg
}

func (cfg *HwFusionClientConfig) CloudproviderConfig(cpcfg cloudprovider.ProviderConfig) *HwFusionClientConfig {
	cfg.cpcfg = cpcfg
	return cfg
}

func (cfg *HwFusionClientConfig) Debug(debug bool) *HwFusionClientConfig {
	cfg.debug = debug
	return cfg
}

type SHwFusionComputeClient struct {
	*HwFusionClientConfig
	authToken string

	regions            []SRegion
	defaultRegionIndex int
}

func NewHwFusionComputeClient(cfg *HwFusionClientConfig) (*SHwFusionComputeClient, error) {
	client := &SHwFusionComputeClient{
		HwFusionClientConfig: cfg,
	}
	// 登录
	err := client.login()
	if err != nil {
		return nil, err
	}
	// 站点列表
	client.regions, err = client.GetRegions()
	if err != nil {
		return nil, err
	}
	for i := range client.regions {
		client.regions[i].client = client
	}

	return client, nil
}

func (self *SHwFusionComputeClient) GetAuthToken() string {
	return self.authToken
}

func (self *SHwFusionComputeClient) GetRegion(id string) (*SRegion, error) {
	for i := range self.regions {
		if self.regions[i].Urn == id {
			self.defaultRegionIndex = i
			return &self.regions[i], nil
		}
	}
	if len(id) == 0 {
		return &self.regions[0], nil
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SHwFusionComputeClient) getDefaultClient(timeout time.Duration) *http.Client {
	client := httputils.GetDefaultClient()
	if timeout > 0 {
		client = httputils.GetTimeoutClient(timeout)
	}
	proxy := func(req *http.Request) (*url.URL, error) {
		if self.cpcfg.ProxyFunc != nil {
			self.cpcfg.ProxyFunc(req)
		}
		return nil, nil
	}

	httputils.SetClientProxyFunc(client, proxy)
	return client
}

func (self *SHwFusionComputeClient) invoke(method, uri string, params map[string]string) (jsonutils.JSONObject, error) {
	req, _ := http.NewRequest(method, self.HwFusionClientConfig.endpoint+uri, nil)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json;version=6.0;charset=UTF-8")
	req.Header.Set("X-Auth-Token", self.authToken)
	cli := self.getDefaultClient(0)
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if self.debug {
		log.Debugf("response: %s", string(result))
	}

	obj, err := jsonutils.Parse(result)
	if err != nil {
		return nil, errors.Wrapf(err, "jsonutils.Parse")
	}

	respKey := uri + "Response"
	if obj.Contains(respKey) {
		obj, err = obj.Get(respKey)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (self *SHwFusionComputeClient) login() error {
	req, _ := http.NewRequest("POST", self.HwFusionClientConfig.endpoint+LOGIN_URI, nil)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json;version=6.0;charset=UTF-8")
	req.Header.Set("X-Auth-User", self.HwFusionClientConfig.user)
	req.Header.Set("X-Auth-Key", self.HwFusionClientConfig.password)
	req.Header.Set("X-ENCRIPT-ALGORITHM", "1")
	req.Header.Set("X-Auth-UserType", "0")

	cli := self.getDefaultClient(0)
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	self.authToken = resp.Header.Get("X-Auth-Token")
	return nil
}
