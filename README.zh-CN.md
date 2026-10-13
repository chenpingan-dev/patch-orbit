# Patch Orbit · Go 改动影响地图

**改了公共包，哪些地方值得复查？沿 import 链给出依据**

[在线体验](https://chenpingan-dev.github.io/patch-orbit/) · [English](README.md)

![改动影响地图](assets/screenshot.png)

面向维护 Go 模块的开发者和代码评审者。对比两个已提交的版本，找到改动包及其直接、间接依赖者，点开节点就能看到最短导入链和对应文件行号

## 使用

需要 Git 和 Go 1.23+，也可以从 [Releases](https://github.com/chenpingan-dev/patch-orbit/releases) 下载编译好的程序

```sh
go install github.com/chenpingan-dev/patch-orbit/cmd/orbit@v0.1.0
orbit -repo 你的Go仓库路径 -base HEAD~1 -head HEAD -html impact.html -json impact.json
```

浏览器打开 `impact.html`，可以搜索包、查看依赖依据、复制测试命令、下载复查清单。已有报告需要 `-force` 才能覆盖。两个版本都要有根目录 `go.mod`

## 范围

依赖图合并改动前后的 import，所以删除的包和移除的依赖不会直接消失。包含测试文件和全部构建条件；依赖清单变化时把整个模块纳入复查

当前按单模块的包导入关系分析，不覆盖运行时调用、反射、嵌入资源、外部模块和未提交改动。未匹配包不代表绝对安全。报告生成过程不会执行源码、测试、构建或远程拉取

跳过 vendor、隐藏目录、下划线目录、testdata、嵌套模块和符号链接。无效 import、缺失本地包和过大文件会留下扫描提示。每个快照最多选择 3,000 个文件 / 32 MiB，超过 1 MiB 的单个文件会跳过。地图最多展示 80 个节点，完整结果仍保留在列表和 JSON

在线示例来自两个真实的本地合成 Git 提交，仅用来演示分析过程，不包含业务源码。界面支持中英文，使用微软雅黑，可保存单文件离线报告

```sh
go test ./...
go vet ./...
go run ./cmd/demo
```

独立、AI 辅助开发的 MIT 开源项目
