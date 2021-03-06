[![Docker Repository on Quay](https://quay.io/repository/cybozu/hyperkube/status "Docker Repository on Quay")](https://quay.io/repository/cybozu/hyperkube)

hyperkube container
===================

[hyperkube](https://github.com/kubernetes/kubernetes/tree/master/cluster/images/hyperkube) an all-in-one binary for the Kubernetes components.

Usage
-----

To launch `apiserver` by `docker run`:

    $ docker run --net=host --name apiserver -d \
        quay.io/cybozu/hyperkube:1.16 kube-apiserver \
        --advertise-address=192.168.1.101 \
        --insecure-bind-address=0.0.0.0 \
        --insecure-port=8080 \
        --enable-bootstrap-token-auth=true \
        --etcd-servers=http://192.168.1.101:2379,http://192.168.1.102:2379,http://192.168.1.103:2379 \
        --storage-backend=etcd3

To use `kubectl` by `docker run`:

    $ docker run --rm -it quay.io/cybozu/hyperkube:1.16 kubectl cluster-info
