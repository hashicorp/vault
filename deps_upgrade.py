import os

filename = "dep.txt"
with open(filename) as f:
    content = f.readlines()
    for l in content:
        name = l.split()[0]
        print(name)
        if "k8s" not in name:
            os.system("go get -u " + name)
            # google.golang.org/grpc also 
            # github.com/hashicorp/go-discover (go mod vendor CI issue)
        else:
            print("skippin: " + name)