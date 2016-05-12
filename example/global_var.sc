int a, data[10];

int main() {
  int *p;
  a = 42;

  data[0] = 1;
  data[1] = 2;
  data[2] = 3;

  p = &a;
  print(*p == 42);
  print(data[0] + data[2] == 4);
}
