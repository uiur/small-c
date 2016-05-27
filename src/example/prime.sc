int prime[10];

int mod(int x, int y) {
  return x - (x / y) * y;
}

int main() {
  int x, i, j, k, N;

  N = 10;
  x = 1;
  k = 1;
  prime[0] = 2;

  while (k < N) {
    int m;

    x = x + 2;
    j = 0;

    while (j < k && mod(x, prime[j]) != 0) {
      j = j + 1;
    }

    if (j == k) {
      prime[k] = x;
      k = k + 1;
    }
  }

  for (i = 0; i < N; i = i + 1) {
    print(prime[i]);
    putchar(' ');
  }
}
