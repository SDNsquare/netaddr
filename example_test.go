// Copyright 2020 The Inet.Af AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netaddr_test

import (
	"fmt"
	"os"
	"text/tabwriter"

	"SDNsquare/netaddr"
)

func ExampleIP() {
	ip, err := netaddr.ParseIP("192.0.2.3")
	if err != nil {
		panic(err)
	}

	// netaddr.IP supports comparison using ==
	fmt.Println(ip == netaddr.IPv4(192, 0, 2, 3))

	// netaddr.IP can be used as a map key
	hosts := map[netaddr.IP]string{ip: "example.net"}
	fmt.Println(hosts)
	// Output:
	// true
	// map[192.0.2.3:example.net]
}

func ExampleIP_properties() {
	var zeroIP netaddr.IP
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "String()\tZone()\tIsZero()\tIs4()\tIs6()\tIs4in6()")
	for _, ip := range []netaddr.IP{
		zeroIP,
		netaddr.MustParseIP("192.0.2.3"),
		netaddr.MustParseIP("2001:db8::68"),
		netaddr.MustParseIP("2001:db8::68%eth0"),
		netaddr.MustParseIP("::ffff:192.0.2.3"),
	} {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", ip, ip.Zone(), ip.IsZero(), ip.Is4(), ip.Is6(), ip.Is4in6())
	}
	w.Flush()
	// Output:
	// String()           Zone()  IsZero()  Is4()  Is6()  Is4in6()
	// invalid IP                 true      false  false  false
	// 192.0.2.3                  false     true   false  false
	// 2001:db8::68               false     false  true   false
	// 2001:db8::68%eth0  eth0    false     false  true   false
	// ::ffff:c000:203            false     false  true   true
}

func ExampleIP_Is4() {
	var zeroIP netaddr.IP
	ipv4 := netaddr.MustParseIP("192.0.2.3")
	ipv6 := netaddr.MustParseIP("2001:db8::68")
	ip4in6 := netaddr.MustParseIP("::ffff:192.0.2.3")

	fmt.Printf("IP{}.Is4() -> %v\n", zeroIP.Is4())
	fmt.Printf("(%v).Is4() -> %v\n", ipv4, ipv4.Is4())
	fmt.Printf("(%v).Is4() -> %v\n", ipv6, ipv6.Is4())
	fmt.Printf("(%v).Is4() -> %v\n", ip4in6, ip4in6.Is4())
	// Output:
	// IP{}.Is4() -> false
	// (192.0.2.3).Is4() -> true
	// (2001:db8::68).Is4() -> false
	// (::ffff:c000:203).Is4() -> false
}

func ExampleIP_Is4in6() {
	var zeroIP netaddr.IP
	ipv4 := netaddr.MustParseIP("192.0.2.3")
	ipv6 := netaddr.MustParseIP("2001:db8::68")
	ip4in6 := netaddr.MustParseIP("::ffff:192.0.2.3")

	fmt.Printf("IP{}.Is4in6() -> %v\n", zeroIP.Is4in6())
	fmt.Printf("(%v).Is4in6() -> %v\n", ipv4, ipv4.Is4in6())
	fmt.Printf("(%v).Is4in6() -> %v\n", ipv6, ipv6.Is4in6())
	fmt.Printf("(%v).Is4in6() -> %v\n", ip4in6, ip4in6.Is4in6())
	// Output:
	// IP{}.Is4in6() -> false
	// (192.0.2.3).Is4in6() -> false
	// (2001:db8::68).Is4in6() -> false
	// (::ffff:c000:203).Is4in6() -> true
}

func ExampleIP_Is6() {
	var zeroIP netaddr.IP
	ipv4 := netaddr.MustParseIP("192.0.2.3")
	ipv6 := netaddr.MustParseIP("2001:db8::68")
	ip4in6 := netaddr.MustParseIP("::ffff:192.0.2.3")

	fmt.Printf("IP{}.Is6() -> %v\n", zeroIP.Is4in6())
	fmt.Printf("(%v).Is6() -> %v\n", ipv4, ipv4.Is6())
	fmt.Printf("(%v).Is6() -> %v\n", ipv6, ipv6.Is6())
	fmt.Printf("(%v).Is6() -> %v\n", ip4in6, ip4in6.Is6())
	// Output:
	// IP{}.Is6() -> false
	// (192.0.2.3).Is6() -> false
	// (2001:db8::68).Is6() -> true
	// (::ffff:c000:203).Is6() -> true
}

