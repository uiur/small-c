int  val[201];

void init_val() {
  int i;

  for(i = 0; i <= 200; i = i + 1)
    val[i] = 200 - i;
}

int sum(int *arr) {
  int x;

  if((x = *arr) != 0)
    return x + sum(arr + 1);
  else
    return 0;
}

int sum2(int *arr, int acc) {
  int x;

  if((x = *arr) != 0)
    return sum2(arr + 1, acc + x);
  else
    return acc;
}

void main() {
  init_val();
  print(sum(val) == 20100 &&
        sum2(val, 0) == 20100);
}
