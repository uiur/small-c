int ack(int x, int y) {
  if(x == 0)
    return y+1;
  if(y == 0)
    return ack(x-1, 1);
  return ack(x-1, ack(x, y-1));
}

void main() {
  print(ack(3, 3) ==  61 &&
        ack(3, 4) == 125 &&
        ack(3, 5) == 253);
}
