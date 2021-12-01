FROM therecipe/qt:linux_static_base as fcitx

# RUN apt-get update -y

# ENV QT_CI_PACKAGES="qt.qt5.5130.gcc_64,qt.qt5.5130.qtwebglplugin,qt.qt5.5130.qtwebengine"


# RUN echo yoo $QT_CI_LOGIN
# COPY ./extract-qt-installer /root
# COPY ./install-qt /root
# RUN chmod +x /root/extract-qt-installer && chmod +x /root/install-qt
# RUN cd /root && ./install-qt 5.13.0
# #RUN rm -rf /qt/Docs /qt/Tools /qt/Examples /qt/5.12.0/gcc_64/doc/

ENV QT_DIR=/opt/Qt
ENV QT_VERSION=5.13.0
ENV QT_API=5.13.0
ENV PATH="/usr/local/go/bin:${PATH}"
ENV QT_DOCKER true
ENV QT_STATIC true
RUN apt-get update -y
RUN apt-get -qq update && apt-get --no-install-recommends -qq -y install libfontconfig1-dev libfreetype6-dev libx11-dev libxext-dev libxfixes-dev libxi-dev libxrender-dev libxcb1-dev libx11-xcb-dev libxcb-glx0-dev
RUN apt-get install wget bison build-essential gperf flex ruby python libasound2-dev libbz2-dev libcap-dev \
     libcups2-dev libdrm-dev libegl1-mesa-dev libgcrypt11-dev libnss3-dev libpci-dev libpulse-dev libudev-dev \
     libxtst-dev gyp ninja-build libfreetype6-dev libfontconfig-dev libevent-dev -y
RUN apt-get -y install libglu1-mesa-dev  libglib2.0-dev curl git
RUN GO=go1.17.3.linux-amd64.tar.gz; curl -sL https://dl.google.com/go/$GO |  tar -xz -C /usr/local

ENV GO111MODULE=off
RUN export GO111MODULE=off; go get -v github.com/therecipe/qt/cmd/...
RUN echo nocache=3
COPY src /root/go/src/github.com/threefoldtech/yggdrasil-desktop-client 
RUN cd /root/go/src/github.com/threefoldtech/yggdrasil-desktop-client && go get || true
ARG QT_CI_LOGIN
ARG QT_CI_PASSWORD
ENV QT_CI_LOGIN=$QT_CI_LOGIN
ENV QT_CI_PASSWORD=$QT_CI_PASSWORD

RUN $(go env GOPATH)/bin/qtdeploy test desktop /root/go/src/github.com/threefoldtech/yggdrasil-desktop-client
#RUN  go get -v github.com/therecipe/qt/internal/examples/widgets/systray
#RUN $(go env GOPATH)/bin/qtdeploy test desktop github.com/therecipe/examples/widgets/systray/systray.go
