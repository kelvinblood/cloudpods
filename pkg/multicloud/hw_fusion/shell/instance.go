package shell

import (
	"yunion.io/x/onecloud/pkg/multicloud/hw_fusion"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type VmsListOptions struct {
		RegionID string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
	}
	shellutils.R(&VmsListOptions{}, "instance-list", "List all instances(vms)", func(cli *hw_fusion.SRegion, args *VmsListOptions) error {
		vms, err := cli.GetClient().GetInstances(args.RegionID)
		if err != nil {
			return err
		}
		printList(vms, 0, 0, 0, nil)
		return nil
	})

	type VmIdOptions struct {
		RegionID string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
		VmId     string `help:"VmId" default:"$HUAWEI_FUSION_COMPUTE_VMID" metavar:"HUAWEI_FUSION_COMPUTE_VMID"`
	}
	shellutils.R(&VmIdOptions{}, "instance-show", "Show the instance's detailed information", func(cli *hw_fusion.SRegion, args *VmIdOptions) error {
		vm, err := cli.GetClient().GetInstanceDetail(args.RegionID, args.VmId)
		if err != nil {
			return err
		}
		printObject(vm)
		return nil
	})

	type VmsStatisticsOptions struct {
		RegionID string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
	}
	shellutils.R(&VmsStatisticsOptions{}, "instance-statistics", "Show the instance's statistical information", func(cli *hw_fusion.SRegion, args *VmsStatisticsOptions) error {
		vm, err := cli.GetClient().GetInstanceStatistics(args.RegionID)
		if err != nil {
			return err
		}
		printObject(vm)
		return nil
	})
}
