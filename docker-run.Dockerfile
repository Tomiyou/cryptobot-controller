FROM tomiyou/crypto-arbitrage:latest
ARG configPath

WORKDIR /home/appuser/

# Copy all runtime configs
COPY ./keys keys/
COPY $configPath config.yaml

CMD ["config.yaml"]