# TODOs

* CLI
* Some way to monitor the blink on terminal
* Draw circle if requested
* Stats? Graph? Time?

* The size of the world can be smaller than 4K, we render it in high quality
* Param for scale, check if the requested resolution fits
* Use labels for common res

```
  3.840 = 2^8 × 3 × 5    ;   2.160 = 2^4 × 3^3 × 5

( 3.840 = 120 * 32 ) * 1 ; ( 2.160 = 120 * 1 ) * 18
( 3.840 = 240 * 16 ) * 1 ; ( 2.160 = 240 * 1 ) *  9
( 1.920 = 120 * 16 ) * 2 ; ( 1.080 = 120 * 2 ) *  9
( 1.280 =  80 * 16 ) * 3 ; (   720 =  80 * 3 ) *  9
(   960 =  60 * 16 ) * 4 ; (   540 = 120 * 4 ) *  9
(   768 =  48 * 16 ) * 5 ; (   432 =  48 * 5 ) *  9
(   640 =  40 * 16 ) * 6 ; (   360 =  40 * 6 ) *  9
    ^      ^           ^
    |      |           |- scale factor
    |      |-    cell size
    |- world size
```
