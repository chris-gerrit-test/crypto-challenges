#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

int main() {
    char key[16] = "YELLOW SUBMARINE";
    char nonce[16] = {0};

    char s[] = "L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==";
    char *enc = calloc(strlen(s), 1);
    int n = base64_to_bytes(s, enc, strlen(s));
    assert(n > 0);
    enc[n] = '\0';
    ctr_crypt(enc, enc, n, nonce, key);
    printf("%s\n", enc);
}
