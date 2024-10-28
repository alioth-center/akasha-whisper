# Akasha Whisper 

![Go Version](https://img.shields.io/github/go-mod/go-version/alioth-center/akasha-whisper)
![Release](https://img.shields.io/github/v/release/alioth-center/akasha-whisper)
![Go Report Card](https://goreportcard.com/badge/github.com/alioth-center/akasha-whisper)
![GitHub Actions](https://img.shields.io/github/actions/workflow/status/alioth-center/akasha-whisper/build-docker.yml?branch=main)
![License](https://img.shields.io/github/license/alioth-center/akasha-whisper)

## Summary

`akasha-whisper` is a Golang project offering a unified, user-friendly API for integrating and interacting with multiple AI models.

While similar to [one-api](https://github.com/songquanpeng/one-api), `akasha-whisper` provides additional features:

- Dynamic client weight configuration and intelligent load balancing.
- Support for multiple providers for each model, allowing integration of various sources like OpenAI and Azure, with automatic provider selection based on predefined weight and pricing configurations.

## Document

For detailed documentation, visit the [Akasha-whisper](https://docs.alioth.center/akasha-whisper.html)  documentation.

## Compatible APIs

|     Name      |                 Description                  |           URL            | Status |
|:-------------:|:--------------------------------------------:|:------------------------:|:------:|
|     Chat      |       complete chat with given prompt        |  `v1/chat/completions`   |   ✅    |
|     Model     |            list available models             |       `/v1/models`       |   ✅    |
|  Embeddings   |        get embeddings for given text         |     `/v1/embeddings`     |   ✅    |
|     Image     |       generate image with given prompt       | `v1/images/generations`  |  WIP   |
|    Speech     |    generate speech audio from given text     |    `v1/audio/speech`     |  WIP   |
| Transcription |     generate text from given audio file      | `v1/audio/transcription` |  WIP   |

## Support Models

The following models and providers have been thoroughly tested by the development team and confirmed to function reliably.

|  Name   | Provider |                           Home Page                            |                       BaseURL                       |
|:-------:|:--------:|:--------------------------------------------------------------:|:---------------------------------------------------:|
|   GPT   |  OpenAI  |                  [OpenAI](https://openai.com)                  |             `https://api.openai.com/v1`             |
|   GLM   | ZhipuAI  |              [ZhipuAI](https://open.bigmodel.cn)               |       `https://open.bigmodel.cn/api/paas/v4`        |
|  qwen   | Alibaba  |              [Tongyi](https://tongyi.aliyun.com)               | `https://dashscope.aliyuncs.com/compatible-mode/v1` |
| hunyuan | Tencent  | [Hunyuan](https://cloud.tencent.com/act/pro/Hunyuan-promotion) |     `https://api.hunyuan.cloud.tencent.com/v1`      |              

## Thanks

Thanks to JetBrains for providing [Open Source development license(s)](https://www.jetbrains.com/community/opensource/#support) for this project.

## Contributors

![Contributors](https://contrib.rocks/image?repo=alioth-center/akasha-whisper&max=1000)

