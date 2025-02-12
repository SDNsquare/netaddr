// Copyright 2020 The Inet.Af AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netaddr

import "sort"

// IPSetBuilder builds an immutable IPSet.
//
// The zero value is a valid value representing a set of no IPs.
//
// The Add and Remove methods add or remove IPs to/from the set.
// Removals only affect the current membership of the set, so in
// general Adds should be called first. Input ranges may overlap in
// any way.
type IPSetBuilder struct {
	// in are the ranges in the set.
	in []IPRange

	// out are the ranges to be removed from 'in'.
	out []IPRange
}

// normalize normalizes s: s.in becomes the minimal sorted list of
// ranges required to describe s, and s.out becomes empty.
func (s *IPSetBuilder) normalize() {
	const debug = false
	if debug {
		debugf("ranges start in=%v out=%v", s.in, s.out)
	}
	in, ok := mergeIPRanges(s.in)
	if !ok {
		return
	}
	out, ok := mergeIPRanges(s.out)
	if !ok {
		return
	}
	if debug {
		debugf("ranges sort  in=%v out=%v", in, out)
	}

	// in and out are sorted in ascending range order, and have no
	// overlaps within each other. We can run a merge of the two lists
	// in one pass.

	min := make([]IPRange, 0, len(in))
	for len(in) > 0 && len(out) > 0 {
		rin, rout := in[0], out[0]
		if debug {
			debugf("step in=%v out=%v", rin, rout)
		}

		switch {
		case !rout.Valid() || !rin.Valid():
			// mergeIPRanges should have prevented invalid ranges from
			// sneaking in.
			panic("invalid IPRanges during Ranges merge")
		case rout.entirelyBefore(rin):
			// "out" is entirely before "in".
			//
			//    out         in
			// f-------t   f-------t
			out = out[1:]
			if debug {
				debugf("out before in; drop out")
			}
		case rin.entirelyBefore(rout):
			// "in" is entirely before "out".
			//
			//    in         out
			// f------t   f-------t
			min = append(min, rin)
			in = in[1:]
			if debug {
				debugf("in before out; append in")
				debugf("min=%v", min)
			}
		case rin.coveredBy(rout):
			// "out" entirely covers "in".
			//
			//       out
			// f-------------t
			//    f------t
			//       in
			in = in[1:]
			if debug {
				debugf("in inside out; drop in")
			}
		case rout.inMiddleOf(rin):
			// "in" entirely covers "out".
			//
			//       in
			// f-------------t
			//    f------t
			//       out
			min = append(min, IPRange{from: rin.from, to: rout.from.Prior()})
			// Adjust in[0], not ir, because we want to consider the
			// mutated range on the next iteration.
			in[0].from = rout.to.Next()
			out = out[1:]
			if debug {
				debugf("out inside in; split in, append first in, drop out, adjust second in")
				debugf("min=%v", min)
			}
		case rout.overlapsStartOf(rin):
			// "out" overlaps start of "in".
			//
			//   out
			// f------t
			//    f------t
			//       in
			in[0].from = rout.to.Next()
			// Can't move ir onto min yet, another later out might
			// trim it further. Just discard or and continue.
			out = out[1:]
			if debug {
				debugf("out cuts start of in; adjust in, drop out")
			}
		case rout.overlapsEndOf(rin):
			// "out" overlaps end of "in".
			//
			//           out
			//        f------t
			//    f------t
			//       in
			min = append(min, IPRange{from: rin.from, to: rout.from.Prior()})
			in = in[1:]
			if debug {
				debugf("merge out cuts end of in; append shortened in")
				debugf("min=%v", min)
			}
		default:
			// The above should account for all combinations of in and
			// out overlapping, but insert a panic to be sure.
			panic("unexpected additional overlap scenario")
		}
	}
	if len(in) > 0 {
		// Ran out of removals before the end of in.
		min = append(min, in...)
		if debug {
			debugf("min=%v", min)
		}
	}

	s.in = min
	s.out = nil
}

// Clone returns a copy of s that shares no memory with s.
func (s *IPSetBuilder) Clone() *IPSetBuilder {
	return &IPSetBuilder{
		in:  append([]IPRange(nil), s.in...),
		out: append([]IPRange(nil), s.out...),
	}
}

// Add adds ip to s.
func (s *IPSetBuilder) Add(ip IP) { s.AddRange(IPRange{ip, ip}) }

// AddPrefix adds all IPs in p to s.
func (s *IPSetBuilder) AddPrefix(p IPPrefix) { s.AddRange(p.Range()) }

// AddRange adds r to s.
// If r is not Valid, AddRange does nothing.
func (s *IPSetBuilder) AddRange(r IPRange) {
	if !r.Valid() {
		return
	}
	// If there are any removals (s.out), then we need to compact the set
	// first to get the order right.
	if len(s.out) > 0 {
		s.normalize()
	}
	s.in = append(s.in, r)
}

// AddSet adds all IPs in b to s.
func (s *IPSetBuilder) AddSet(b *IPSet) {
	for _, r := range b.rr {
		s.AddRange(r)
	}
}

// Remove removes ip from s.
func (s *IPSetBuilder) Remove(ip IP) { s.RemoveRange(IPRange{ip, ip}) }

// RemovePrefix removes all IPs in p from s.
func (s *IPSetBuilder) RemovePrefix(p IPPrefix) { s.RemoveRange(p.Range()) }

