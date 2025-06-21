
# `resto`

[![Code Quality](https://img.shields.io/badge/Code_Quality-A+-brightgreen?style=for-the-badge&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/lucasepe/resto)

> A minimalist CLI REST client that calls APIs, waits for conditions, and retries intelligently.

## Overview

`resto` is a tool that allows you to make HTTP calls with retry capability. 

While it can be used for general retry scenarios, it’s especially useful when you need to ensure that a REST API returning JSON objects has marked those objects with a desired condition or status.

`resto` lets you retry requests until a specified jq condition evaluates to true. 

This feature is particularly handy when working with objects managed by Kubernetes APIs, for example, but it’s broadly applicable to any REST API that accepts an operation and then updates the resource’s status accordingly.

Makes scripting and automation of REST API calls simpler and more reliable
in CI/CD pipelines and development workflows.


## 🔧 Usage

```sh
resto [FLAGS] URL
```

For complete help including all flags, supported environment variables, and usage examples, type:

```sh
resto --help
```

## 👍 Support

All tools are completely free to use, with every feature fully unlocked and accessible.

If you find one or more of these tool helpful, please consider supporting its development with a donation.

Your contribution, no matter the amount, helps cover the time and effort dedicated to creating and maintaining these tools, ensuring they remain free and receive continuous improvements.

Every bit of support makes a meaningful difference and allows me to focus on building more tools that solve real-world challenges.

Thank you for your generosity and for being part of this journey!

[![Donate with PayPal](https://img.shields.io/badge/💸-Tip%20me%20on%20PayPal-0070ba?style=for-the-badge&logo=paypal&logoColor=white)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=FV575PVWGXZBY&source=url)


## 🛠️ How To Install

### Download the latest binaries from the [releases page](https://github.com/lucasepe/resto/releases/latest):

- [macOS](https://github.com/lucasepe/resto/releases/latest)
- [Windows](https://github.com/lucasepe/resto/releases/latest)
- [Linux (arm64)](https://github.com/lucasepe/resto/releases/latest)
- [Linux (amd64)](https://github.com/lucasepe/resto/releases/latest)

### Using a Package Manager

» macOS » [Homebrew](https://brew.sh/)

```sh
brew tap lucasepe/cli-tools
brew install resto
```
