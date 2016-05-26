int cmp1(int a, int b) {
  return a + 1 > b;
}

int cmp2(int a, int b) {
  return a >= b + 1;
}

int cmp3(int a, int b) {
  return a > b + 1;
}

void main() {
  print(cmp1(2, 3) == 0 &&
        cmp2(4, 3) == 1 &&
        cmp3(4, 3) == 0);
}
