// https://github.com/lufia/iostat/blob/9f7362b77ad333b26c01c99de52a11bdb650ded2/iostat_darwin.h
typedef struct DriveStats DriveStats;
typedef struct CPUStats CPUStats;

enum {
	NDRIVE = 16,
	NAMELEN = 31
};

struct DriveStats {
	char name[NAMELEN+1];
	int64_t size;
	int64_t blocksize;

	int64_t read;
	int64_t written;
	int64_t nread;
	int64_t nwrite;
	int64_t readtime;
	int64_t writetime;
	int64_t readlat;
	int64_t writelat;
};

struct CPUStats {
	natural_t user;
	natural_t nice;
	natural_t sys;
	natural_t idle;
};

extern int readdrivestat(DriveStats a[], int n);
extern int readcpustat(CPUStats *cpu);
