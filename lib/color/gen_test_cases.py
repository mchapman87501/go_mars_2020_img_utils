#!/usr/bin/env python3

import colorsys as cs


def print_hsv_to_rgb_cases():
    print("# HSV to RGB:")
    values = [0.0, 0.25, 0.5, 0.75, 1.0]
    for h in values:
        for s in values:
            for v in values:
                rgb_f = cs.hsv_to_rgb(h, s, v)
                r, g, b = [int(f * 0xffffffff) for f in rgb_f]
                print(f"    {{{h:.3f}, {s:.3f}, {v:.3f},  "
                      f"0x{r:x}, 0x{g:x}, 0x{b:x}}},")


def print_rgb_to_hsv_cases():
    print("# RGB to HSV:")
    values = list(range(0, 256, 51))
    for r in values:
        for g in values:
            for b in values:
                rn = r / 0xff
                gn = g / 0xff
                bn = b / 0xff
                h, s, v = cs.rgb_to_hsv(rn, gn, bn)
                r32, g32, b32 = [int(f * 0xffffffff) for f in [rn, gn, bn]]
                print(f"    {{0x{r32:x}, 0x{g32:x}, 0x{b32:x},  "
                      f"{h:.3f}, {s:.3f}, {v:.3f}}},")


def main():
    # print_hsv_to_rgb_cases()
    print_rgb_to_hsv_cases()


if __name__ == "__main__":
    main()
