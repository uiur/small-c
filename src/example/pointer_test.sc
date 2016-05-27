int main() {
  int a;
  int b;
  int *p;
  int data[2];
  a = 1;

  p = &a;
  *p = 0;
  b = *p;

  data[0] = 1;
  data[1] = 2;

  print(a == 0 && b == 0 && *(1 + data) == 2);
}
