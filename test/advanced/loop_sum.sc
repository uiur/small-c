void init_v(int *v) {
  int i;

  for(i = 0; i <= 100; i = i + 1)
    v[i] = 100 - i;
}

int sum(int *arr) {
  int s, x;

  s = 0;
  while(x = *arr) {
    s = s + x;
    arr = arr + 1;
  }

  return s;
}

void main() {
  int  v[101];
  init_v(v);
  print(sum(v) == 5050);
}
