int x, y[2];

int f(int x) {
  return x + y[0] + y[1];
}

int g(int *y) {
  return x + y[0] + y[1];
}

int main() {
  int z[2];

  x = 1;
  y[0] = 2;
  y[1] = 3;
  z[0] = 4;
  z[1] = 5;

  print(f(6) == 11 && g(z) == 10);
}
