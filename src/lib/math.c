void xor(byte* in1, byte *in2, size_t num_bytes, byte *out) {
	for (int i = 0; i < num_bytes; i++) {
		out[i] = in1[i] ^ in2[i];
	}
}

void repeated_xor(byte *in, size_t num_bytes, byte *key, size_t key_len, byte *out) {
	for (int i = 0; i < num_bytes; i++) {
		out[i] = in[i] ^ key[i % key_len];
	}
}