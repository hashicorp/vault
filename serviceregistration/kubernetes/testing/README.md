# How to Test Manually

- `$ minikube start`
- In the Vault folder, `$ make dev XC_ARCH=amd64 XC_OS=linux XC_OSARCH=linux/amd64`
- Create a file called `vault-test.yaml` with the following contents:

```
apiVersion: v1
kind: Pod
metadata:
  name: vault
spec:
  containers:
    - name: nginx
      image: nginx
      command: [ "sh", "-c"]
      args:
      - while true; do
          echo -en '\n';
          printenv VAULT_K8S_POD_NAME VAULT_K8S_NAMESPACE;
          sleep 10;
        done;
      env:
        - name: VAULT_K8S_POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: VAULT_K8S_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
  restartPolicy: Never
```

- Create the pod: `$ kubectl apply -f vault-test.yaml`
- View the full initial state of the pod: `$ kubectl get pod vault -o=yaml > initialstate.txt`
- Drop the Vault binary into the pod: `$ kubectl cp bin/vault /vault:/`
- Drop to the shell within the pod: `$ kubectl exec -it vault -- /bin/bash`
- Install a text editor: `$ apt-get update`, `$ apt-get install nano`
- Write a test Vault config to `vault.config` like:

```
storage "inmem" {}
service_registration "kubernetes" {}
disable_mlock = true
api_addr = "http://127.0.0.1:8200"
log_level = "debug"

ui_config {
  enabled = true
}
```

- Run Vault: `$ ./vault server -config=vault.config -dev -dev-root-token-id=root`
- If 403's are received, you may need to grant RBAC, example here: https://github.com/fabric8io/fabric8/issues/6840#issuecomment-307560275
- In a separate window outside the pod, view the resulting state of the pod: `$ kubectl get pod vault -o=yaml > currentstate.txt`
- View the differences: `$ diff initialstate.txt currentstate.txt`
