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

	"yunion.io/x/onecloud/pkg/cloudprovider"
	"yunion.io/x/onecloud/pkg/multicloud/apsara"
	_ "yunion.io/x/onecloud/pkg/multicloud/apsara/shell"
	"yunion.io/x/onecloud/pkg/util/shellutils"
)

// Alibaba Cloud Apsara Stack 阿里云专用云（分布式结构）
type BaseOptions struct {
	Debug                          bool   `help:"debug mode"`
	Help                           bool   `help:"Show help"`
	AccessKey                      string `help:"Access key" default:"$APSARA_ACCESS_KEY" metavar:"APSARA_ACCESS_KEY"`
	Secret                         string `help:"Secret" default:"$APSARA_SECRET" metavar:"APSARA_SECRET"`
	cloudprovider.SApsaraEndpoints `help:"Endpoints for Apsara"`
	Endpoint                       string `help:"Apsara endpoint" default:"$APSARA_ENDPOINT" metavar:"APSARA_ENDPOINT"`
	RegionId                       string `help:"RegionId" default:"$APSARA_REGION" metavar:"APSARA_REGION"`
	DEFAULT_REGION                 string `help:"Default region" default:"$APSARA_DEFAULT_REGION"`
	SUBCOMMAND                     string `help:"apsaracli subcommand" subcommand:"true"`
}

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	parse, e := structarg.NewArgumentParser(&BaseOptions{},
		"apsaracli",
		"Command-line interface to apsara API.",
		`See "apsaracli help COMMAND" for help on a specific command.`)

	if e != nil {
		return nil, e
	}

	subcmd := parse.GetSubcommand()
	if subcmd == nil {
		return nil, fmt.Errorf("No subcommand argument.")
	}
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

func newClient(options *BaseOptions) (*apsara.SRegion, error) {
	if len(options.AccessKey) == 0 {
		return nil, fmt.Errorf("Missing accessKey")
	}

	if len(options.Secret) == 0 {
		return nil, fmt.Errorf("Missing secret")
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

	cli, err := apsara.NewApsaraClient(
		apsara.NewApsaraClientConfig(
			options.AccessKey,
			options.Secret,
			options.Endpoint,
			options.SApsaraEndpoints,
		).Debug(options.Debug).
			CloudproviderConfig(
				cloudprovider.ProviderConfig{
					ProxyFunc:     proxyFunc,
					URL:           options.Endpoint,
					DefaultRegion: options.DEFAULT_REGION,
				},
			),
	)
	if err != nil {
		return nil, err
	}

	region := cli.GetRegion(options.RegionId)
	if region == nil {
		return nil, fmt.Errorf("No such region %s", options.RegionId)
	}

	return region, nil
}

func main() {
	parser, e := getSubcommandParser()
	if e != nil {
		showErrorAndExit(e)
	}
	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*BaseOptions)

	if options.Help {
		fmt.Print(parser.HelpString())
		return
	}
	subcmd := parser.GetSubcommand()
	subparser := subcmd.GetSubParser()
	if e != nil {
		if subparser != nil {
			fmt.Print(subparser.Usage())
		} else {
			fmt.Print(parser.Usage())
		}
		showErrorAndExit(e)
	}
	suboptions := subparser.Options()
	if options.SUBCOMMAND == "help" {
		e = subcmd.Invoke(suboptions)
	} else {
		var region *apsara.SRegion
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
