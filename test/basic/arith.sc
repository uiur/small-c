int ff(int x, int y, int z) {
  return x + y + z;
}

int gg(int x) {
  return x * ff(1, 2, 3);
}

int hh(int x) {
  return x / ff(1, 2, 3);
}

int ii() {
  return ff(1, 2, 3) - gg(4);
}

int jj(int x, int y) {
  int i;

  i = 10;
  x = x - y;
  i = i - x - 1;

  return x + i;
}

void main() {
  print(ff(1, 2, 3) ==   6 &&
        gg(10)      ==  60 &&
        hh(40)      ==   6 &&
        ii()        == -18 &&
        jj(2, 4)    ==   9);
}
