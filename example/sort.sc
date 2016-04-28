int main() {
  int data[4];
  int i;
  int j;
  int tmp;
  int size;

  data[0] = 4;
  data[1] = 2;
  data[2] = 1;
  data[3] = 3;

  size = 4;

  for (i = 0; i < size; i = i + 1) {
    for (j = 1; j < size; j = j + 1) {
      if (data[j] < data[j-1]) {
        tmp = data[j];
        data[j] = data[j-1];
        data[j-1] = tmp;
      }
    }
  }

  for (i = 0; i < size; i = i + 1) {
    print(data[i]);
  }
}
