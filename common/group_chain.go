/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2022-2026, daeuniverse Organization <dae@v2raya.org>
 */

package common

import (
	"fmt"
	"strings"
)

// GroupChain describes the only supported group-based chain form:
// group(NAME) -> one concrete node.
type GroupChain struct {
	Name       string
	EntryGroup string
	ExitLink   string
	Link       string
}

// ParseGroupChain recognizes group-based node chains without changing the
// parsing behavior of ordinary nodes and node-to-node chains.
func ParseGroupChain(link string) (*GroupChain, bool, error) {
	name, linklike := GetTagFromLinkLikePlaintext(link)
	parts := strings.Split(linklike, "->")
	groupAt := -1
	groupName := ""
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if !strings.HasPrefix(part, "group(") {
			continue
		}
		if !strings.HasSuffix(part, ")") {
			return nil, true, fmt.Errorf("invalid group chain entry %q", part)
		}
		groupAt = i
		groupName = strings.TrimSpace(part[len("group(") : len(part)-1])
		break
	}
	if groupAt < 0 {
		return nil, false, nil
	}
	if len(parts) != 2 || groupAt != 0 {
		return nil, true, fmt.Errorf("group chain must have the form group(NAME) -> node")
	}
	if groupName == "" {
		return nil, true, fmt.Errorf("group chain entry name is empty")
	}
	exitLink := strings.TrimSpace(parts[1])
	if exitLink == "" {
		return nil, true, fmt.Errorf("group chain exit node is empty")
	}
	if strings.HasPrefix(exitLink, "group(") {
		return nil, true, fmt.Errorf("group chain exit must be a concrete node")
	}
	return &GroupChain{
		Name:       strings.TrimSpace(name),
		EntryGroup: groupName,
		ExitLink:   exitLink,
		Link:       strings.TrimSpace(linklike),
	}, true, nil
}
