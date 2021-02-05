import os

filename = "dep.txt"
with open(filename) as f:
    content = f.readlines()
    for l in content:
        name = l.split()[0]
        print(name)
        if "k8s" not in name:
            os.system("go get -u " + name)
        else:
            print("skippin: " + name)