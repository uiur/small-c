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

int l[6], r[6];
int bound;

void merge(int *a, int len) {
  int i;

  for(i = 0; i < len/2; l[i] = a[i], i = i + 1) ;
  l[i] = bound;

  for(i = len/2; i < len; r[i - len/2] = a[i], i = i + 1) ;
  r[i - len/2] = bound;

  {
    int j, k;
    for(i = 0, j = 0, k = 0; i < len; i = i + 1)
      if(is_greater(l[j], r[k])) {
        a[i] = r[k];
        k = k + 1;
      } else {
        a[i] = l[j];
        j = j + 1;
      }
  }
}

void sort(int *a, int len) {
  if(len < 2) return;
  sort(a, len/2);
  sort(a + len/2, len - len/2);
  merge(a, len);
}

void main() {
  int a[10];

  bound = 11;
  
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
    sort(a, 10);
    print(r == 0 && sorted(a, 10) && sum(a, 10) == 55);
  }
}
