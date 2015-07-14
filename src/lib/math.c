void xor(char* in1, char *in2, size_t num_bytes, char *out) {
	for (size_t i = 0; i < num_bytes; i++) {
		out[i] = in1[i] ^ in2[i];
	}
}

void repeated_xor(char *in, size_t num_bytes, char *key, size_t key_len, char *out) {
	for (size_t i = 0; i < num_bytes; i++) {
		out[i] = in[i] ^ key[i % key_len];
	}
}