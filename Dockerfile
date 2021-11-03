FROM ubuntu:latest

RUN apt update && \
    apt install -y --no-install-recommends nasm gcc gcc-multilib gdb sasm
RUN apt clean

CMD ["/usr/bin/sasm"]
