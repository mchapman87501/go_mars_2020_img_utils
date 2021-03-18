#!/usr/bin/env python
# coding: utf-8

# In[1]:


# Trying to gin up a human-readable, simple-minded (bilinear interpolation) algorithm for de-mosaicing a
# sensor readout that has an RGGB color filter array (CFA).

# Red filters lie over cells whose x coordinate is even and whose y coordinate is even:  even, even
# Blue filters: odd, odd
# Green filters: even, odd *and* odd, even.


# In[2]:


import numpy as np
from PIL import Image


# In[13]:


# Image dimensions
width = 255
height = 255

# Dummy image data is grayscale - single component, 0..255.
# Build it up as a gradient.
# Give it a demosaiced red tinge by boosting pixels that should be
# under a red filter in the Bayer image pattern.
dummy_image_data = []
for y in range(height):
    row = []
    for x in range(width):
        red_boost = 100 if (x % 2, y % 2) == (0, 0) else 0
        row.append(min(255, x + red_boost))
    dummy_image_data.append(row)


gray_image_data = np.array(dummy_image_data, dtype=np.uint8)
print("Dummy image data:", gray_image_data)
# PIL seems to be ignoring my mode, dangit.
gray_img = Image.fromarray(gray_image_data, mode="L")
gray_img.show()

print("Converted back to numpy array:")
print(np.asarray(gray_img))


# In[14]:


# Offset of each color component within a pixel:
R = 0
G = 1
B = 2

# filter pattern, addressable as [y][x]
pattern = [
    [R, G],
    [G, B]
]

# Demosaiced image data is RGB - three components.
demosaiced = []
for y in range(height):
    row = [[0, 0, 0] for x in range(width)]
    demosaiced.append(row)


def indices(v, limit):
    result = []
    for offset in [-1, 0, 1]:
        index = v + offset
        if 0 <= index < limit:
            result.append(index)
    return result


def channel(x, y):
    x_pattern = x % 2
    y_pattern = y % 2
    return pattern[y_pattern][x_pattern]


def demosaic(sensor_image, demosaiced, width, height):
    for x_image in range(width):
        x_indices = indices(x_image, width)
        for y_image in range(height):
            y_indices = indices(y_image, height)

            sums = {R: 0, G: 0, B: 0}
            counts = {R: 0, G: 0, B: 0}

            for x in x_indices:
                for y in y_indices:
                    c = channel(x, y)
                    sums[c] += sensor_image[y][x]
                    counts[c] += 1
            for c in [R, G, B]:
                intensity = sums[c] / counts[c] if counts[c] > 0 else 0
                # May as well convert to 8-bit integer.
                pixel_value = min(255, max(0, int(intensity)))
                demosaiced[y_image][x_image][c] = pixel_value


demosaic(dummy_image_data, demosaiced, width, height)


# In[15]:


color_img = Image.fromarray(np.array(demosaiced, dtype=np.uint8), mode="RGB")
color_img.show()
