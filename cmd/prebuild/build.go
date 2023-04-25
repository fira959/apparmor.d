// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regABI           = regexp.MustCompile(`abi <abi/[0-9.]*>,\n`)
	regAttachments   = regexp.MustCompile(`(profile .* @{exec_path})`)
	regFlag          = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {`)
)

// Set complain flag on all profiles
func BuildComplain(profile string) string {
	if !Complain {
		return profile
	}

	flags := []string{}
	matches := regFlag.FindStringSubmatch(profile)
	if len(matches) != 0 {
		flags = strings.Split(matches[1], ",")
		if util.InSlice("complain", flags) {
			return profile
		}
	}
	flags = append(flags, "complain")
	strFlags := " flags=(" + strings.Join(flags, ",") + ") {"

	// Remove all flags definition, then set manifest' flags
	profile = regFlag.ReplaceAllLiteralString(profile, "")
	return regProfileHeader.ReplaceAllLiteralString(profile, strFlags)
}

// Bypass userspace tools restriction
func BuildUserspace(profile string) string {
	p := aa.NewAppArmorProfile(profile)
	p.ParseVariables()
	p.ResolveAttachments()
	att := p.NestAttachments()
	matches := regAttachments.FindAllString(profile, -1)
	if len(matches) > 0 {
		strheader := strings.Replace(matches[0], "@{exec_path}", att, -1)
		return regAttachments.ReplaceAllLiteralString(profile, strheader)
	}
	return profile
}

// Remove abi header for distributions that don't support it
func BuildABI(profile string) string {
	switch Distribution {
	case "debian", "whonix":
		return regABI.ReplaceAllLiteralString(profile, "")
	default:
		return profile
	}
}