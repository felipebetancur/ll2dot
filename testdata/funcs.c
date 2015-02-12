int bar(int x) {
	while (x < 1000) {
		x *= 2;
	}
	return x;
}

int main(int argc, char **argv) {
	int x;
	if (x < 3) {
		return bar(x);
	}
	return bar(x*2);
}
