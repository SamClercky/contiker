FROM contiker/contiki-ng

RUN apt update && DEBIAN_FRONTEND=noninteractive apt install -y \
  dos2unix \
  wireshark \
  firefox \
  curl \
  mosquitto \
  mosquitto-clients \
  && rm -rf /var/lib/apt/lists/*

# Very bad way of installing node 22
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ENV BASH_ENV=~/.bashrc
ENV VOLTA_HOME="/root/.volta"
ENV PATH="$VOLTA_HOME/bin:$PATH"
RUN curl https://get.volta.sh | bash
RUN volta install node@22.5.1

# Install the coap tools
RUN npm install coap-cli -g && \
  ln -s /usr/bin/nodejs /usr/bin/node

COPY oobe.sh /etc/
COPY ./wsl-distribution.conf /etc/
COPY ./wsl.conf /etc/
