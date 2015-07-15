#include <openssl/aes.h>

int pkcs7(char *buf, size_t buf_size, size_t block_size) {
	for (size_t i = 0; i < buf_size; ++i) {
		if (buf[i] == '\0') {
			size_t padding_size = block_size - (i % block_size);
			if (i + padding_size + 1 < buf_size) {
				return -1;
			}
			for (size_t j = 0; j < padding_size; ++j) {
				buf[i + j] = padding_size;
			}
			buf[i + padding_size] = '\0';
			return i + padding_size;
		}
	}
	return -1;
}

void aes_decrypt(char *encrypted, char* decrypted, size_t num_bytes, char *key) {
    AES_KEY aes_key;
    AES_set_decrypt_key((unsigned char*)key, 128, &aes_key);
    for (size_t offset = 0; offset < num_bytes; offset += 16) {
        AES_decrypt((unsigned char*)encrypted + offset, (unsigned char*)decrypted + offset, &aes_key);
    }
}
