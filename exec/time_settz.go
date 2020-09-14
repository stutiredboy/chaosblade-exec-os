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

// SetTimeZoneActionCommandSpec is the Action to run in time executor
type SetTimeZoneActionCommandSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewSetTimeZoneActionCommandSpec() spec.ExpActionCommandSpec {
	return &SetTimeZoneActionCommandSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: []spec.ExpFlagSpec{},
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:     "timezone",
					Desc:     "timezone, such as: Asia/Hong_kong",
					Required: true,
				},
			},
			ActionExecutor: &TimeZoneExecutor{},
		},
	}
}

func (*SetTimeZoneActionCommandSpec) Name() string {
	return "settz"
}

func (*SetTimeZoneActionCommandSpec) Aliases() []string {
	return []string{}
}

func (*SetTimeZoneActionCommandSpec) ShortDesc() string {
	return "set timezone action"
}

func (*SetTimeZoneActionCommandSpec) LongDesc() string {
	return "set OS timezone to specific timezone"
}

func (*SetTimeZoneActionCommandSpec) Matchers() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

func (*SetTimeZoneActionCommandSpec) Flags() []spec.ExpFlagSpec {
	return []spec.ExpFlagSpec{}
}

// TimeZoneExecutor is one of the OS executor
type TimeZoneExecutor struct {
	channel spec.Channel
}

func (te *TimeZoneExecutor) Name() string {
	return "time"
}

func (te *TimeZoneExecutor) SetChannel(channel spec.Channel) {
	te.channel = channel
}

// setTZBin is the command of settime
const setTZBin = "chaos_settz"

func (te *TimeZoneExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	err := checkTZExpEnv()
	if err != nil {
		return spec.ReturnFail(spec.Code[spec.CommandNotFound], err.Error())
	}
	if te.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}
	timeZone := model.ActionFlags["timezone"]
	if _, ok := spec.IsDestroy(ctx); ok {
		return te.stop(ctx, uid)
	}

	return te.start(ctx, uid, timeZone)
}

func (te *TimeZoneExecutor) start(ctx context.Context, uid string, timeZone string) *spec.Response {
	args := fmt.Sprintf("--start --debug=%t --timezone='%s' --uid=%s", util.Debug, timeZone, uid)
	return te.channel.Run(ctx, path.Join(te.channel.GetScriptPath(), setTZBin), args)
}

func (te *TimeZoneExecutor) stop(ctx context.Context, uid string) *spec.Response {
	args := fmt.Sprintf("--stop --debug=%t --uid=%s", util.Debug, uid)
	return te.channel.Run(ctx, path.Join(te.channel.GetScriptPath(), setTZBin), args)
}

// checkTZExpEnv check the commands depended exists or not.
func checkTZExpEnv() error {
	commands := []string{"timedatectl", "fgrep", "awk"}
	for _, command := range commands {
		if !channel.NewLocalChannel().IsCommandAvailable(command) {
			return fmt.Errorf("%s command not found", command)
		}
	}
	return nil
}
