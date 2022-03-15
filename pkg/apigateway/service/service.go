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

package service

import (
	"net"
	"os"
	"strconv"

	"yunion.io/x/log"

	"yunion.io/x/onecloud/pkg/apigateway/app"
	"yunion.io/x/onecloud/pkg/apigateway/clientman"
	"yunion.io/x/onecloud/pkg/apigateway/options"
	api "yunion.io/x/onecloud/pkg/apis/apigateway"
	app_common "yunion.io/x/onecloud/pkg/cloudcommon/app"
	common_options "yunion.io/x/onecloud/pkg/cloudcommon/options"
	"yunion.io/x/onecloud/pkg/mcclient"
)

func StartService() {
	options.Options = &options.GatewayOptions{}
	opts := options.Options
	baseOpts := &opts.BaseOptions
	commonOpts := &opts.CommonOptions
	common_options.ParseOptions(opts, os.Args, "./apigateway.conf", api.SERVICE_TYPE)

	// 认证模块初始化（异步），初步猜测cloudmon是针对所有云平台的监控相关的共有代码模块，该认证是apigateway作为keystone的admin访问keystone
	app_common.InitAuth(commonOpts, func() {
		log.Infof("Auth complete.")
	})

	// 同步所有底层服务的配置，目测 apigateway 这个服务组件是所有服务组件中最晚启动的，并且底层服务的配置变动或地址变动应该能够同步到 apigateway
	// 所有api底层的manager注册的地方
	common_options.StartOptionManager(opts, opts.ConfigSyncPeriodSeconds, api.SERVICE_TYPE, api.SERVICE_VERSION, options.OnOptionsChange)

	if opts.DisableModuleApiVersion {
		mcclient.DisableApiVersionByModule()
	}

	// 读取SSL密钥文件，设置rsa密钥
	if err := clientman.InitClient(); err != nil {
		log.Fatalf("Init client token manager: %v", err)
	}

	serviceApp := app.NewApp(app_common.InitApp(baseOpts, false))
	serviceApp.InitHandlers().Bind()

	// mods, jmods := modulebase.GetRegisterdModules()
	// log.Infof("Modules: %s", jsonutils.Marshal(mods).PrettyString())
	// log.Infof("Modules: %s", jsonutils.Marshal(jmods).PrettyString())

	listenAddr := net.JoinHostPort(options.Options.Address, strconv.Itoa(options.Options.Port))
	if opts.EnableSsl {
		serviceApp.ListenAndServeTLS(listenAddr, opts.SslCertfile, opts.SslKeyfile)
	} else {
		serviceApp.ListenAndServe(listenAddr)
	}
}
