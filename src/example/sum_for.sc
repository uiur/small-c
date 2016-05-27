int main() {
  int sum;
  int i;

  sum = 0;

  for (i = 0; i < 10; i = i + 1) {
    sum = sum + i;
  }

  print(sum);
}
