#!/usr/bin/env python3

from colormath.color_objects import sRGBColor, LabColor, XYZColor
from colormath.color_conversions import convert_color


comp_values = [0.0, 0.33, 0.67, 1.0]


def clip(val, min_val=0.0, max_val=1.0):
    return min(max_val, max(min_val, val))


def to_uint16(fval):
    return int(fval * 0xffff)


def to_float(uval):
    return uval / 0xffff


_rl_preamble = """var rgbToLabTestCases = []struct {
    r, g, b          uint32
    labL, laba, labb float64
}{"""


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
                print(f"    {{"
                      f"0x{r16:04x}, 0x{g16:04x}, 0x{b16:04x}, "
                      f"{x:.6f}, {y:.6f}, {z:.6f}"
                      f"}},")
    print("}")


_lr_preamble = """var labToRGBTestCases = []struct {
    labL, laba, labb float64
    r, g, b          uint32
}{"""


def in_gamut(rgb):
    components = (rgb.rgb_r, rgb.rgb_g, rgb.rgb_b)
    return all((0.0 <= comp <= 1.0) for comp in components)


def print_lab_rgb_test_cases():
    print(_lr_preamble)
    for labl in [0.0, 33.0, 67.0, 100.0]:
        for laba in [-128.0, -31.5, 0.0, 64.0, 127.0]:
            for labb in [-128.0, -64.0, 0.0, 64.0, 127.0]:
                # Colormath Lab colors default to D50.
                # For purposes of this Go library, use standard illuminante
                # D65 throughout.
                lab = LabColor(labl, laba, labb, illuminant="d65")
                xyz = convert_color(lab, XYZColor)
                rgbn = convert_color(xyz, sRGBColor)
                # CIE Lab covers a larger gamut than sRGB.  Omit points
                # can't be expressed in sRGB.
                if in_gamut(rgbn):
                    rn, gn, bn = rgbn.rgb_r, rgbn.rgb_g, rgbn.rgb_b
                    r16 = to_uint16(clip(rn))
                    g16 = to_uint16(clip(gn))
                    b16 = to_uint16(clip(bn))
                    print(f"    {{"
                          f"{labl:.2f}, {laba:.2f}, {labb:.2f}, "
                          f"0x{r16:04x}, 0x{g16:04x}, 0x{b16:04x}"
                          f"}},")
                # else:
                #     print(f"   // Out of gamut: {lab}")
    print("}")


def main():
    # print_rgb_lab_test_cases()
    print_lab_rgb_test_cases()


if __name__ == "__main__":
    main()