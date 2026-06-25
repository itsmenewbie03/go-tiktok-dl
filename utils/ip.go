// Package utils
package utils

import (
	"math/rand/v2"
	"net"
)

var usPrefixes = []string{
	"3.0.0.0/8",  // Amazon
	"4.0.0.0/8",  // Level 3
	"8.0.0.0/8",  // Level 3
	"12.0.0.0/8", // AT&T
	"13.0.0.0/8", // Xerox
	"18.0.0.0/8", // MIT
	"20.0.0.0/8", // DoD
	"23.0.0.0/8", // Akamai
	"44.0.0.0/8", // Amateur Radio
	"52.0.0.0/8", // Amazon AWS
}

func RandomUSIP() net.IP {
	cidr := usPrefixes[rand.IntN(len(usPrefixes))]

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	ip = ip.To4()

	for i := range 4 {
		maskByte := ipnet.Mask[i]
		ip[i] |= byte(rand.IntN(256)) &^ maskByte
	}

	return ip
}
