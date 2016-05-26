int mod(int x, int y) {
  return x - (x / y) * y;
}

void puts(int *s) {
  while (*s != 0) {
    putchar(*s);
    s = s + 1;
  }
}

int main() {
  int i;
  int fizz[5];
  int buzz[5];

  fizz[0] = 'F'; fizz[1] = 'i'; fizz[2] = 'z'; fizz[3] = 'z'; fizz[4] = 0;
  buzz[0] = 'B'; buzz[1] = 'u'; buzz[2] = 'z'; buzz[3] = 'z'; buzz[4] = 0;

  for (i = 1; i <= 30; i = i + 1) {
    if (mod(i, 3) == 0) {
      puts(fizz);
    }

    if (mod(i, 5) == 0) {
      puts(buzz);
    }

    if (mod(i, 3) != 0 && mod(i, 5) != 0) {
      print(i);
    }

    putchar(' ');
  }
}
