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

package identity

import (
	"yunion.io/x/onecloud/pkg/mcclient/modulebase"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

var (
	Services   modulebase.ResourceManager
	ServicesV3 modulebase.ResourceManager
)

// 说明服务注册这一步并不是在main函数的执行流程里面的
func init() {
	Services = modules.NewIdentityManager("OS-KSADM:service",
		"OS-KSADM:services",
		[]string{},
		[]string{"ID", "Name", "Type", "Description"})

	modules.Register(&Services)

	ServicesV3 = modules.NewIdentityV3Manager("service",
		"services",
		[]string{},
		[]string{"ID", "Name", "Type", "Description"})

	modules.Register(&ServicesV3)
}
