int data[8];

int main() {
  int i;
  int j;
  int tmp;
  int size;

  size = 8;
  for (i = 0; i < size; i = i + 1) {
    for (j = 1; j < size; j = j + 1) {
      if (data[j] < data[j-1]) {
        tmp = data[j];
        data[j] = data[j-1];
        data[j-1] = tmp;
      }
    }
  }
}
