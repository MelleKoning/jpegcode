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
