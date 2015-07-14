#include <openssl/aes.h>

void aes_decrypt(char *encrypted, char* decrypted, size_t num_bytes, char *key) {
    AES_KEY aes_key;
    AES_set_decrypt_key((unsigned char*)key, 128, &aes_key);
    for (size_t offset = 0; offset < num_bytes; offset += 16) {
        AES_decrypt((unsigned char*)encrypted + offset, (unsigned char*)decrypted + offset, &aes_key);
    }
}
