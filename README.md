# Small C
[![Build Status](https://travis-ci.org/uiureo/small-c.svg?branch=master)](https://travis-ci.org/uiureo/small-c)

Small C compiler in Go. "Small C" is a small subset of C.

The target assembly language is MIPS.

This compiler is for [京都大学工学部情報学科計算機科学コース / 計算機科学実験及び演習3](http://www.fos.kuis.kyoto-u.ac.jp/~umatani/le3b/).

## Run

``` sh
make
./small-c example/quick_sort.sc
```

## Test
The test command requires [spim CLI](https://github.com/ymyzk/spim-for-kuis).

```sh
make test
```

## Example
This example shows bubble sort algorithm in Small C.

example/bubble_sort.sc:
``` c
void bubble_sort(int *p, int size) {
  int i, j, tmp;

  for (i = 0; i < size; i = i + 1) {
    for (j = 1; j < size; j = j + 1) {
      int current;
      int prev;

      current = *(p + j);
      prev = *(p + j - 1);
      if (current < prev) {
        tmp = current;
        *(p + j) = prev;
        *(p + j - 1) = tmp;
      }
    }
  }
}
```

There are other examples in `example/`

## License
MIT
