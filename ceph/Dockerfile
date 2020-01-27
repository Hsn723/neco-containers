# Ceph container image

# Stage1: build from source
FROM quay.io/cybozu/ubuntu:18.04 AS build

ARG CEPH_VERSION=14.2.6

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN apt-get update && apt-get install --no-install-recommends -y git
RUN git clone https://github.com/ceph/ceph.git /work/ceph

WORKDIR /work/ceph
RUN git checkout v${CEPH_VERSION} \
    && git submodule update --init --recursive
RUN ./install-deps.sh \
    && ./do_cmake.sh -DCMAKE_BUILD_TYPE=RelWithDebInfo

WORKDIR /work/ceph/build
RUN make -j20

# Stage2: setup runtime container
FROM quay.io/cybozu/ubuntu:18.04

COPY --from=build /work/ceph/build /ceph
WORKDIR /ceph
RUN make install

WORKDIR /
# install dependencies

USER 10000:10000
