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

package app

import (
	"context"
	"time"

	"yunion.io/x/log"
	"yunion.io/x/pkg/utils"

	"yunion.io/x/onecloud/pkg/apis"
	"yunion.io/x/onecloud/pkg/apis/identity"
	"yunion.io/x/onecloud/pkg/cloudcommon/consts"
	"yunion.io/x/onecloud/pkg/cloudcommon/notifyclient"
	common_options "yunion.io/x/onecloud/pkg/cloudcommon/options"
	"yunion.io/x/onecloud/pkg/cloudcommon/policy"
	"yunion.io/x/onecloud/pkg/mcclient/auth"
	identity_modules "yunion.io/x/onecloud/pkg/mcclient/modules/identity"
)

func InitAuth(options *common_options.CommonOptions, authComplete auth.AuthCompletedCallback) {

	if len(options.AuthURL) == 0 {
		log.Fatalln("Missing AuthURL")
	}

	if len(options.AdminUser) == 0 {
		log.Fatalln("Mising AdminUser")
	}

	if len(options.AdminPassword) == 0 {
		log.Fatalln("Missing AdminPasswd")
	}

	if len(options.AdminProject) == 0 {
		log.Fatalln("Missing AdminProject")
	}

	a := auth.NewAuthInfo(
		options.AuthURL,
		options.AdminDomain,
		options.AdminUser,
		options.AdminPassword,
		options.AdminProject,
		options.AdminProjectDomain,
	)

	// debug := options.LogLevel == "debug"

	if options.SessionEndpointType != "" {
		if !utils.IsInStringArray(options.SessionEndpointType,
			[]string{identity.EndpointInterfacePublic, identity.EndpointInterfaceInternal}) {
			log.Fatalf("Invalid session endpoint type %s", options.SessionEndpointType)
		}
		auth.SetEndpointType(options.SessionEndpointType)
	}

	// 异步初始化一个 mcclient ，传给 authManager 作为认证备用，使用 LRU 算法缓存keystone获取的权限数据记录
	auth.Init(a, options.DebugClient, true, options.SslCertfile, options.SslKeyfile) // , authComplete)

	// 设置通知用户和通知组，notify client模块
	users := options.NotifyAdminUsers
	groups := options.NotifyAdminGroups
	if len(users) == 0 && len(groups) == 0 {
		users = []string{"sysadmin"}
	}
	notifyclient.FetchNotifyAdminRecipients(context.Background(), options.Region, users, groups)

	// 回调函数 callback ，打印 Auth complete.
	if authComplete != nil {
		authComplete()
	}

	// 设置缓存过期时间
	consts.SetTenantCacheExpireSeconds(options.TenantCacheExpireSeconds)

	// 初始化 rbac 策略计算模块，加载缓存 policymanager: pkg/cloudcommon/policy/policy.go
	InitBaseAuth(&options.BaseOptions)

	// endpoint更新管理器（两小时更新一次），监听资源管理器，执行某些子任务
	watcher := newEndpointChangeManager()
	watcher.StartWatching(&identity_modules.EndpointsV3)

	// 开启etcd puller
	startEtcdEndpointPuller()
}

func InitBaseAuth(options *common_options.BaseOptions) {
	policy.EnableGlobalRbac(
		time.Second*time.Duration(options.RbacPolicyRefreshIntervalSeconds),
		options.RbacDebug,
	)
	consts.SetNonDefaultDomainProjects(options.NonDefaultDomainProjects)
}

func FetchEtcdServiceInfo() (*identity.EndpointDetails, error) {
	s := auth.GetAdminSession(context.Background(), "", "")
	return s.GetCommonEtcdEndpoint()
}

func startEtcdEndpointPuller() {
	// 疑问点，这个etcdURL获取后在哪里使用
	retryInterval := 60
	etecdUrl, err := auth.GetServiceURL(apis.SERVICE_TYPE_ETCD, consts.GetRegion(), "", "")
	if err != nil {
		log.Errorf("[etcd] GetServiceURL fail %s, retry after %d seconds", err, retryInterval)
	} else if len(etecdUrl) == 0 {
		log.Errorf("[etcd] no service url found, retry after %d seconds", retryInterval)
	} else {
		return
	}
	auth.ReAuth()
	time.AfterFunc(time.Second*time.Duration(retryInterval), startEtcdEndpointPuller)
}
