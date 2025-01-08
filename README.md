# Localhook

**Webhook Testing Made Simple**

Localhook is a lightweight tool built in Go to receive and debug webhooks locally. Perfect for developers who need to test webhooks during development.

---

## Features

- 🚀 **Lightweight**: Built in Go for fast and efficient performance.
- 🖥️ **Cross-Platform**: Supports macOS, Linux, and Windows.
- 📝 **Easy to Use**: Simply run the binary and start receiving webhooks.
- 🔍 **Debugging Made Easy**: Logs incoming requests for easy inspection.

---

## Download Pre-Built Binaries

You can download pre-built binaries for your platform from the table below:

| Platform      | Architecture   | Download Link                                                                            |
| ------------- | -------------- | ---------------------------------------------------------------------------------------- |
| macOS (Intel) | 64-bit         | [Download](https://github.com/ahmetsabri/localhook/releases/download/V1/mac_m1.zip)      |
| macOS (ARM)   | 64-bit (M1/M2) | [Download](https://github.com/ahmetsabri/localhook/releases/download/V1/mac_m1.zip)      |
| Linux (Intel) | 64-bit         | [Download](https://github.com/ahmetsabri/localhook/releases/download/V1/linux_intel.zip) |
| Linux (ARM)   | 64-bit         | [Download](https://github.com/ahmetsabri/localhook/releases/download/V1/linux_arm.zip)   |
| Windows       | 64-bit         | [Download](https://github.com/ahmetsabri/localhook/releases/download/V1/windows.zip)     |

---

## Usage

1. Download the appropriate pre-built binary for your platform from the [Downloads](#download-pre-built-binaries) section.
2. Unzip the downloaded file and navigate to the unzipped folder:
   ```bash
   unzip <platform>.zip
   cd <platform>
   sudo mv localhook /usr/local/bin
3. Reopen terminal and run the localhook command:
   ```bash
   localhook -f <your local URL>
5. Use the provided `webhook url` in any api and go on ... 🚀
