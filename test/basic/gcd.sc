int gcd(int i, int j);

int main() {
  print(gcd(315, 189) == 63);
}

int gcd(int a, int b) {
  if(a == b)
    return a;
  else if(a > b)
    return gcd(a-b, b);
  else
    return gcd(a, b-a);
}

