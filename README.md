# Domain Watch

<img src="./assets/icon.svg" alt="domain-watch Icon" width="92" align="right">

[![Build](https://github.com/gabe565/domain-watch/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/domain-watch/actions/workflows/build.yml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/gabe565)](https://artifacthub.io/packages/helm/gabe565/domain-watch)

Get notified about domain changes as they happen.

Domain Watch fetches public whois records for each configured domain on a schedule,
and sends the following notifications:
- Expiration date has passed a threshold
- Status code has changed

Supported notification providers:
- Gotify
- Telegram

## Install
See [Installation](https://github.com/gabe565/domain-watch/wiki/Installation).

## Usage
See [Usage](https://github.com/gabe565/domain-watch/wiki/Usage).
