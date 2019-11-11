// https://github.com/lufia/iostat/blob/9f7362b77ad333b26c01c99de52a11bdb650ded2/iostat_darwin.c
#include <stdint.h>
#include <CoreFoundation/CoreFoundation.h>
#include "disk_darwin.h"

#define IOKIT	1	/* to get io_name_t in device_types.h */

#include <IOKit/IOKitLib.h>
#include <IOKit/storage/IOBlockStorageDriver.h>
#include <IOKit/storage/IOMedia.h>
#include <IOKit/IOBSD.h>

#include <mach/mach_host.h>

static int getdrivestat(io_registry_entry_t d, DriveStats *stat);
static int fillstat(io_registry_entry_t d, DriveStats *stat);

int
readdrivestat(DriveStats a[], int n)
{
	mach_port_t port;
	CFMutableDictionaryRef match;
	io_iterator_t drives;
	io_registry_entry_t d;
	kern_return_t status;
	int na, rv;

	IOMasterPort(bootstrap_port, &port);
	match = IOServiceMatching("IOMedia");
	CFDictionaryAddValue(match, CFSTR(kIOMediaWholeKey), kCFBooleanTrue);
	status = IOServiceGetMatchingServices(port, match, &drives);
	if(status != KERN_SUCCESS)
		return -1;

	na = 0;
	while(na < n && (d=IOIteratorNext(drives)) > 0){
		rv = getdrivestat(d, &a[na]);
		if(rv < 0)
			return -1;
		if(rv > 0)
			na++;
		IOObjectRelease(d);
	}
	IOObjectRelease(drives);
	return na;
}

static int
getdrivestat(io_registry_entry_t d, DriveStats *stat)
{
	io_registry_entry_t parent;
	kern_return_t status;
	CFDictionaryRef props;
	CFStringRef name;
	CFNumberRef num;
	int rv;

	memset(stat, 0, sizeof *stat);
	status = IORegistryEntryGetParentEntry(d, kIOServicePlane, &parent);
	if(status != KERN_SUCCESS)
		return -1;
	if(!IOObjectConformsTo(parent, "IOBlockStorageDriver")){
		IOObjectRelease(parent);
		return 0;
	}

	status = IORegistryEntryCreateCFProperties(d, (CFMutableDictionaryRef *)&props, kCFAllocatorDefault, kNilOptions);
	if(status != KERN_SUCCESS){
		IOObjectRelease(parent);
		return -1;
	}
	name = (CFStringRef)CFDictionaryGetValue(props, CFSTR(kIOBSDNameKey));
	CFStringGetCString(name, stat->name, NAMELEN, CFStringGetSystemEncoding());
	num = (CFNumberRef)CFDictionaryGetValue(props, CFSTR(kIOMediaSizeKey));
	CFNumberGetValue(num, kCFNumberSInt64Type, &stat->size);
	num = (CFNumberRef)CFDictionaryGetValue(props, CFSTR(kIOMediaPreferredBlockSizeKey));
	CFNumberGetValue(num, kCFNumberSInt64Type, &stat->blocksize);
	CFRelease(props);

	rv = fillstat(parent, stat);
	IOObjectRelease(parent);
	if(rv < 0)
		return -1;
	return 1;
}

static struct {
	char *key;
	size_t off;
} statstab[] = {
	{kIOBlockStorageDriverStatisticsBytesReadKey, offsetof(DriveStats, read)},
	{kIOBlockStorageDriverStatisticsBytesWrittenKey, offsetof(DriveStats, written)},
	{kIOBlockStorageDriverStatisticsReadsKey, offsetof(DriveStats, nread)},
	{kIOBlockStorageDriverStatisticsWritesKey, offsetof(DriveStats, nwrite)},
	{kIOBlockStorageDriverStatisticsTotalReadTimeKey, offsetof(DriveStats, readtime)},
	{kIOBlockStorageDriverStatisticsTotalWriteTimeKey, offsetof(DriveStats, writetime)},
	{kIOBlockStorageDriverStatisticsLatentReadTimeKey, offsetof(DriveStats, readlat)},
	{kIOBlockStorageDriverStatisticsLatentWriteTimeKey, offsetof(DriveStats, writelat)},
};

static int
fillstat(io_registry_entry_t d, DriveStats *stat)
{
	CFDictionaryRef props, v;
	CFNumberRef num;
	kern_return_t status;
	typeof(statstab[0]) *bp, *ep;

	status = IORegistryEntryCreateCFProperties(d, (CFMutableDictionaryRef *)&props, kCFAllocatorDefault, kNilOptions);
	if(status != KERN_SUCCESS)
		return -1;
	v = (CFDictionaryRef)CFDictionaryGetValue(props, CFSTR(kIOBlockStorageDriverStatisticsKey));
	if(v == NULL){
		CFRelease(props);
		return -1;
	}

	ep = &statstab[sizeof(statstab)/sizeof(statstab[0])];
	for(bp = &statstab[0]; bp < ep; bp++){
		CFStringRef s;

		s = CFStringCreateWithCString(kCFAllocatorDefault, bp->key, CFStringGetSystemEncoding());
		num = (CFNumberRef)CFDictionaryGetValue(v, s);
		if(num)
			CFNumberGetValue(num, kCFNumberSInt64Type, ((char*)stat)+bp->off);
		CFRelease(s);
	}

	CFRelease(props);
	return 0;
}
