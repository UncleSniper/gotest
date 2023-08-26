#!/usr/bin/python3

print('// CAUTION: This file is generated. Do not edit it!')
print()
print('package gotest')

for args in range(11):
	for rets in range(5):
		atpl = alst = rtpl = rlst = ''
		for aidx in range(args):
			if len(atpl) > 0:
				atpl += ', '
			if len(alst) > 0:
				alst += ', '
			aname = 'Argument' + str(aidx) + 'T'
			atpl += aname + ' any'
			alst += aname
		for ridx in range(rets):
			if len(rtpl) > 0:
				rtpl += ', '
			if len(rlst) > 0:
				rlst += ', '
			rname = 'Return' + str(ridx) + 'T'
			rtpl += rname + ' any'
			rlst += rname
		alltpl = atpl
		if len(atpl) > 0 and len(rtpl) > 0:
			alltpl += ', '
		alltpl += rtpl
		if len(alltpl) > 0:
			alltpl = '[' + alltpl + ']'
		if len(rlst) == 0:
			allrlst = ''
		else:
			allrlst = ' (' + rlst + ')'
		print()
		print('func IsFuncNil' + str(args) + 'to' + str(rets) + alltpl + '(f func(' + alst + ')' + allrlst + ') bool {')
		print('\treturn f == nil')
		print('}')
