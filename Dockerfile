FROM i386/alpine as base

RUN apk update && \
    apk add nasm gcc gdb fontconfig musl-dev

RUN rm /var/cache/apk/*

FROM base as build

RUN apk add build-base qt5-qtbase-dev unzip curl \
    msttcorefonts-installer && \
    update-ms-fonts && \
    fc-cache -f && \
    curl -L -o sasm.zip https://github.com/Dman95/SASM/archive/refs/heads/master.zip && \
    unzip sasm.zip -d /home && \
    cd /home/SASM-master && qmake-qt5 && make && make install


# Source: https://gist.github.com/bcardiff/85ae47e66ff0df35a78697508fcb49af#gistcomment-2078660 
RUN ldd /usr/bin/sasm | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname /home/deps%); cp % /home/deps%;'

RUN ldd /usr/lib/qt5/plugins/platforms/libqxcb.so | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname /home/deps%); cp % /home/deps%;'

FROM base as runtime    

# Copy fonts
COPY --from=build /usr/share/fonts/truetype/msttcorefonts/Courier* /usr/share/fonts/truetype/msttcorefonts/
COPY --from=build /usr/share/fonts/truetype/msttcorefonts/cour* /usr/share/fonts/truetype/msttcorefonts/
COPY --from=build /usr/share/fonts/truetype/msttcorefonts/Arial* /usr/share/fonts/truetype/msttcorefonts/
COPY --from=build /usr/share/fonts/truetype/msttcorefonts/arial* /usr/share/fonts/truetype/msttcorefonts/
RUN fc-cache -f -v

# Copy dependencies
COPY --from=build /home/deps /
COPY --from=build /usr/lib/qt5/plugins/platforms/libqxcb.so /usr/lib/qt5/plugins/platforms/libqxcb.so

# Copy sasm
COPY --from=build /usr/bin/sasm /usr/bin/sasm
COPY --from=build /usr/share/sasm /usr/share/sasm

CMD ["sasm"]