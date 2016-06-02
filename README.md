# Small C
[![Build Status](https://travis-ci.org/uiureo/small-c.svg?branch=master)](https://travis-ci.org/uiureo/small-c) [![Coverage Status](https://coveralls.io/repos/github/uiureo/small-c/badge.svg?branch=master)](https://coveralls.io/github/uiureo/small-c?branch=master)

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

example/fizzbuzz.sc:
```c
int mod(int x, int y) {
  return x - (x / y) * y;
}

void puts(int *s) {
  while (*s != 0) {
    putchar(*s);
    s = s + 1;
  }
}

int main() {
  int i;
  int fizz[5];
  int buzz[5];

  fizz[0] = 'F'; fizz[1] = 'i'; fizz[2] = 'z'; fizz[3] = 'z'; fizz[4] = 0;
  buzz[0] = 'B'; buzz[1] = 'u'; buzz[2] = 'z'; buzz[3] = 'z'; buzz[4] = 0;

  for (i = 1; i <= 30; i = i + 1) {
    if (mod(i, 3) == 0) {
      puts(fizz);
    }

    if (mod(i, 5) == 0) {
      puts(buzz);
    }

    if (mod(i, 3) != 0 && mod(i, 5) != 0) {
      print(i);
    }

    putchar(' ');
  }
}

```

There are other examples in `src/example/`

## License
MIT
