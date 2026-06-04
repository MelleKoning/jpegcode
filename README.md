# Jpeg in golang

We make use of `github.com/davidbyttow/govips/v2/vips`

so on ubuntu we need

```bash
apt install libvips-dev
```

## loading images in the browser

To test progressive jpg, serve the /data folder and set the browser developer tools to a slow connection to see if a jpeg is progressively loaded.

```bash
npx serve data
```

Usually that makes the data folder accessible at `http://localhost:3000`

in the browser go to dev tools, network and at "no throttling" select a slower connection, and load any of the example pages to see the difference between the progressive jpeg, which first loads a rough images and then slowly improves in the browser, and the non progressive jpeg that loads "line-by-line" top-down till finished.

## Ways to identify details of the image

install imagemagick, it comes with several tools.

```bash
identify -verbose ./data/lionheadbig.jpg
```

Shows that the original file has 4:4:4 full colour information, and keeping the 1x1, 1x1, 1x1 full information for the entire image (each pixel), with a JPEG quality factor of 100, meaning that JPEG rounding cant use the Discrete Cosine Transform (DCT) to shrunk the colour info.

You will notice that when you load the image and only reduce the JPEG quality from 100 to 95, this will reduce the filesize in half, even more then in half, while the human eye can't see any difference.

### Reasons of image reduction

- JPEG image quality
- Progressive loading

This small 5-point difference is actually massive. In JPEG compression math, Quality 100 is a trap. JPEG uses an algorithm called the Discrete Cosine Transform (DCT) to turn pixels into mathematical frequencies. It then uses a "Quantization Table" to divide these frequencies by a specific factor and round off the remainders. This rounding is where the file size shrinks.

When you select Quality 100, the software sets the quantization table divider to 1. This means virtually no mathematical rounding happens. The algorithm forces itself to store every microscopic bit of high-frequency noise (like sensor grain or tiny variations in the lion's fur), even if it is completely invisible to the human eye.

Dropping the quality down to 95 tells the algorithm it is allowed to start rounding numbers off. Because high-frequency noise takes up a huge chunk of data, filtering it out slashes the storage requirement in half, while leaving the actual image looking identical to your eyes.

While Quality 95 is doing the heavy lifting, two other changes helped push the size down:

#### Progressive Optimization

Original: Interlace: None (Baseline)
GoVips Output: Interlace: JPEG (Progressive)

Progressive JPEGs group similar frequencies together across the entire file. For a highly detailed image at full 4:4:4 color sampling, this grouping allows the final Huffman compression pass to compress the data stream much more tightly than a top-to-bottom baseline layout can. On top of that, for web page viewing, a progressive image loads the full height and weight of the image on screen as soon as possible and while the image is loading fills in the details, while Baseline will always load the image line by line top down. Counterintuitively, for images larger than about 10KB, progressive JPEGs are often 2% to 10% smaller than baseline JPEGs.
