FROM contiker/contiki-ng

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y \
  dos2unix \
  wireshark \
  firefox \
  && rm -rf /var/lib/apt/lists/*

COPY oobe.sh /etc/
COPY ./wsl-distribution.conf /etc/
COPY ./wsl.conf /etc/
