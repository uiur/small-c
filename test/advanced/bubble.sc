int a[10];

void swap(int *a, int*b) {
  int tmp;
  tmp = *a;
  *a = *b;
  *b = tmp;
}

int sum(int *a, int len) {
  int i, s;

  for(s = 0, i = 0; i < len; i = i + 1)
    s = s + a[i];

  return s;
}

int is_greater(int a, int b) {
  return a > b;
}

int sorted(int *a, int len) {
  if(len < 2)
    return 1;
  else {
    int i;

    for(i = 0; i < len - 1; i = i + 1)
      if(is_greater(a[i], a[i+1])) return 0;

    return 1;
  }
}

void sort(int *a, int len) {
  int i;
  for(i = len - 1; i > 0; i =  i - 1) {
    int j;
    for(j = 0; j < i; j = j + 1)
      if(is_greater(a[j], a[j+1]))
        swap(&a[j], &a[j+1]);
  }
}

void main() {
  int b[10];

  a[0] = b[0] = 8;
  a[1] = b[1] = 3;
  a[2] = b[2] = 1;
  a[3] = b[3] = 7;
  a[4] = b[4] = 4;
  a[5] = b[5] = 9;
  a[6] = b[6] = 10;
  a[7] = b[7] = 2;
  a[8] = b[8] = 6;
  a[9] = b[9] = 5;

  {
    int r;
    r = sorted(b, 10);
    sort(b, 10);
    print(sorted(a, 10) == 0 && r == 0 &&
          sorted(b, 10) && sum(a, 10) == 55
          && sum(b, 10) == 55);
  }
}