func ExampleIP_IsZero() {
	var zeroIP netaddr.IP
	ipv4AllZeroes := netaddr.MustParseIP("0.0.0.0")
	ipv6AllZeroes := netaddr.MustParseIP("::")

	fmt.Printf("IP{}.IsZero() -> %v\n", zeroIP.IsZero())
	fmt.Printf("(%v).IsZero() -> %v\n", ipv4AllZeroes, ipv4AllZeroes.IsZero())
	fmt.Printf("(%v).IsZero() -> %v\n", ipv6AllZeroes, ipv6AllZeroes.IsZero())
	// Output:
	// IP{}.IsZero() -> true
	// (0.0.0.0).IsZero() -> false
	// (::).IsZero() -> false
}

func ExampleIP_String() {
	ipv4 := netaddr.MustParseIP("192.0.2.3")
	ipv6 := netaddr.MustParseIP("2001:db8::68")
	ip4in6 := netaddr.MustParseIP("::ffff:192.0.2.3")

	fmt.Printf("(%v).String() -> %v\n", ipv4, ipv4.String())
	fmt.Printf("(%v).String() -> %v\n", ipv6, ipv6.String())
	fmt.Printf("(%v).String() -> %v\n", ip4in6, ip4in6.String())
	// Output:
	// (192.0.2.3).String() -> 192.0.2.3
	// (2001:db8::68).String() -> 2001:db8::68
	// (::ffff:c000:203).String() -> ::ffff:c000:203
}

func ExampleIP_Unmap() {
	ipv4 := netaddr.MustParseIP("192.0.2.3")
	ipv6 := netaddr.MustParseIP("2001:db8::68")
	ip4in6 := netaddr.MustParseIP("::ffff:192.0.2.3")

	fmt.Printf("(%v).Unmap() -> %v\n", ipv4, ipv4.Unmap())
	fmt.Printf("(%v).Unmap() -> %v\n", ipv6, ipv6.Unmap())
	fmt.Printf("(%v).Unmap() -> %v\n", ip4in6, ip4in6.Unmap())
	// Output:
	// (192.0.2.3).Unmap() -> 192.0.2.3
	// (2001:db8::68).Unmap() -> 2001:db8::68
	// (::ffff:c000:203).Unmap() -> 192.0.2.3
}

func ExampleIP_WithZone() {
	ipv4 := netaddr.MustParseIP("192.0.2.3")
	ipv6 := netaddr.MustParseIP("2001:db8::68")
	ipv6Zoned := netaddr.MustParseIP("2001:db8::68%eth0")

	fmt.Printf("(%v).WithZone(\"newzone\") -> %v\n", ipv4, ipv4.WithZone("newzone"))
	fmt.Printf("(%v).WithZone(\"newzone\") -> %v\n", ipv6, ipv6.WithZone("newzone"))
	fmt.Printf("(%v).WithZone(\"newzone\") -> %v\n", ipv6Zoned, ipv6Zoned.WithZone("newzone"))
	fmt.Printf("(%v).WithZone(\"\") -> %v\n", ipv6Zoned, ipv6Zoned.WithZone(""))
	// Output:
	// (192.0.2.3).WithZone("newzone") -> 192.0.2.3
	// (2001:db8::68).WithZone("newzone") -> 2001:db8::68%newzone
	// (2001:db8::68%eth0).WithZone("newzone") -> 2001:db8::68%newzone
	// (2001:db8::68%eth0).WithZone("") -> 2001:db8::68
}

func ExampleIPSet() {
	var b netaddr.IPSetBuilder

	b.AddPrefix(netaddr.MustParseIPPrefix("10.0.0.0/8"))
	b.RemovePrefix(netaddr.MustParseIPPrefix("10.0.0.0/16"))

	b.AddRange(netaddr.IPRangeFrom(
		netaddr.MustParseIP("fed0::0400"),
		netaddr.MustParseIP("fed0::04ff"),
	))

	s := b.IPSet()

	fmt.Println("Ranges:")
	for _, r := range s.Ranges() {
		fmt.Printf("  %s - %s\n", r.From(), r.To())
	}

	fmt.Println("Prefixes:")
	for _, p := range s.Prefixes() {
		fmt.Printf("  %s\n", p)
	}
	// Output:
	// Ranges:
	//   10.1.0.0 - 10.255.255.255
	//   fed0::400 - fed0::4ff
	// Prefixes:
	//   10.1.0.0/16
	//   10.2.0.0/15
	//   10.4.0.0/14
	//   10.8.0.0/13
	//   10.16.0.0/12
	//   10.32.0.0/11
	//   10.64.0.0/10
	//   10.128.0.0/9
	//   fed0::400/120
}
