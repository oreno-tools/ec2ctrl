#!/usr/bin/env python

import re
import sys

def tasks():
    for line in sys.stdin:
        match = re.match(r'^([a-z\-]*):(?:.*)#\s*(.*)$', line)
        if match:
            print match.group(1).ljust(12) + '....... ' +match.group(2)
