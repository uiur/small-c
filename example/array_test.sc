int main() {
  int data[4];
  int i;

  for (i = 0; i < 10; i = i + 1) {
    data[i] = i;
  }

  for (i = 0; i < 10; i = i + 1) {
    print(data[i]);
  }
}
