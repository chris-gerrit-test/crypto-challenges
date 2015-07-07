void xor(byte* in1, byte* in2, size_t num_bytes, byte* out) {
	for (int i = 0; i < num_bytes; i++) {
		out[i] = in1[i] ^ in2[i];
	}
}
