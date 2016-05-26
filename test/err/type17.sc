int f(int *a, int b) {
  return *a + b; 
}

void main() {
  int a;
  f(&a);
}