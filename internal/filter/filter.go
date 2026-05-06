// Package filter provides port filtering utilities for portwatch.
// It allows users to include or exclude specific ports or port ranges
// from scan and watch operations.
package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule represents a single port filter rule.
type Rule struct {
	Low  int
	High int
}

// Filter holds a set of inclusion and exclusion rules.
type Filter struct {
	include []Rule
	exclude []Rule
}

// New creates a Filter from include and exclude range strings.
// Each string is a comma-separated list of ports or ranges (e.g. "22,80,8000-9000").
func New(include, exclude string) (*Filter, error) {
	f := &Filter{}
	var err error

	if include != "" {
		f.include, err = parseRules(include)
		if err != nil {
			return nil, fmt.Errorf("include: %w", err)
		}
	}

	if exclude != "" {
		f.exclude, err = parseRules(exclude)
		if err != nil {
			return nil, fmt.Errorf("exclude: %w", err)
		}
	}

	return f, nil
}

// Allow returns true if the given port passes the filter rules.
// If include rules are set, the port must match at least one.
// If exclude rules are set, the port must not match any.
func (f *Filter) Allow(port int) bool {
	if len(f.include) > 0 {
		matched := false
		for _, r := range f.include {
			if port >= r.Low && port <= r.High {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	for _, r := range f.exclude {
		if port >= r.Low && port <= r.High {
			return false
		}
	}

	return true
}

// Apply returns a filtered copy of the given port list.
func (f *Filter) Apply(ports []int) []int {
	out := make([]int, 0, len(ports))
	for _, p := range ports {
		if f.Allow(p) {
			out = append(out, p)
		}
	}
	return out
}

func parseRules(s string) ([]Rule, error) {
	parts := strings.Split(s, ",")
	rules := make([]Rule, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			hi, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			if lo > hi {
				return nil, fmt.Errorf("range low > high in %q", part)
			}
			rules = append(rules, Rule{Low: lo, High: hi})
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port %q", part)
			}
			rules = append(rules, Rule{Low: p, High: p})
		}
	}
	return rules, nil
}
