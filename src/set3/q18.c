#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"

int main() {
    char key[16] = "YELLOW SUBMARINE";
    char nonce[8] = {0};

    char s[] = "L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==";
    int n = base64_to_bytes(s, s, strlen(s));
    assert(n > 0);
    s[n] = '\0';
    ctr_crypt(s, s, n, nonce, key);
    printf("%s\n", s);
}
