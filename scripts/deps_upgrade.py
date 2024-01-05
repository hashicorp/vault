# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

import os
import sys

filename = sys.argv[1]
with open(filename) as f:
    content = f.readlines()
    for l in content:
        name = l.split()[0]
        print(name)
        os.system("go get " + name + "@latest")