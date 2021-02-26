import os
import sys

filename = sys.argv[1]
with open(filename) as f:
    content = f.readlines()
    for l in content:
        name = l.split()[0]
        print(name)
        os.system("go get " + name + "@latest")
        # if "k8s" not in name:
        #     os.system("go get " + name + "@latest")
        #     # google.golang.org/grpc also 
        #     # github.com/hashicorp/go-discover (go mod vendor CI issue)
        # else:
        #     print("skippin: " + name)