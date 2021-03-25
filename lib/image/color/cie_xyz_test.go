package color

import (
	"fmt"
	"strings"
	"testing"
)

func tupleStr(valueStrs []string) string {
	return "(" + strings.Join(valueStrs, ", ") + ")"
}

func fTupleStr(values ...float64) string {
	valueStrs := make([]string, 0, len(values))
	for _, v := range values {
		valueStrs = append(valueStrs, fmt.Sprintf("%.4f", v))
	}
	return tupleStr(valueStrs)
}

func hexTupleStr(values ...uint32) string {
	valueStrs := make([]string, 0, len(values))
	for _, v := range values {
		valueStrs = append(valueStrs, fmt.Sprintf("0x%04x", v))
	}
	return tupleStr(valueStrs)
}

// These test cases were produced using Python colormath:
// https://github.com/gtaylor/python-colormath.git
var xyzTestCases = []struct {
	r, g, b uint32  // input
	x, y, z float64 // expected output
}{
	{0x0000, 0x0000, 0x0000, 0.0000, 0.0000, 0.0000},
	{0x0000, 0x0000, 0x547a, 0.0161, 0.0064, 0.0846},
	{0x0000, 0x0000, 0xab84, 0.0733, 0.0293, 0.3863},
	{0x0000, 0x0000, 0xffff, 0.1805, 0.0722, 0.9504},
	{0x0000, 0x547a, 0x0000, 0.0318, 0.0636, 0.0106},
	{0x0000, 0x547a, 0x547a, 0.0479, 0.0701, 0.0952},
	{0x0000, 0x547a, 0xab84, 0.1052, 0.0930, 0.3969},
	{0x0000, 0x547a, 0xffff, 0.2123, 0.1358, 0.9610},
	{0x0000, 0xab84, 0x0000, 0.1453, 0.2907, 0.0484},
	{0x0000, 0xab84, 0x547a, 0.1614, 0.2971, 0.1330},
	{0x0000, 0xab84, 0xab84, 0.2187, 0.3200, 0.4348},
	{0x0000, 0xab84, 0xffff, 0.3258, 0.3629, 0.9989},
	{0x0000, 0xffff, 0x0000, 0.3576, 0.7152, 0.1192},
	{0x0000, 0xffff, 0x547a, 0.3736, 0.7216, 0.2038},
	{0x0000, 0xffff, 0xab84, 0.4309, 0.7445, 0.5055},
	{0x0000, 0xffff, 0xffff, 0.5380, 0.7873, 1.0696},
	{0x547a, 0x0000, 0x0000, 0.0367, 0.0189, 0.0017},
	{0x547a, 0x0000, 0x547a, 0.0528, 0.0253, 0.0863},
	{0x547a, 0x0000, 0xab84, 0.1100, 0.0483, 0.3880},
	{0x547a, 0x0000, 0xffff, 0.2172, 0.0911, 0.9522},
	{0x547a, 0x547a, 0x0000, 0.0685, 0.0826, 0.0123},
	{0x547a, 0x547a, 0x547a, 0.0846, 0.0890, 0.0969},
	{0x547a, 0x547a, 0xab84, 0.1419, 0.1119, 0.3986},
	{0x547a, 0x547a, 0xffff, 0.2490, 0.1547, 0.9628},
	{0x547a, 0xab84, 0x0000, 0.1820, 0.3096, 0.0502},
	{0x547a, 0xab84, 0x547a, 0.1981, 0.3160, 0.1347},
	{0x547a, 0xab84, 0xab84, 0.2554, 0.3389, 0.4365},
	{0x547a, 0xab84, 0xffff, 0.3625, 0.3818, 1.0006},
	{0x547a, 0xffff, 0x0000, 0.3943, 0.7341, 0.1209},
	{0x547a, 0xffff, 0x547a, 0.4103, 0.7405, 0.2055},
	{0x547a, 0xffff, 0xab84, 0.4676, 0.7634, 0.5072},
	{0x547a, 0xffff, 0xffff, 0.5747, 0.8063, 1.0714},
	{0xab84, 0x0000, 0x0000, 0.1676, 0.0864, 0.0079},
	{0xab84, 0x0000, 0x547a, 0.1837, 0.0929, 0.0924},
	{0xab84, 0x0000, 0xab84, 0.2410, 0.1158, 0.3942},
	{0xab84, 0x0000, 0xffff, 0.3481, 0.1586, 0.9583},
	{0xab84, 0x547a, 0x0000, 0.1994, 0.1501, 0.0185},
	{0xab84, 0x547a, 0x547a, 0.2155, 0.1565, 0.1030},
	{0xab84, 0x547a, 0xab84, 0.2728, 0.1794, 0.4048},
	{0xab84, 0x547a, 0xffff, 0.3799, 0.2223, 0.9689},
	{0xab84, 0xab84, 0x0000, 0.3130, 0.3771, 0.0563},
	{0xab84, 0xab84, 0x547a, 0.3290, 0.3835, 0.1409},
	{0xab84, 0xab84, 0xab84, 0.3863, 0.4064, 0.4426},
	{0xab84, 0xab84, 0xffff, 0.4934, 0.4493, 1.0067},
	{0xab84, 0xffff, 0x0000, 0.5252, 0.8016, 0.1271},
	{0xab84, 0xffff, 0x547a, 0.5413, 0.8080, 0.2116},
	{0xab84, 0xffff, 0xab84, 0.5986, 0.8309, 0.5134},
	{0xab84, 0xffff, 0xffff, 0.7057, 0.8738, 1.0775},
	{0xffff, 0x0000, 0x0000, 0.4124, 0.2127, 0.0193},
	{0xffff, 0x0000, 0x547a, 0.4285, 0.2191, 0.1039},
	{0xffff, 0x0000, 0xab84, 0.4858, 0.2420, 0.4056},
	{0xffff, 0x0000, 0xffff, 0.5929, 0.2848, 0.9698},
	{0xffff, 0x547a, 0x0000, 0.4442, 0.2763, 0.0299},
	{0xffff, 0x547a, 0x547a, 0.4603, 0.2827, 0.1145},
	{0xffff, 0x547a, 0xab84, 0.5176, 0.3056, 0.4162},
	{0xffff, 0x547a, 0xffff, 0.6247, 0.3485, 0.9804},
	{0xffff, 0xab84, 0x0000, 0.5578, 0.5033, 0.0678},
	{0xffff, 0xab84, 0x547a, 0.5738, 0.5098, 0.1524},
	{0xffff, 0xab84, 0xab84, 0.6311, 0.5327, 0.4541},
	{0xffff, 0xab84, 0xffff, 0.7382, 0.5755, 1.0182},
	{0xffff, 0xffff, 0x0000, 0.7700, 0.9278, 0.1385},
	{0xffff, 0xffff, 0x547a, 0.7861, 0.9342, 0.2231},
	{0xffff, 0xffff, 0xab84, 0.8434, 0.9572, 0.5248},
	{0xffff, 0xffff, 0xffff, 0.9505, 1.0000, 1.0890},
}

