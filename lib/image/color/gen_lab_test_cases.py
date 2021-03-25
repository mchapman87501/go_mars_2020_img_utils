#!/usr/bin/env python3

from colormath.color_objects import sRGBColor, LabColor
from colormath.color_conversions import convert_color


comp_values = [0.0, 0.33, 0.67, 1.0]


def to_uint16(fval):
    return int(fval * 0xffff)

def to_float(uval):
    return uval / 0xffff


def print_rgb_lab_test_cases():
    print("// sRGB to CIE Lab")
    for r in comp_values:
        r16 = to_uint16(r)
        for g in comp_values:
            g16 = to_uint16(g)
            for b in comp_values:
                b16 = to_uint16(b)
                rgb = sRGBColor(to_float(r16), to_float(g16), to_float(b16))
                lab = convert_color(rgb, LabColor)
                x = lab.lab_l
                y = lab.lab_a
                z = lab.lab_b
                print(f"        {{"
                      f"0x{r16:04x}, 0x{g16:04x}, 0x{b16:04x}, "
                      f"{x:.6f}, {y:.6f}, {z:.6f}"
                      f"}},")

                # Round-trip check.
                # rgb_rt = convert_color(lab, sRGBColor)
                # rgb_rt16 = [f"0x{to_uint16(v):04x}" for v in [rgb_rt.rgb_r, rgb_rt.rgb_g, rgb_rt.rgb_b]]
                # print(f"          // {rgb}->{lab}->{rgb_rt}")
    print("")


def main():
    print_rgb_lab_test_cases()


if __name__ == "__main__":
    main()