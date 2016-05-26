int is_odd(int n);

int is_even(int n) {
  if(n == 0) return 1;
  if(n == 1) return 0;
  return is_odd(n-1);
}

int is_odd(int n) {
  if(n == 0) return 0;
  if(n == 1) return 1;
  return is_even(n-1);
}

int all_even(int *a) {
  if(*a < 0)
    return 1;
  else 
    return is_even(*a) && all_even(a+1);
}

int end;

int *find_odd(int *a) {
  for(;;) {
    if(*a < 0) return &end;
    if(is_odd(*a)) return a;
    a = a + 1;
  }
}

void odd_to_even(int *a) {
  int *p;
  while((p = find_odd(a)) != &end)
    *p = *p + 1;
}

void init(int *a) {
  int i;
  for(i = 1; i <= 10; i = i + 1, a = a + 1)
    *a = i;
  *a = -1;
}

void main() {
  int a[11], r;

  init(a);
  r = all_even(a);
  odd_to_even(a);
  print(r == 0 && all_even(a));
}
