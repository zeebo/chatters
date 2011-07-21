import glob, datetime, time
import re, sys

match = re.compile(r'\[(\d\d):(\d\d):(\d\d)\] (\w+?): (.*)')

data = []

files = glob.glob('logs/*.txt')
for fname in files:
	with open(fname) as hnd:
		year, month, day = map(int, fname[5:-4].split('-'))
		for line in hnd:
			m = match.match(line)
			if m is not None:
				groups = list(m.groups())
				hours, mins, secs, name, message = map(int, groups[:3]) + groups[3:]
				dtime = datetime.datetime(year, month, day, hours, mins, secs)
				data.append(' '.join((str(int(time.mktime(dtime.timetuple()))), name, message)))

with open('log.txt', 'w') as hnd:
	hnd.write('\n'.join(data))

#(year, month, day, hour, minute, second)