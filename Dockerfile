FROM i386/alpine as base

# We don't need other stuff
# RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf && \
#     echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf
    # RUN echo "CONFIG += static \n QTPLUGIN += qsqloci qgif" >> /home/SASM-master/SASM.pro
RUN apk update && \
    apk add nasm gcc gdb fontconfig && \
    rm /var/cache/apk/*

FROM base as build

RUN apk update && \
    apk add nasm gcc gdb fontconfig g++ make qt5-qtbase-dev unzip curl musl-dev \
    msttcorefonts-installer && \
    echo "#### Setup fonts ####" && \
    update-ms-fonts && \
    fc-cache -f && \
    echo "#### Install sasm ####" && \
    curl -L -o sasm.zip https://github.com/Dman95/SASM/archive/refs/heads/master.zip && \
    unzip sasm.zip -d /home && rm sasm.zip && \
    cd /home/SASM-master && qmake-qt5 && make && make install && \
    mv /home/SASM-master/sasm /home/SASM-master/Linux/sasm && \
    rm /var/cache/apk/*
# RUN mv /home/SASM-master/Linux /home
# Clean stuff up

CMD ["/usr/bin/sasm"]

# FROM base as runtime    

# # Copy fonts
# COPY --from=build /usr/share/fonts/truetype/msttcorefonts/Courier* /usr/share/fonts/truetype/msttcorefonts/
# COPY --from=build /usr/share/fonts/truetype/msttcorefonts/cour* /usr/share/fonts/truetype/msttcorefonts/
# COPY --from=build /usr/share/fonts/truetype/msttcorefonts/Arial* /usr/share/fonts/truetype/msttcorefonts/
# COPY --from=build /usr/share/fonts/truetype/msttcorefonts/arial* /usr/share/fonts/truetype/msttcorefonts/
# RUN fc-cache -f -v

# # Copy sasm
# COPY --from=build /home/SASM-master/Linux /home/sasm

# # Copy dependencies
# # Find them out by running ldd /home/sasm/sasm > and ldd /usr/lib/libQt5XcbQpa.so.5
# # ldd /home/sasm/sasm > dep
# # ldd /usr/lib/libQt5XcbQpa.so.5 >> dep
# # Then remove duplicates
# # awk '!seen[$0]++' dep
# # 
# # Now use regex to remove the beginning and end
# # Using multiline selections 
# # add COPY --from=build in the front and duplicate the paths

# COPY --from=build /usr/lib/qt5 /usr/lib/qt5
# COPY --from=build /usr/lib/libQt5XcbQpa.so.5 /usr/lib/libQt5XcbQpa.so.5
# COPY --from=build /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1 
# COPY --from=build /usr/lib/libfontconfig.so.1 /usr/lib/libfontconfig.so.1 
# COPY --from=build /usr/lib/libfreetype.so.6 /usr/lib/libfreetype.so.6 
# COPY --from=build /usr/lib/libQt5Gui.so.5 /usr/lib/libQt5Gui.so.5 
# COPY --from=build /usr/lib/libQt5DBus.so.5 /usr/lib/libQt5DBus.so.5 
# COPY --from=build /usr/lib/libQt5Core.so.5 /usr/lib/libQt5Core.so.5 
# COPY --from=build /usr/lib/libX11-xcb.so.1 /usr/lib/libX11-xcb.so.1 
# COPY --from=build /usr/lib/libxcb-icccm.so.4 /usr/lib/libxcb-icccm.so.4 
# COPY --from=build /usr/lib/libxcb-image.so.0 /usr/lib/libxcb-image.so.0 
# COPY --from=build /usr/lib/libxcb-shm.so.0 /usr/lib/libxcb-shm.so.0 
# COPY --from=build /usr/lib/libxcb-keysyms.so.1 /usr/lib/libxcb-keysyms.so.1 
# COPY --from=build /usr/lib/libxcb-randr.so.0 /usr/lib/libxcb-randr.so.0 
# COPY --from=build /usr/lib/libxcb-render-util.so.0 /usr/lib/libxcb-render-util.so.0 
# COPY --from=build /usr/lib/libxcb-render.so.0 /usr/lib/libxcb-render.so.0 
# COPY --from=build /usr/lib/libxcb-shape.so.0 /usr/lib/libxcb-shape.so.0 
# COPY --from=build /usr/lib/libxcb-sync.so.1 /usr/lib/libxcb-sync.so.1 
# COPY --from=build /usr/lib/libxcb-xfixes.so.0 /usr/lib/libxcb-xfixes.so.0 
# COPY --from=build /usr/lib/libxcb-xinerama.so.0 /usr/lib/libxcb-xinerama.so.0 
# COPY --from=build /usr/lib/libxcb-xkb.so.1 /usr/lib/libxcb-xkb.so.1 
# COPY --from=build /usr/lib/libxcb-xinput.so.0 /usr/lib/libxcb-xinput.so.0 
# COPY --from=build /usr/lib/libxcb.so.1 /usr/lib/libxcb.so.1 
# COPY --from=build /usr/lib/libX11.so.6 /usr/lib/libX11.so.6 
# COPY --from=build /usr/lib/libSM.so.6 /usr/lib/libSM.so.6 
# COPY --from=build /usr/lib/libICE.so.6 /usr/lib/libICE.so.6 
# COPY --from=build /usr/lib/libxkbcommon-x11.so.0 /usr/lib/libxkbcommon-x11.so.0 
# COPY --from=build /usr/lib/libxkbcommon.so.0 /usr/lib/libxkbcommon.so.0 
# COPY --from=build /usr/lib/libglib-2.0.so.0 /usr/lib/libglib-2.0.so.0 
# COPY --from=build /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6 
# COPY --from=build /usr/lib/libgcc_s.so.1 /usr/lib/libgcc_s.so.1 
# COPY --from=build /usr/lib/libexpat.so.1 /usr/lib/libexpat.so.1 
# COPY --from=build /lib/libuuid.so.1 /lib/libuuid.so.1 
# COPY --from=build /usr/lib/libbz2.so.1 /usr/lib/libbz2.so.1 
# COPY --from=build /usr/lib/libpng16.so.16 /usr/lib/libpng16.so.16 
# COPY --from=build /lib/libz.so.1 /lib/libz.so.1 
# COPY --from=build /usr/lib/libbrotlidec.so.1 /usr/lib/libbrotlidec.so.1 
# COPY --from=build /usr/lib/libGL.so.1 /usr/lib/libGL.so.1 
# COPY --from=build /usr/lib/libharfbuzz.so.0 /usr/lib/libharfbuzz.so.0 
# COPY --from=build /usr/lib/libdbus-1.so.3 /usr/lib/libdbus-1.so.3 
# COPY --from=build /usr/lib/libicui18n.so.67 /usr/lib/libicui18n.so.67 
# COPY --from=build /usr/lib/libicuuc.so.67 /usr/lib/libicuuc.so.67 
# COPY --from=build /usr/lib/libpcre2-16.so.0 /usr/lib/libpcre2-16.so.0 
# COPY --from=build /usr/lib/libzstd.so.1 /usr/lib/libzstd.so.1 
# COPY --from=build /usr/lib/libxcb-util.so.1 /usr/lib/libxcb-util.so.1 
# COPY --from=build /usr/lib/libXau.so.6 /usr/lib/libXau.so.6 
# COPY --from=build /usr/lib/libXdmcp.so.6 /usr/lib/libXdmcp.so.6 
# COPY --from=build /usr/lib/libpcre.so.1 /usr/lib/libpcre.so.1 
# COPY --from=build /usr/lib/libintl.so.8 /usr/lib/libintl.so.8 
# COPY --from=build /usr/lib/libbrotlicommon.so.1 /usr/lib/libbrotlicommon.so.1 
# COPY --from=build /usr/lib/libglapi.so.0 /usr/lib/libglapi.so.0 
# COPY --from=build /usr/lib/libdrm.so.2 /usr/lib/libdrm.so.2 
# COPY --from=build /usr/lib/libxcb-glx.so.0 /usr/lib/libxcb-glx.so.0 
# COPY --from=build /usr/lib/libxcb-dri2.so.0 /usr/lib/libxcb-dri2.so.0 
# COPY --from=build /usr/lib/libXext.so.6 /usr/lib/libXext.so.6 
# COPY --from=build /usr/lib/libXfixes.so.3 /usr/lib/libXfixes.so.3 
# COPY --from=build /usr/lib/libXxf86vm.so.1 /usr/lib/libXxf86vm.so.1 
# COPY --from=build /usr/lib/libxcb-dri3.so.0 /usr/lib/libxcb-dri3.so.0 
# COPY --from=build /usr/lib/libxcb-present.so.0 /usr/lib/libxcb-present.so.0 
# COPY --from=build /usr/lib/libxshmfence.so.1 /usr/lib/libxshmfence.so.1 
# COPY --from=build /usr/lib/libgraphite2.so.3 /usr/lib/libgraphite2.so.3 
# COPY --from=build /usr/lib/libicudata.so.67 /usr/lib/libicudata.so.67 
# COPY --from=build /usr/lib/libbsd.so.0 /usr/lib/libbsd.so.0 
# COPY --from=build /usr/lib/libmd.so.0 /usr/lib/libmd.so.0 
# COPY --from=build /usr/lib/libQt5Widgets.so.5 /usr/lib/libQt5Widgets.so.5 
# COPY --from=build /usr/lib/libQt5Network.so.5 /usr/lib/libQt5Network.so.5 
# COPY --from=build /lib/libssl.so.1.1 /lib/libssl.so.1.1 
# COPY --from=build /lib/libcrypto.so.1.1 /lib/libcrypto.so.1.1 

# CMD ["/home/sasm/sasm"]

# #/usr/lib/qt5/plugins/platforms/libqxcb.so