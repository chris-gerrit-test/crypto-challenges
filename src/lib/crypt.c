#include <string.h>
#include <openssl/aes.h>

#include "math.c"

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

char *strip_pkcs7(char *buf, size_t buf_size, size_t block_size) {
    if (buf_size % block_size != 0) {
        return NULL;
    }
    int padding_length = buf[buf_size - 1];
    if (padding_length <= 0 || (size_t)padding_length > block_size) {
        return NULL;
    }
    for (int i = 1; i < padding_length; ++i) {
        if (buf[buf_size - i - 1] != padding_length) {
            return NULL;
        }
    }
    return buf + buf_size - padding_length;
}

void aes_decrypt(char *encrypted, char *decrypted, size_t num_bytes, char *key) {
    AES_KEY aes_key;
    AES_set_decrypt_key((unsigned char*)key, 128, &aes_key);
    for (size_t offset = 0; offset < num_bytes; offset += 16) {
        AES_decrypt((unsigned char*)encrypted + offset, (unsigned char*)decrypted + offset, &aes_key);
    }
}

void aes_encrypt(char *decrypted, char *encrypted, size_t num_bytes, char *key) {
    AES_KEY aes_key;
    AES_set_encrypt_key((unsigned char*)key, 128, &aes_key);
    for (size_t offset = 0; offset < num_bytes; offset += 16) {
        AES_encrypt((unsigned char*)decrypted + offset, (unsigned char*)encrypted + offset, &aes_key);
    }
}

void cbc_decrypt(char *encrypted, char *decrypted, size_t num_bytes, char *iv, char *key) {
    AES_KEY aes_key;
    AES_set_decrypt_key((unsigned char*)key, 128, &aes_key);
    char prev[16];
    char buf[16];
    memcpy(prev, iv, 16);
    for (size_t offset = 0; offset < num_bytes; offset += 16) {
    	memcpy(buf, encrypted + offset, 16);
        AES_decrypt((unsigned char*)encrypted + offset, (unsigned char*)decrypted + offset, &aes_key);
    	xor(decrypted + offset, prev, 16, decrypted + offset);
        memcpy(prev, buf, 16);
    }
}

void cbc_encrypt(char *decrypted, char *encrypted, size_t num_bytes, char *iv, char *key) {
    AES_KEY aes_key;
    AES_set_encrypt_key((unsigned char*)key, 128, &aes_key);
    char buf[16];
    for (size_t offset = 0; offset < num_bytes; offset += 16) {
    	memcpy(buf, decrypted + offset, 16);
    	xor(buf, iv, 16, buf);
        AES_encrypt((unsigned char*)buf, (unsigned char*)encrypted + offset, &aes_key);
        iv = encrypted + offset;
    }
}
