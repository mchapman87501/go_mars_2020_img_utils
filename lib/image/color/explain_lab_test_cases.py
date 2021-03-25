#!/usr/bin/env python3

from colormath.color_objects import sRGBColor, LabColor, XYZColor
from colormath.color_conversions import convert_color


comp_values = [0.0, 0.33, 0.67, 1.0]


def to_uint16(fval):
    return int(fval * 0xffff)


def to_float(uval):
    return uval / 0xffff


def explain_lab_rgb():
    labl_vals = [0.0]  # [0.0, 33.0, 67.0, 100.0]
    laba_vals = [-128.0]  # [-128.0, -31.5, 0.0, 64.0, 127.0]
    labb_vals = [-128.0]  # [-128.0, -64.0, 0.0, 64.0, 127.0]

    for labl in labl_vals:
        for laba in laba_vals:
            for labb in labb_vals:
                print(f"Lab: {labl:.1f}, {laba:.1f}, {labb:.1f}")
                # Colormath Lab colors default to D50.
                # For purposes of this Go library, use standard illuminante
                # D65 throughout.
                lab = LabColor(labl, laba, labb, illuminant="d65")

                # Why is colormath using D65 in this direction?
                xyz = convert_color(lab, XYZColor)
                print(f"XYZ: {xyz.xyz_x:.4f}, {xyz.xyz_y:.4f}, {xyz.xyz_x:.4f}")

                rgb = convert_color(xyz, sRGBColor)
                r = to_uint16(rgb.rgb_r)
                g = to_uint16(rgb.rgb_g)
                b = to_uint16(rgb.rgb_b)
                print(f"RGB: 0x{r:04x}, 0x{g:04x}, 0x{b:04x}")
                print("")


def main():
    explain_lab_rgb()


if __name__ == "__main__":
    main()