package shell

import (
	"yunion.io/x/onecloud/pkg/multicloud/hw_fusion"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

func init() {
	type RegionListOptions struct {
	}
	shellutils.R(&RegionListOptions{}, "region-list", "List all regions(sites)", func(cli *hw_fusion.SRegion, args *RegionListOptions) error {
		regions, err := cli.GetClient().GetRegions()
		if err != nil {
			return err
		}
		printList(regions, 0, 0, 0, nil)
		return nil
	})

	type RegionIdOptions struct {
		RegionID string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
	}
	shellutils.R(&RegionIdOptions{}, "region-show", "Show the region(site) detailed information", func(cli *hw_fusion.SRegion, args *RegionIdOptions) error {
		regions, err := cli.GetClient().GetRegionDetail(args.RegionID)
		if err != nil {
			return err
		}
		printObject(regions)
		return nil
	})
}
