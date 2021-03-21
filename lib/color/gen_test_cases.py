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


def main():
    print_hsv_to_rgb_cases()


if __name__ == "__main__":
    main()
