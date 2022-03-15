// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/http/httpproxy"

	"yunion.io/x/structarg"

	"yunion.io/x/onecloud/pkg/cloudprovider"                //
	"yunion.io/x/onecloud/pkg/multicloud/hw_fusion"         // hw_fusion
	_ "yunion.io/x/onecloud/pkg/multicloud/hw_fusion/shell" // hw_fusion
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

// X-Auth-User:ronly
// X-Auth-Key:Admin@huawei123$
// https://10.10.1.52:7443
// 建议第三方系统将每分钟接口调用次数控制在300次以内，以保证系统的整体性能和稳定。QPS控制在300q/s
// TODO 优先完成虚拟机部门的查询功能，查询虚拟机统计信息

// 华为FusionCompute 私有云
type BaseOptions struct {
	Help       bool   `help:"Show help" default:"false"`
	Debug      bool   `help:"Show debug" default:"false"`
	EndPoint   string `help:"Endpoint" default:"$HUAWEI_FUSION_COMPUTE_ENDPOINT" metavar:"$HUAWEI_FUSION_COMPUTE_ENDPOINT"`
	User       string `help:"User" default:"$HUAWEI_FUSION_COMPUTE_USER" metavar:"HUAWEI_FUSION_COMPUTE_USER"`
	Password   string `help:"Password" default:"$HUAWEI_FUSION_COMPUTE_PASSWORD" metavar:"HUAWEI_FUSION_COMPUTE_PASSWORD"`
	RegionId   string `help:"RegionId" default:"$HUAWEI_FUSION_COMPUTE_REGION" metavar:"HUAWEI_FUSION_COMPUTE_REGION"`
	SUBCOMMAND string `help:"hw_fusioncli subcommand" subcommand:"true"`
}

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	// 生成一个命令行工具解析器
	parse, e := structarg.NewArgumentParser(&BaseOptions{},
		"hw_fusioncli",
		"Command-line interface to Huawei Fusion Compute API.",
		`See "hw_fusioncli help COMMAND" for help on a specific command.`)

	if e != nil {
		return nil, e
	}

	subcmd := parse.GetSubcommand()
	if subcmd == nil {
		return nil, fmt.Errorf("No subcommand argument.")
	}

	// help 函数需要使用到subcmd的闭包
	type HelpOptions struct {
		SUBCOMMAND string `help:"sub-command name"`
	}
	shellutils.R(&HelpOptions{}, "help", "Show help of a subcommand", func(args *HelpOptions) error {
		helpstr, e := subcmd.SubHelpString(args.SUBCOMMAND)
		if e != nil {
			return e
		} else {
			fmt.Print(helpstr)
			return nil
		}
	})

	for _, v := range shellutils.CommandTable {
		_, e := subcmd.AddSubParser(v.Options, v.Command, v.Desc, v.Callback)
		if e != nil {
			return nil, e
		}
	}
	return parse, nil
}

func showErrorAndExit(e error) {
	fmt.Fprintf(os.Stderr, "%s", e)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func newClient(options *BaseOptions) (*hw_fusion.SRegion, error) {
	if len(options.EndPoint) == 0 {
		return nil, fmt.Errorf("Missing endpoint")
	}

	if len(options.User) == 0 {
		return nil, fmt.Errorf("Missing user")
	}

	if len(options.Password) == 0 {
		return nil, fmt.Errorf("Missing password")
	}

	cfg := &httpproxy.Config{
		HTTPProxy:  os.Getenv("HTTP_PROXY"),
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		NoProxy:    os.Getenv("NO_PROXY"),
	}
	cfgProxyFunc := cfg.ProxyFunc()
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		return cfgProxyFunc(req.URL)
	}

	cli, err := hw_fusion.NewHwFusionComputeClient(
		hw_fusion.NewHwFusionClientConfig(
			options.EndPoint,
			options.User,
			options.Password,
		).Debug(options.Debug).
			CloudproviderConfig(
				cloudprovider.ProviderConfig{
					ProxyFunc: proxyFunc,
				},
			),
	)
	if err != nil {
		return nil, err
	}

	return cli.GetRegion(options.RegionId)
}

func main() {
	// 获取解析器
	parser, e := getSubcommandParser()
	if e != nil {
		showErrorAndExit(e)
	}
	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*BaseOptions)

	if options.Help {
		fmt.Print(parser.HelpString())
	} else {
		subcmd := parser.GetSubcommand()
		subparser := subcmd.GetSubParser()
		if e != nil {
			if subparser != nil {
				fmt.Print(subparser.Usage())
			} else {
				fmt.Print(parser.Usage())
			}
			showErrorAndExit(e)
		} else {
			suboptions := subparser.Options()
			if options.SUBCOMMAND == "help" {
				e = subcmd.Invoke(suboptions)
			} else {
				var region *hw_fusion.SRegion
				region, e = newClient(options)
				if e != nil {
					showErrorAndExit(e)
				}
				e = subcmd.Invoke(region, suboptions)
			}
			if e != nil {
				showErrorAndExit(e)
			}
		}
	}
}
