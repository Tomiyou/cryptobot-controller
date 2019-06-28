FROM tomiyou/crypto-arbitrage:latest
ARG configPath

ARG orgConfigName
ENV CONFIG_NAME=$orgConfigName

WORKDIR /home/appuser/

# Copy all runtime configs
COPY ./keys keys/
COPY $configPath config.yaml