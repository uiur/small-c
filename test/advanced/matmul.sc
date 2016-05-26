int n;

int get(int *a, int i, int j) {
  return a[i * n + j];
}

void set(int *a, int i, int j, int val) {
  a[i * n + j] = val;
}

void matmul(int *a, int *b, int *c) {
  int i, j, k;

  for(i = 0; i < n; i = i + 1)
    for(j = 0; j < n; j = j + 1) {
      int tmp;
      tmp = 0;
      for(k = 0; k < n; k = k + 1)
        tmp = tmp + get(a, i, k) * get(b, k, j);
      set(c, i, j, tmp);
    }
}

void main() {
  int a[9], b[9], c[9];
  
  n = 3;

  set(a, 0, 0,  2);
  set(a, 0, 1,  3);
  set(a, 0, 2,  2);
  set(a, 1, 0,  1);
  set(a, 1, 1,  4);
  set(a, 1, 2, -1);
  set(a, 2, 0, -2);
  set(a, 2, 1,  1);
  set(a, 2, 2, -3);

  set(b, 0, 0, -3);
  set(b, 0, 1,  1);
  set(b, 0, 2,  2);
  set(b, 1, 0, -2);
  set(b, 1, 1, -4);
  set(b, 1, 2,  2);
  set(b, 2, 0,  4);
  set(b, 2, 1,  3);
  set(b, 2, 2,  1);

  matmul(a, b, c);

  print(get(c, 0, 0) ==  -4 &&
        get(c, 0, 1) ==  -4 &&
        get(c, 0, 2) ==  12 &&
        get(c, 1, 0) == -15 &&
        get(c, 1, 1) == -18 &&
        get(c, 1, 2) ==   9 &&
        get(c, 2, 0) ==  -8 &&
        get(c, 2, 1) == -15 &&
        get(c, 2, 2) ==  -5);
}
