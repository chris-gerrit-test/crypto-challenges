#include <assert.h>
#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "math.c"
#include "score.c"

char key[16];
char nonce[8];

void edit(char *encrypted, size_t len, size_t offset, char *new_text) {
    size_t n = strlen(new_text);
    assert(offset + n <= len);
    ctr_crypt(encrypted, encrypted, len, nonce, key);
    memcpy(encrypted + offset, new_text, n);
    ctr_crypt(encrypted, encrypted, len, nonce, key);
}

void init_keys() {
    for (size_t i = 0; i < sizeof(key); ++i) {
        key[i] = randn(256) - 128;
    }
    for (size_t i = 0; i < sizeof(nonce); ++i) {
        nonce[i] = randn(256) - 128;
    }
}

int main() {
    srand(44);
    init_keys();

    char *doc = calloc(1, 1);
    size_t doc_len = 0;
    char *line = NULL;
    size_t n = 0;
    while (getline(&line, &n, stdin) != -1) {
    	size_t line_len = strlen(line);
    	doc = realloc(doc, doc_len + line_len);
    	strncpy(doc + doc_len, line, line_len - 1);
    	doc_len += line_len - 1;
    	doc[doc_len] = '\0';
    }
    size_t num_bytes = num_bytes_from_base64(doc_len);
    char *as_bytes = calloc(num_bytes + 1, 1);
    num_bytes = base64_to_bytes(doc, as_bytes, num_bytes);

    char *decrypted = calloc(num_bytes + 1, 1);
    char *encrypted = calloc(num_bytes + 1, 1);
    aes_decrypt(as_bytes, decrypted, num_bytes, "YELLOW SUBMARINE");

    ctr_crypt(decrypted, encrypted, num_bytes, nonce, key);  //  E = C = P ^ K

    char *a = calloc(num_bytes + 1, 1);
    memset(a, 'a', num_bytes);
    memcpy(decrypted, encrypted, num_bytes);
    edit(encrypted, num_bytes, 0, a);  // D = A ^ K; 
    xor(decrypted, encrypted, num_bytes, decrypted);  // D = (A ^ K) ^ (P ^ K) = A ^ P
    xor(decrypted, a, num_bytes, decrypted);  // D = (A ^ P) ^ A = P

    printf("%s\n", decrypted);

    free(doc);
    free(as_bytes);
    free(encrypted);
    free(decrypted);
    free(a);
}
