#include <stdio.h>

#include "crypt.c"
#include "encoding.c"
#include "math.c"
#include "score.c"

int main() {
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
    aes_decrypt(as_bytes, decrypted, num_bytes, "YELLOW SUBMARINE");
    printf("%s\n", decrypted);
}
