void swap(int *p, int *q);

void quick_sort(int *p, int left, int right) {
  int i, j, pivot;

  if (left >= right) {
    return;
  }

  pivot = *(p + left + (right - left) / 2);

  i = left;
  j = right;

  while (i < j) {
    while (*(p + i) < pivot) {
      i = i + 1;
    }

    while (pivot < *(p + j)) {
      j = j - 1;
    }

    if (i < j) {
      swap(p + i, p + j);

      i = i + 1;
      j = j - 1;
    }
  }

  quick_sort(p, left, i - 1);
  quick_sort(p, j + 1, right);
}

void swap(int *p, int *q) {
  int tmp;

  tmp = *p;
  *p = *q;
  *q = tmp;
}

int main() {
  int i, size, data[8];
  size = 8;

  for (i = 0; i < size; i = i + 1) {
    data[i] = size - i;
  }

  swap(data, data + 4);
  swap(data + 1, data + 5);
  swap(data + 3, data + 7);

  quick_sort(data, 0, size);

  for (i = 0; i < size; i = i + 1) {
    print(data[i]);
  }
}
