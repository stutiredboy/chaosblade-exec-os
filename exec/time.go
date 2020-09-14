/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package exec

import (
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

// TimeCommandModelSpec is time model spec
type TimeCommandModelSpec struct {
	spec.BaseExpModelCommandSpec
}

// NewTimeCommandModelSpec generate TimeCommandModelSpec
func NewTimeCommandModelSpec() spec.ExpModelCommandSpec {
	return &TimeCommandModelSpec{
		spec.BaseExpModelCommandSpec{
			ExpActions: []spec.ExpActionCommandSpec{
        NewSetTimeActionCommandSpec(),
        NewSetTimeZoneActionCommandSpec(),
			},
			ExpFlags: []spec.ExpFlagSpec{},
		},
	}
}

func (*TimeCommandModelSpec) Name() string {
	return "time"
}

func (*TimeCommandModelSpec) ShortDesc() string {
	return "time experiment"
}

func (*TimeCommandModelSpec) LongDesc() string {
	return "time experiment, for example set OS time to a wrong value"
}

func (*TimeCommandModelSpec) Example() string {
	return "time settime --datetime='2020-09-09 16:48:00'"
}
