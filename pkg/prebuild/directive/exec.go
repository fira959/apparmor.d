// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Exec struct {
	cfg.Base
}

func init() {
	RegisterDirective(&Exec{
		Base: cfg.Base{
			Keyword: "exec",
			Msg:     "Exec directive applied",
			Help:    Keyword + `exec [P|U|p|u|PU|pu|] profiles...`,
		},
	})
}

func (d Exec) Apply(opt *Option, profileRaw string) string {
	transition := "Px"
	transitions := []string{"P", "U", "p", "u", "PU", "pu"}
	t := opt.ArgList[0]
	if slices.Contains(transitions, t) {
		transition = t + "x"
		delete(opt.ArgMap, t)
	}

	rules := aa.Rules{}
	for name := range opt.ArgMap {
		profiletoTransition := util.MustReadFile(cfg.RootApparmord.Join(name))
		dstProfile := aa.DefaultTunables()
		dstProfile.ParseVariables(profiletoTransition)
		for _, variable := range dstProfile.Variables {
			if variable.Name == "exec_path" {
				for _, v := range variable.Values {
					rules = append(rules, &aa.File{
						Path:   v,
						Access: []string{transition},
					})
				}
				break
			}
		}
	}

	aa.TemplateIndentationLevel = strings.Count(
		strings.SplitN(opt.Raw, Keyword, 1)[0], aa.TemplateIndentation,
	)
	rules.Sort()
	new := rules.String()
	new = new[:len(new)-1]
	return strings.Replace(profileRaw, opt.Raw, new, -1)
}
