# IDEAs

* CLI
* Some way to monitor the blink on terminal
* Draw circle if requested
* Stats? Graph? Time?

* Create folder with start timestamp

* The size of the world can be smaller than 4K, we render it in high quality
* Param for scale, check if the requested resolution fits
* Use labels for common res

* First you scale up then you quantize the position.

* `elemColor/templateFirefly` should be in a separate file `const.go`

```
  3.840 = 2^8 × 3 × 5    ;   2.160 = 2^4 × 3^3 × 5

( 3.840 = 120 * 32 ) * 1 ; ( 2.160 = 120 * 18 ) * 1
( 3.840 = 240 * 16 ) * 1 ; ( 2.160 = 240 *  9 ) * 1
( 1.920 = 120 * 16 ) * 2 ; ( 1.080 = 120 *  9 ) * 2
( 1.280 =  80 * 16 ) * 3 ; (   720 =  80 *  9 ) * 3
(   960 =  60 * 16 ) * 4 ; (   540 = 120 *  9 ) * 4
(   768 =  48 * 16 ) * 5 ; (   432 =  48 *  9 ) * 5
(   640 =  40 * 16 ) * 6 ; (   360 =  40 *  9 ) * 6
    ^      ^           ^
    |      |           ╵- scale factor
    |      ╵- cell size
    ╵- world size
```

* bmp of the fireflies in several directions
* Use HCL https://github.com/lucasb-eyer/go-colorful
* Save min and max L value for each pixel

* We are using sane colorspaces _specifically_ to have uniform luminosity bro

### Image to video

* https://stackoverflow.com/questions/46397240/ffmpeg-image2pipe-producing-broken-video
* https://github.com/leixiaohua1020/simplest_ffmpeg_video_encoder


# TODOs

1. Max/min HCL map
2. Cache all templates
3. Implement blitting of rotated/shifted templates on larger image


* Create a grid with the colors blended (https://github.com/lucasb-eyer/go-colorful#blending-colors)
* or a preview of the glowing fireflies
* method GetBlent(t) that returns the blended color?
