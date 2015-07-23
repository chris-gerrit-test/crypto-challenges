#include <stdio.h>

#include <openssl/sha.h>

#include "crypt.c"
#include "encoding.c"

int main() {
  char mac[SHA_DIGEST_LENGTH];
  sha1_mac("christopher", 11, "hundt", 5, mac);
  print_bytes(mac, SHA_DIGEST_LENGTH);
}
