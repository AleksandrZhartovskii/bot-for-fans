#!/bin/sh

# Path for storing logs
readonly logsDir='/opt/logs'

# Deploy bot service
docker build -t bot-for-fans/bot:latest -f bot.Dockerfile .
docker stop bot && docker rm bot
mkdir -p $logsDir/bot
docker run --name=bot --net=host -v $logsDir/bot:/app/_logs -d bot-for-fans/bot:latest

# Deploy release-tracker service
docker build -t bot-for-fans/tracker:latest -f tracker.Dockerfile .
docker stop tracker && docker rm tracker
mkdir -p $logsDir/tracker
docker run --name=tracker --net=host -v $logsDir/tracker:/app/_logs -d bot-for-fans/tracker:latest