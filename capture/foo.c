int main() {
  int i;
  for (i = 0; i < 6; i++) {
    printf("pixel %i y:%d cb:%d cr:%d\n",
           i, i * 2,
           (i & 0xfffffffe) * 2 + 1,
           (i & 0xfffffffe) * 2 + 3);
  }
}
