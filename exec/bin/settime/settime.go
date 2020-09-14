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

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/chaosblade-io/chaosblade-exec-os/exec/bin"
	"github.com/chaosblade-io/chaosblade-spec-go/channel"
)

var dstDateTime string

var setDateTime, clearDateTime bool

func main() {
	flag.StringVar(&dstDateTime, "datetime", "", "datetime")
	flag.BoolVar(&setDateTime, "start", false, "set datetime to destination value")
	flag.BoolVar(&clearDateTime, "stop", false, "sync datetime from ntp server")
	bin.ParseFlagAndInitLog()

	if setDateTime == clearDateTime {
		bin.PrintErrAndExit("must add --start or --stop flag")
	}

	if setDateTime {
		if dstDateTime == "" {
			bin.PrintErrAndExit("must add --datetime flag")
		}
		doSetDateTime(dstDateTime)
	} else if clearDateTime {
		doClearDateTime()
	} else {
		bin.PrintErrAndExit("less --start or --stop flag")
	}
}

var cl = channel.NewLocalChannel()

func doSetDateTime(datetime string) {
	var ctx = context.WithValue(context.Background(), channel.ExcludeProcessKey, "blade")
	// systemd may need few seconds to stop systemd-timesyncd service
	args := fmt.Sprintf(`--no-pager --no-ask-password set-ntp false && \
		sleep 5 && \
		timedatectl --no-pager --no-ask-password set-time "%s"`, datetime)
	response := channel.NewLocalChannel().Run(ctx, "timedatectl", args)
	if !response.Success {
		bin.PrintErrAndExit(response.Err)
	}
	bin.PrintOutputAndExit(response.Result.(string))
}

func doClearDateTime() {
	var ctx = context.WithValue(context.Background(), channel.ExcludeProcessKey, "blade")
	args := fmt.Sprintf("--no-pager --no-ask-password set-ntp true")
	response := channel.NewLocalChannel().Run(ctx, "timedatectl", args)
	if !response.Success {
		bin.PrintErrAndExit(response.Err)
	}
	bin.PrintOutputAndExit(response.Result.(string))
}