// RemoveRange removes all IPs in r from s.
func (s *IPSetBuilder) RemoveRange(r IPRange) {
	if r.Valid() {
		s.out = append(s.out, r)
	}
}

// RemoveSet removes all IPs in o from s.
func (s *IPSetBuilder) RemoveSet(b *IPSet) {
	for _, r := range b.rr {
		s.RemoveRange(r)
	}
}

// removeBuilder removes all IPs in b from s.
func (s *IPSetBuilder) removeBuilder(b *IPSetBuilder) {
	b.normalize()
	for _, r := range b.in {
		s.RemoveRange(r)
	}
}

// Complement updates s to contain the complement of its current
// contents.
func (s *IPSetBuilder) Complement() {
	s.normalize()
	s.out = s.in
	s.in = []IPRange{
		IPPrefix{ip: IPv4(0, 0, 0, 0), bits: 0}.Range(),
		IPPrefix{ip: IPv6Unspecified(), bits: 0}.Range(),
	}
}

// Intersect updates s to the set intersection of s and b.
func (s *IPSetBuilder) Intersect(b *IPSet) {
	var o IPSetBuilder
	o.Complement()
	o.RemoveSet(b)
	s.removeBuilder(&o)
}

func discardf(format string, args ...interface{}) {}

// debugf is reassigned by tests.
var debugf = discardf

// IPSet returns an immutable IPSet representing the current state of
// s. The builder remains usable after calling IPSet.
func (s *IPSetBuilder) IPSet() *IPSet {
	s.normalize()
	return &IPSet{
		rr: append([]IPRange{}, s.in...),
	}
}

// IPSet represents a set of IP addresses.
//
// IPSet is safe for concurrent use.
// The zero value is a valid value representing a set of no IPs.
// Use IPSetBuilder to construct IPSets.
type IPSet struct {
	// rr is the set of IPs that belong to this IPSet. The IPRanges
	// are normalized according to IPSetBuilder.normalize, meaning
	// they are a sorted, minimal representation (no overlapping
	// ranges, no contiguous ranges). The implementation of various
	// methods rely on this property.
	rr []IPRange
}

// Ranges returns the minimum and sorted set of IP
// ranges that covers s.
func (s *IPSet) Ranges() []IPRange {
	return append([]IPRange{}, s.rr...)
}

// Prefixes returns the minimum and sorted set of IP prefixes
// that covers s.
func (s *IPSet) Prefixes() []IPPrefix {
	out := make([]IPPrefix, 0, len(s.rr))
	for _, r := range s.rr {
		out = append(out, r.Prefixes()...)
	}
	return out
}

// Equal reports whether s and o represent the same set of IP
// addresses.
func (s *IPSet) Equal(o *IPSet) bool {
	if len(s.rr) != len(o.rr) {
		return false
	}
	for i := range s.rr {
		if s.rr[i] != o.rr[i] {
			return false
		}
	}
	return true
}

// Contains reports whether ip is in s.
func (s *IPSet) Contains(ip IP) bool {
	// TODO: data structure permitting more efficient lookups:
	// https://github.com/inetaf/netaddr/issues/139
	i := sort.Search(len(s.rr), func(i int) bool {
		return ip.Less(s.rr[i].from)
	})
	if i == 0 {
		return false
	}
	i--
	return s.rr[i].contains(ip)
}

// ContainsRange reports whether all IPs in r are in s.
func (s *IPSet) ContainsRange(r IPRange) bool {
	for _, x := range s.rr {
		if r.coveredBy(x) {
			return true
		}
	}
	return false
}

// ContainsPrefix reports whether all IPs in p are in s.
func (s *IPSet) ContainsPrefix(p IPPrefix) bool {
	return s.ContainsRange(p.Range())
}

// Overlaps reports whether any IP in b is also in s.
func (s *IPSet) Overlaps(b *IPSet) bool {
	// TODO: sorted ranges lets us do this in O(n+m)
	for _, r := range s.rr {
		for _, or := range b.rr {
			if r.Overlaps(or) {
				return true
			}
		}
	}
	return false
}

// OverlapsRange reports whether any IP in r is also in s.
func (s *IPSet) OverlapsRange(r IPRange) bool {
	// TODO: sorted ranges lets us do this more efficiently.
	for _, x := range s.rr {
		if x.Overlaps(r) {
			return true
		}
	}
	return false
}

// OverlapsPrefix reports whether any IP in p is also in s.
func (s *IPSet) OverlapsPrefix(p IPPrefix) bool {
	return s.OverlapsRange(p.Range())
}

// RemoveFreePrefix splits s into a Prefix of length bitLen and a new
// IPSet with that prefix removed.
//
// If no contiguous prefix of length bitLen exists in s,
// RemoveFreePrefix returns ok=false.
func (s *IPSet) RemoveFreePrefix(bitLen uint8) (p IPPrefix, newSet *IPSet, ok bool) {
	var bestFit IPPrefix
	for _, r := range s.rr {
		for _, prefix := range r.Prefixes() {
			if prefix.bits > bitLen {
				continue
			}
			if bestFit.ip.IsZero() || prefix.bits > bestFit.bits {
				bestFit = prefix
				if bestFit.bits == bitLen {
					// exact match, done.
					break
				}
			}
		}
	}

	if bestFit.ip.IsZero() {
		return IPPrefix{}, s, false
	}

	prefix := IPPrefix{ip: bestFit.ip, bits: bitLen}

	var b IPSetBuilder
	b.AddSet(s)
	b.RemovePrefix(prefix)
	return prefix, b.IPSet(), true
}
