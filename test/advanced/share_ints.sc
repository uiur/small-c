int init_v(int *a) {
  int i;

  for(i = 0; i < 10; i = i + 1) 
    a[i] = i + 1;

  return i;
}

void main() {
  int i, j, size, end, sum, sum2, v[16], *s[16], *t[16];

  size = init_v(v);

  for(i = 0; i < size; i = i + 1)
    s[i] = &v[i];
  s[i] = &end;

  for(i = 0, sum = 0; s[i] != &end; i = i + 1)
    sum = sum + *s[i];

  for(i = 0, j = 0; i < size; i = i + 2, j = j + 1)
    t[j] = s[i];
  t[j] = &end;

  for(j = 0; t[j] != &end; j = j + 1)
    *t[j] = *t[j] * *t[j];

  for(i = 0, sum2 = 0; s[i] != &end; i = i + 1)
    sum2 = sum2 + *s[i];

  print(sum == 55 && sum2 == 195);
}
