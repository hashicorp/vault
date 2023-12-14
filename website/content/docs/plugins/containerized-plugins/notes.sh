vault plugin register -sha256="${SHA256}" -oci_image=hashicorp/vault-plugin-secrets-kv -version="v1.0+container" secret my-kv-plugin2


vault plugin runtime deregister -type=container runsc

vault plugin runtime register -oci_runtime=runsc -type=container runsc
vault plugin runtime register -oci_runtime=runc -type=container runc


vault plugin register -runtime=runc -sha256="${SHA256}" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-runc2

vault plugin register -runtime=runsc -sha256="${SHA256}" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-runsc2

vault plugin register -sha256="${SHA256}" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-runc3


vault secrets enable scc-runc3
vault secrets enable scc-runc2
vault secrets enable scc-runsc2



vault plugin register -runtime=runc -sha256="${SHA256}" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-runc4
vault secrets enable scc-runc4



vault plugin register -runtime=runsc -sha256="${SHA256}" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-runsc5
vault secrets enable scc-runsc5



podman build -t secrets-kv:kv-pod .

vault plugin register -runtime=podrt -sha256="$(podman images --no-trunc --format="{{ .ID }}" localhost/secrets-kv:kv-pod | cut -d: -f2)" -oci_image=localhost/secrets-kv secret scc-podrt3
vault secrets enable scc-podrt3



export SHA256=$( podman images --no-trunc --format="{{ .ID }}" secrets-kv:kv-pod | cut -d: -f2)


vault plugin runtime register -oci_runtime=runc -type=container allyourcons
vault plugin runtime register -oci_runtime=crun -type=container allyourpods


cd vault-plugin-secrets-kv ; docker build -t hashicorp/vault-plugin-secrets-kv:dock-container . ; cd ..
cd vault-plugin-secrets-kv ; podman build -t hashicorp/vault-plugin-secrets-kv:pod-container . ; cd ..


docker images --no-trunc hashicorp/vault-plugin-secrets-kv:dock-container
podman images --no-trunc hashicorp/vault-plugin-secrets-kv:pod-container 

docker images --no-trunc --format="{{ .ID }}" hashicorp/vault-plugin-secrets-kv:dock-container | sort | uniq | cut -d: -f2
podman images --no-trunc --format="{{ .ID }}" hashicorp/vault-plugin-secrets-kv:pod-container | sort | uniq | cut -d: -f2


vault plugin register -runtime=allyourcons -sha256="$(docker images --no-trunc --format="{{ .ID }}" hashicorp/vault-plugin-secrets-kv:dock-container | cut -d: -f2)" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-dock
vault secrets enable scc-dock


vault plugin register -runtime=allyourpods -sha256="$(podman images --no-trunc --format="{{ .ID }}" hashicorp/vault-plugin-secrets-kv:pod-container | cut -d: -f2)" -oci_image=hashicorp/vault-plugin-secrets-kv secret scc-pod
vault secrets enable scc-pod



IGNOR GVISOR, JUST USE THE DEFAULT RUNTIMES >_<

- Need to test headless
- Need to add explicit runtime registration
- Need to add a tab group




