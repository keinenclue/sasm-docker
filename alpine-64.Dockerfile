# alpine-64 image

FROM alpine as base

RUN apk update && \
    apk add nasm gcc gdb fontconfig musl-dev libc-dev &&\
    rm /var/cache/apk/*

FROM base as build

RUN apk add build-base qt5-qtbase-dev unzip curl \
    msttcorefonts-installer && \
    update-ms-fonts && \
    fc-cache -f && \
    curl -L -o sasm.zip https://codeload.github.com/schreiberx/SASM/zip/refs/heads/master && \
    unzip sasm.zip -d /home && \
    cd /home/SASM-master && qmake-qt5 && make && make install


# Source: https://gist.github.com/bcardiff/85ae47e66ff0df35a78697508fcb49af#gistcomment-2078660 
RUN ldd /usr/bin/sasm | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname /home/deps%); cp % /home/deps%;'

RUN ldd /usr/lib/qt5/plugins/platforms/libqxcb.so | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname /home/deps%); cp % /home/deps%;'


RUN mkdir -p /home/deps/usr/share/fonts/truetype/msttcorefonts
RUN cp /usr/share/fonts/truetype/msttcorefonts/Courier* /home/deps/usr/share/fonts/truetype/msttcorefonts/
RUN cp /usr/share/fonts/truetype/msttcorefonts/Arial* /home/deps/usr/share/fonts/truetype/msttcorefonts/

FROM base as runtime    

# Copy dependencies
COPY --from=build /home/deps /
COPY --from=build /usr/lib/qt5/plugins/platforms/libqxcb.so /usr/lib/qt5/plugins/platforms/libqxcb.so

# Copy sasm
COPY --from=build /usr/bin/sasm /usr/bin/sasm
COPY --from=build /usr/share/sasm /usr/share/sasm

# Install fonts
RUN fc-cache -f -v

CMD ["sasm"]
