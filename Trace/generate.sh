#! /bin/bash
# 删除
rm channel-artifacts crypto-config -rf
# 生成证书
cryptogen generate --config crypto-config.yaml
# 创建目录
mkdir channel-artifacts
# 生成创始区块文件
configtxgen -profile Genesis -outputBlock ./channel-artifacts/genesis.block
# 生成通道文件
configtxgen -profile Channel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID tracechannel
