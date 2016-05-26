int w(int x) {
  int r;
  r = 0;

  while (x > 0) {
    r = r + x;
    x = x - 1;
  }

  return r;
}

int main() {
  print(w(10) == 55);
}
