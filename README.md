# 🚀 优选 IP 自动采集器 (Go)

[![Update IP List](https://github.com/camel52zhang/yxip/actions/workflows/main.yml/badge.svg)](https://github.com/camel52zhang/yxip/actions/workflows/main.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> 每天自动从多个数据源采集优选 IPv4 地址，聚合去重后输出 `ip.txt`，供科学上网、CDN 加速等场景使用。

---

## 📑 目录

- [项目背景](#-项目背景)
- [架构设计](#-架构设计)
- [数据采集流程](#-数据采集流程)
- [采集源列表](#-采集源列表)
- [项目结构](#-项目结构)
- [GitHub Actions 部署](#-github-actions-部署)
- [本地运行](#-本地运行)
- [输出格式](#-输出格式)
- [核心依赖](#-核心依赖)
- [常见问题](#-常见问题)
- [贡献指南](#-贡献指南)

---

## 📖 项目背景

### 什么是优选 IP？

**优选 IP** 是指经过测试筛选出的、在当前网络环境下表现最优的IP地址。

### 解决的问题

- 🎯 **自动聚合**：从多个可信数据源抓取优选 IP，避免人工筛选
- 🔄 **定时更新**：通过 GitHub Actions 每天 0 点（北京时间）自动更新
- 🧹 **去重处理**：自动去除重复 IP，确保输出干净可用
- ⚡ **高性能**：Go 语言实现，编译型语言，启动快、内存占用低

---

## 🏗️ 架构设计
