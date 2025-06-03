FROM craftslab/ubuntu:22.04-aosp

USER aosp
WORKDIR /home/aosp
RUN mkdir -p .distbuild/boong/bin && \
    mkdir -p .distbuild/prebuilts && \
    mkdir -p .distbuild/workspace && \
    ln -s /home/aosp/.distbuild/prebuilts /home/aosp/.distbuild/workspace/prebuilts

USER aosp
WORKDIR /home/aosp
RUN git clone https://github.com/craftslab/gitclone.git --depth=1 && \
    pushd gitclone && \
    make build && \
    cp bin/clone /home/aosp/.distbuild/boong/bin/ && \
    cp config.yml /home/aosp/.distbuild/boong/bin/clone-config.yml && \
    popd && \
    sudo rm -rf gitclone

USER aosp
WORKDIR /home/aosp
RUN mkdir build
COPY . build/
RUN pushd build && \
    sudo chown -R aosp:aosp * && \
    make build && \
    cp bin/worker /home/aosp/.distbuild/boong/bin/ && \
    popd && \
    sudo rm -rf build

ENV PATH=/home/aosp/.distbuild/boong/bin:$PATH
ENV PATH=/home/aosp/.distbuild/prebuilts/clang/host/linux-x86/clang-r547379/bin:$PATH
CMD ["worker", "--listen-address=:39090", "--workspace-path=/home/aosp/.distbuild/workspace"]
