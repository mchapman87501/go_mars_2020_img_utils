#!/usr/bin/env python3

from colormath.color_objects import sRGBColor, XYZColor
from colormath.color_conversions import convert_color


comp_values = [0.0, 0.33, 0.67, 1.0]


def to_uint16(fval):
    return int(fval * 0xffff)


def print_rgb_xyz_test_cases():
    print("// sRGB to CIE XYZ")
    for r in comp_values:
        r16 = to_uint16(r)
        for g in comp_values:
            g16 = to_uint16(g)
            for b in comp_values:
                b16 = to_uint16(b)
                rgb = sRGBColor(r, g, b)
                xyz = convert_color(rgb, XYZColor)
                x = xyz.xyz_x
                y = xyz.xyz_y
                z = xyz.xyz_z
                print(f"        {{"
                      f"0x{r16:04x}, 0x{g16:04x}, 0x{b16:04x}, "
                      f"{x:.4f}, {y:.4f}, {z:.4f}"
                      f"}},")

                # Round-trip check.
                rgb_rt = convert_color(xyz, sRGBColor)
                rgb_rt16 = [f"0x{to_uint16(v):04x}" for v in [rgb_rt.rgb_r, rgb_rt.rgb_g, rgb_rt.rgb_b]]
                print(f"          // {rgb}->{xyz}->{rgb_rt}")
                # print(f"        // ={rgb_rt16}")
                # print("")
    print("")


def main():
    print_rgb_xyz_test_cases()


if __name__ == "__main__":
    main()