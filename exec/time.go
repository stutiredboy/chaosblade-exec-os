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
	"context"
	"fmt"
	"path"

	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
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
				&setActionCommand{
					spec.BaseExpActionCommandSpec{
						ActionMatchers: []spec.ExpFlagSpec{},
						ActionFlags:    []spec.ExpFlagSpec{},
						ActionExecutor: &timeExecutor{},
					},
				},
			},
			ExpFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:     "datetime",
					Desc:     "datetime to set, such as '2020-09-09 16:53:00'",
					Required: true,
				},
			},
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
	return "time set --datetime='2020-09-09 16:48:00'"
}

// setActionCommand is the Action to run in time executor
type setActionCommand struct {
	spec.BaseExpActionCommandSpec
}

func (*setActionCommand) Name() string {
	return "set"
}

func (*setActionCommand) Aliases() []string {
	return []string{}
}

func (*setActionCommand) ShortDesc() string {
	return "time set"
}

func (*setActionCommand) LongDesc() string {
	return "set OS time to specific date time"
}

func (*setActionCommand) Matchers() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

func (*setActionCommand) Flags() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

// timeExecutor is one of the OS executor
type timeExecutor struct {
	channel spec.Channel
}

func (te *timeExecutor) Name() string {
	return "time"
}

func (te *timeExecutor) SetChannel(channel spec.Channel) {
	te.channel = channel
}

// setTimeBin is the command of settime
const setTimeBin = "chaos_settime"

func (te *timeExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	err := checkTimeExpEnv()
	if err != nil {
		return spec.ReturnFail(spec.Code[spec.CommandNotFound], err.Error())
	}
	if te.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}
	dateTime := model.ActionFlags["datetime"]
	if _, ok := spec.IsDestroy(ctx); ok {
		return te.stop(ctx)
	}

	return te.start(ctx, dateTime)
}

func (te *timeExecutor) start(ctx context.Context, dateTime string) *spec.Response {
	args := fmt.Sprintf("--start --debug=%t --datetime='%s'", util.Debug, dateTime)
	return te.channel.Run(ctx, path.Join(te.channel.GetScriptPath(), setTimeBin), args)
}

func (te *timeExecutor) stop(ctx context.Context) *spec.Response {
	args := fmt.Sprintf("--stop --debug=%t", util.Debug)
	return te.channel.Run(ctx, path.Join(te.channel.GetScriptPath(), setTimeBin), args)
}

// checkTimeExpEnv check the commands depended exists or not.
func checkTimeExpEnv() error {
	commands := []string{"timedatectl"}
	for _, command := range commands {
		if !channel.NewLocalChannel().IsCommandAvailable(command) {
			return fmt.Errorf("%s command not found", command)
		}
	}
	if channel.NewLocalChannel().IsCommandAvailable("ntpd") {
		return fmt.Errorf("only support systemd-timesyncd, so ntpd should no be installed")
	}
	return nil
}
