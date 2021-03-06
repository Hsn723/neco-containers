# Rook container

# Stage1: build from source
FROM quay.io/cybozu/golang:1.13-bionic AS build

ARG ROOK_VERSION=1.2.2-master-19350d62
ARG DEP_VERSION=0.5.4
ARG ROOK_DIR=/work/go/src/github.com/rook/rook

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV GOPATH=/work/go
RUN git clone https://github.com/cybozu-go/rook.git ${ROOK_DIR}

WORKDIR ${ROOK_DIR}

ENV GO111MODULE=off
RUN git checkout v${ROOK_VERSION}
RUN make vendor
RUN make build IMAGES=""

WORKDIR ${ROOK_DIR}/images/ceph
RUN make generate-csv-ceph-templates

# Stage2: setup runtime container
FROM ceph/ceph-amd64:v14.2.6-20200115

ARG TINI_VERSION=v0.16.1
ARG ROOK_DIR=/work/go/src/github.com/rook/rook

COPY --from=build ${ROOK_DIR}/_output/bin/linux_amd64/rook \
    ${ROOK_DIR}/_output/bin/linux_amd64/rookflex \
    ${ROOK_DIR}/images/ceph/toolbox.sh \
    /usr/local/bin/

COPY --from=build ${ROOK_DIR}/cluster/examples/kubernetes/ceph/csi/template /etc/ceph-csi
COPY --from=build ${ROOK_DIR}/cluster/examples/kubernetes/ceph/monitoring /etc/ceph-monitoring
COPY --from=build ${ROOK_DIR}/cluster/olm/ceph/templates /etc/ceph-csv-templates

ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static-amd64 /tini
RUN chmod +x /tini

USER 10000:10000

ENTRYPOINT ["/tini", "--", "/usr/local/bin/rook"]
