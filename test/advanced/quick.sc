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

int sum(int *a, int len) {
  int i, s;

  for(s = 0, i = 0; i < len; i = i + 1)
    s = s + a[i];

  return s;
}

void swap(int *a, int*b) {
  int tmp;
  tmp = *a;
  *a = *b;
  *b = tmp;
}

void sort(int *a, int left, int right) {
  int i, last;

  if(left >= right) return;
  swap(&a[left], &a[(left + right)/2]);
  last = left;
  i = left + 1;
  while (i <= right) {
    if(is_greater(a[left], a[i])) {
      last = last + 1;
      swap(&a[last], &a[i]);
    }
    i = i + 1;
  }
  swap(&a[left], &a[last]);
  sort(a, left, last - 1);
  sort(a, last + 1, right);
}

int main() {
  int a[10];

  a[0] = 8;
  a[1] = 3;
  a[2] = 1;
  a[3] = 7;
  a[4] = 4;
  a[5] = 9;
  a[6] = 10;
  a[7] = 2;
  a[8] = 6;
  a[9] = 5;

  {
    int r;
    r = sorted(a, 10);
    sort(a, 0, 9);

    print(r == 0 && sorted(a, 10) && sum(a, 10) == 55);
  }
}
