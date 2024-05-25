// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
)

const tokPTRACE = "ptrace"

func init() {
	requirements[tokPTRACE] = requirement{
		"access": []string{
			"r", "w", "rw", "read", "readby", "trace", "tracedby",
		},
	}
}

type Ptrace struct {
	RuleBase
	Qualifier
	Access []string
	Peer   string
}

func newPtraceFromLog(log map[string]string) Rule {
	return &Ptrace{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(tokPTRACE, log["requested_mask"])),
		Peer:      log["peer"],
	}
}

func (r *Ptrace) Less(other any) bool {
	o, _ := other.(*Ptrace)
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if r.Peer != o.Peer {
		return r.Peer == o.Peer
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Ptrace) Equals(other any) bool {
	o, _ := other.(*Ptrace)
	return slices.Equal(r.Access, o.Access) && r.Peer == o.Peer &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *Ptrace) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Ptrace) Constraint() constraint {
	return blockKind
}

func (r *Ptrace) Kind() string {
	return tokPTRACE
}
