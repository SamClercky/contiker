FROM contiker/contiki-ng

RUN DEBIAN_FRONTEND=noninteractive apt install -y \
    dos2unix \
    && rm -rf /var/lib/apt/lists/*

COPY oobe.sh /etc/
COPY ./wsl-distribution.conf /etc/
COPY ./wsl.conf /etc/
