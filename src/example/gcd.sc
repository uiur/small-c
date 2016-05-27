// greatest common divisor by euclidean algorithm
int mod(int x, int y) {
  return x - (x / y) * y;
}

int gcd(int x, int y) {
  if (y == 0) {
    return x;
  } else {
    return gcd(y, mod(x, y));
  }
}

int main() {
  print(gcd(1071, 1029)); // 21
}
