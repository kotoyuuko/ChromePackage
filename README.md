# ChromePackage

这是一个使用 Go 语言编写的获取 Chrome 离线安装包信息的程序。

# 支持平台

 - Windows x86/x64
 - macOS

# 使用方法

运行 Go 程序后会生成 `chrome.json` 文件，所有获取到的信息均存储于此文件。

为了方便使用，本项目使用 Travis-CI 每天自动查询最新的安装包，并将 `chrome.json` 同步到 Github Pages：

    https://kotoyuuko.github.io/ChromePackage/chrome.json

你可以直接将此 API 用于你的程序。

# 协议

[The Unlicense](https://unlicense.org/)
