package shell

import (
	"yunion.io/x/onecloud/pkg/multicloud/hw_fusion"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type HostsListOptions struct {
		RegionId string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
	}
	shellutils.R(&HostsListOptions{}, "host-list", "List all hosts", func(cli *hw_fusion.SRegion, args *HostsListOptions) error {
		hosts, err := cli.GetClient().GetHosts(args.RegionId)
		if err != nil {
			return err
		}
		printList(hosts, 0, 0, 0, nil)
		return nil
	})

	type HostIdOptions struct {
		RegionID string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
		HostId   string `help:"HostId" default:"$HUAWEI_FUSION_COMPUTE_HOSTID" metavar:"HUAWEI_FUSION_COMPUTE_HOSTID"`
	}
	shellutils.R(&HostIdOptions{}, "host-show", "Show the instance's detailed information", func(cli *hw_fusion.SRegion, args *HostIdOptions) error {
		host, err := cli.GetClient().GetHostDetail(args.RegionID, args.HostId)
		if err != nil {
			return err
		}
		printObject(host)
		return nil
	})

	type HostsStatisticsOptions struct {
		RegionID string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
	}
	shellutils.R(&HostsStatisticsOptions{}, "host-statistics", "Show the instance's statistical information", func(cli *hw_fusion.SRegion, args *HostsStatisticsOptions) error {
		host, err := cli.GetClient().GetHostsStatistics(args.RegionID)
		if err != nil {
			return err
		}
		printObject(host)
		return nil
	})
}
