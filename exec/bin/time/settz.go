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
	"io/ioutil"
	"os"

	"github.com/chaosblade-io/chaosblade-exec-os/exec/bin"
	"github.com/chaosblade-io/chaosblade-spec-go/channel"
)

var dstTimeZone, expUid string

var setTimeZone, resetTimeZone bool

func main() {
	flag.StringVar(&dstTimeZone, "timezone", "", "timezone")
	flag.StringVar(&expUid, "uid", "settz", "uid")
	flag.BoolVar(&setTimeZone, "start", false, "set timezone to destination value")
	flag.BoolVar(&resetTimeZone, "stop", false, "set timezone to the value before experiment")
	bin.ParseFlagAndInitLog()

	if setTimeZone == resetTimeZone {
		bin.PrintErrAndExit("must add --start or --stop flag")
	}

	if setTimeZone {
		if dstTimeZone == "" {
			bin.PrintErrAndExit("must add --timezone flag")
		}
		doSetTimeZone(expUid, dstTimeZone)
	} else if resetTimeZone {
		doResetTimeZone(expUid)
	} else {
		bin.PrintErrAndExit("less --start or --stop flag")
	}
}

var cl = channel.NewLocalChannel()

func doSetTimeZone(uid string, timezone string) {
	ctx := context.Background()
	tmpFile := fmt.Sprintf("/tmp/chaos-settz-%s.tmp", uid)

	// get current timezone and save it to tmpFile
	response := cl.Run(ctx, "timedatectl", fmt.Sprintf(`| fgrep "Time zone:" | awk '{print $3}' > %s`, tmpFile))
	if !response.Success {
		bin.PrintErrAndExit(response.Err)
	}
	args := fmt.Sprintf("--no-pager --no-ask-password set-timezone %s", timezone)
	response = channel.NewLocalChannel().Run(ctx, "timedatectl", args)
	if !response.Success {
		os.Remove(tmpFile)
		bin.PrintErrAndExit(response.Err)
	}
	bin.PrintOutputAndExit(response.Result.(string))
}

func doResetTimeZone(uid string) {
	ctx := context.Background()
	tmpFile := fmt.Sprintf("/tmp/chaos-settz-%s.tmp", uid)
	defer os.Remove(tmpFile)

	tzinfo, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		bin.PrintErrAndExit(err.Error())
	}
	args := fmt.Sprintf("--no-pager --no-ask-password set-timezone %s", string(tzinfo[:]))
	response := channel.NewLocalChannel().Run(ctx, "timedatectl", args)
	if !response.Success {
		bin.PrintErrAndExit(response.Err)
	}
	bin.PrintOutputAndExit(response.Result.(string))
}
