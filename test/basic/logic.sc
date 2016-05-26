int logi1(int a, int b, int c, int d) {
  if(a < b+1 && b+1 < c || b == 0)
    return 1;
  else
    return 0;
}

int logi2(int a, int b, int c, int d) {
  if(a < b && b < c || d)
    return 1;
  else
    return 0;
}

int logi3(int x, int y) {
  if (y == 3 && y - x) {
    if (x == 0 || x+y < x*y)
      return 1;
    else
      return 0;
  } else
    return 0;
}

int main() {
  print(logi1(1, 2, 4, 3) &&
        logi2(1, 2, 3, 0) &&
        logi3(2, 3));
}
