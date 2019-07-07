FROM tomiyou/crypto-arbitrage:latest
ARG configPath

ARG containerName
ENV CONTAINER_NAME=$containerName

WORKDIR /home/appuser/

# Copy all runtime configs
COPY ./keys keys/
COPY $configPath config.yaml