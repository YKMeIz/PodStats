FROM registry.access.redhat.com/ubi8/ubi-minimal

RUN microdnf --setopt=install_weak_deps=0 --setopt=tsflags=nodocs install -y device-mapper-libs \
            && microdnf clean all

WORKDIR /
COPY ./podstats .
COPY ./dist ./dist

EXPOSE 9023/tcp


ENTRYPOINT ["/podstats"]
