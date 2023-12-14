git clone https://github.com/hashicorp/vault-plugin-secrets-kv.git vault-plugin-secrets-kv
#git clone https://github.com/hashicorp/vault-plugin-secrets-kv.git vault-plugin-secrets-kv-pod

mkdir docker-kv
mkdir podman-kv

cd vault-plugin-secrets-kv ; CGO_ENABLED=0 GOOS=linux go build -o ../docker-kv/kv-dock cmd/vault-plugin-secrets-kv/main.go ; cd ../docker-kv ; touch Containerfile
cd vault-plugin-secrets-kv ; CGO_ENABLED=0 GOOS=linux go build -o ../podman-kv/kv-pod cmd/vault-plugin-secrets-kv/main.go ; cd ../podman-kv ; touch Containerfile

FROM ubuntu
COPY kv /bin/kv
ENTRYPOINT [ "/bin/kv" ]

# Build a docker container
# cd vault-plugin-secrets-kv ; docker build -t hashicorp/vault-plugin-secrets-kv:dock-container . ; cd ..
#docker build -t hashicorp/vault-plugin-secrets-kv:dock-container -f Containerfile .
docker build -t hashicorp/vault-plugin-secrets-kv-dock:dock-container -f Containerfile .

--runtime=runsc

# Build a podman container
# cd vault-plugin-secrets-kv ; podman build -t hashicorp/vault-plugin-secrets-kv:pod-container . ; cd ..
podman build -t hashicorp/vault-plugin-secrets-kv-pod:pod-container -f Containerfile .

# Register the gVisor runtime
vault plugin runtime register -oci_runtime=runsc -type=container gvisor-rt

# Register the default docker runtime
vault plugin runtime register -oci_runtime=runc -type=container docker-rt

# Register the default podman runtime
vault plugin runtime register -oci_runtime=crun -type=container podman-rt

# Confirm our images are available
docker images --no-trunc hashicorp/vault-plugin-secrets-kv-dock:dock-container
podman images --no-trunc hashicorp/vault-plugin-secrets-kv-pod:pod-container 

# Save the SHA
docker images --no-trunc --format="{{ .ID }}" hashicorp/vault-plugin-secrets-kv-dock:dock-container | sort | uniq | cut -d: -f2 > dock-container.sha
podman images --no-trunc --format="{{ .ID }}" hashicorp/vault-plugin-secrets-kv-pod:pod-container | sort | uniq | cut -d: -f2 > pod-container.sha

# Register the docker version of the plugin
cd ~
vault plugin register -runtime=docker-rt -sha256="$(cat docker-kv/dock-container.sha)" -oci_image=hashicorp/vault-plugin-secrets-kv-dock secret scc-dock
vault secrets enable scc-dock

vault plugin register -sha256="$(cat docker-kv/dock-container.sha)" -oci_image=hashicorp/vault-plugin-secrets-kv-dock secret scc-gvis
vault secrets enable scc-gvis

vault plugin register -runtime=gvisor-rt -sha256="$(cat docker-kv/dock-container.sha)" -oci_image=hashicorp/vault-plugin-secrets-kv-dock secret scc-gvis2
vault secrets enable scc-gvis2


vault write scc-dock/testing subject=containers
vault read scc-dock/testing

vault plugin register -runtime=podman-rt -sha256="$(cat podman-kv/pod-container.sha)" -oci_image=hashicorp/vault-plugin-secrets-kv-pod secret scc-pod

vault plugin register -sha256="$(cat podman-kv/pod-container.sha)" -oci_image=hashicorp/vault-plugin-secrets-kv-pod secret scc-pod2
vault secrets enable scc-pod2

vault write scc-pod/testing subject=containers
vault read scc-pod/testing
