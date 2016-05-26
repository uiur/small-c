int loop() {
  loop();
  return -1;
}

int check_and() {
  if(0 && loop()) return 0;
  return 1;
}

int check_or() {
  if(1 || loop()) return 1;
  return 0;
}

void main() {
  print(check_and() && check_or());
}
