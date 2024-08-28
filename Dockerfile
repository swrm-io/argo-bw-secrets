FROM golang:1.22 AS build


WORKDIR /workspace

COPY . .

RUN apt-get update && apt install musl-tools git -y
RUN go mod download

RUN C=musl-gcc CGO_ENABLED=1 go build -trimpath -ldflags '-linkmode external -extldflags "-static -Wl,-unresolved-symbols=ignore-all"' .

FROM debian:latest 
ARG CMP

RUN mkdir -p /home/argocd/cmp-server/config/plugin.yaml
COPY $CMP /home/argocd/cmp-server/config/
COPY --from=build /workspace/argo-bw-secrets /usr/local/bin/argo-bw-secrets

ENTRYPOINT ["/var/run/argocd/argocd-cmp-server"]
