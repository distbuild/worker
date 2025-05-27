FROM craftslab/ubuntu:22.04-aosp

USER aosp
WORKDIR /home/aosp
RUN mkdir -p .distbuild/bin && \
    mkdir -p .distbuild/workspace

USER aosp
WORKDIR /home/aosp
RUN git clone https://github.com/craftslab/gitclone.git --depth=1 && \
    pushd gitclone && \
    make build && \
    cp bin/clone /home/aosp/.distbuild/bin/ && \
    cp config.yml /home/aosp/.distbuild/bin/clone-config.yml && \
    popd && \
    rm -rf gitclone

USER aosp
WORKDIR /home/aosp
RUN mkdir worker
COPY . worker/
RUN pushd worker && \
    make build && \
    cp bin/worker /home/aosp/.distbuild/bin/ && \
    popd && \
    rm -rf worker

ENV PATH=/home/aosp/.distbuild/bin:$PATH
CMD ["worker", "--listen-address=:39090", "--workspace-path=/home/aosp/.distbuild/workspace"]
