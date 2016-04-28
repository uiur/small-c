int main() {
  int a;
  int b;
  int *p;

  a = 1;

  p = &a;
  *p = 0;

  if (a > 0) {
    print(a);
  }

  b = *p;
  if (b > 0) {
    print(b);
  }
}
