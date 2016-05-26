int s() {
  int x;
  x = 1;
  {
    int y;
    y = 2;
    {
      int x;
      x = 3;
      y = y + x;
    }
    x = x + y;
  }
  {
    int y;
    y = 4;
    return x + y;
  }
}

int main() {
  print(s() == 10);
}
