void bubble_sort(int *p, int size) {
  int i, j, tmp;

  for (i = 0; i < size; i = i + 1) {
    for (j = 1; j < size; j = j + 1) {
      int current;
      int prev;

      current = *(p + j);
      prev = *(p + j - 1);
      if (current < prev) {
        tmp = current;
        *(p + j) = prev;
        *(p + j - 1) = tmp;
      }
    }
  }
}

int main() {
  int data[8];
  int size;
  int i;

  size = 8;

  data[0] = 4;
  data[1] = 2;
  data[2] = 1;
  data[3] = 3;
  data[4] = 6;
  data[5] = 8;
  data[6] = 7;
  data[7] = 5;

  bubble_sort(data, size);

  for (i = 0; i < size; i = i + 1) {
    print(data[i]);
  }
}
