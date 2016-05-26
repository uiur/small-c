void swap(int *i, int*j) {
  int tmp;
  tmp = *i;
  *i = *j;
  *j = tmp;
}

void main() {
  int x, y;
  x = 1;
  y = 2;

  swap(&x, &y);
  print(x == 2 && y ==1);
}