func TestRGBToXYZ(t *testing.T) {
	for _, tc := range xyzTestCases {
		x, y, z := rgbToCIEXYZ(tc.r, tc.g, tc.b)
		if !(eq(x, tc.x) && eq(y, tc.y) && eq(z, tc.z)) {
			t.Errorf(
				"rgbToCIEXYZ%v: want %v, got %v",
				hexTupleStr(tc.r, tc.g, tc.b),
				fTupleStr(tc.x, tc.y, tc.z),
				fTupleStr(x, y, z))
		}
	}
}

func eq_u(a, b uint32) bool {
	var diff uint32

	if a < b {
		diff = b - a
	} else {
		diff = a - b
	}

	result := diff <= 3

	if !result {
		fmt.Printf("FAIL eq_u(%d, %d); diff=%d\n", a, b, diff)
	}
	return result
}

func TestRGBToXYZToRGB(t *testing.T) {
	for _, tc := range xyzTestCases {
		x, y, z := rgbToCIEXYZ(tc.r, tc.g, tc.b)
		r, g, b := cieXYZToRGB(x, y, z)
		// If every component is within 0x0001 of wanted, that's close enough.
		if !(eq_u(r, tc.r) && eq_u(g, tc.g) && eq_u(b, tc.b)) {
			t.Errorf(
				"rgbxyz round trip %v: got %v",
				hexTupleStr(tc.r, tc.g, tc.b),
				hexTupleStr(r, g, b))
		}
	}
}
