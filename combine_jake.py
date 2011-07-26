import glob, datetime, time
import re, sys

match = re.compile(r'(\d\d\d\d)-(\d\d)-(\d\d) (\d\d):(\d\d):(\d\d) <(.*?)> (.*)')

data = []
fname = "/Users/zeebo/Downloads/#okcodev.log"

with open(fname) as hnd:
	for line in hnd:
		m = match.match(line)
		if m is not None:
			groups = list(m.groups())
			year, month, day, hours, mins, secs, name, message = map(int, groups[:6]) + groups[6:]
			dtime = datetime.datetime(year, month, day, hours, mins, secs)
			name = name.lstrip("@+ ")
			data.append(' '.join((str(int(time.mktime(dtime.timetuple()))), name, message)))

with open('log.txt', 'w') as hnd:
	hnd.write('\n'.join(data))

#(year, month, day, hour, minute, second)