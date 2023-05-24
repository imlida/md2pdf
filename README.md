# MD to PDF Converter

这是一个将Markdown文件合并并转换为PDF的Go程序。

## 功能

1. 遍历指定的目录及其子目录，查找所有的.md和.mdx文件。
2. 将找到的文件合并为一个.md文件，并在字符数超过指定最大字符数时转换为.pdf文件。

## 安装依赖

这个程序需要 `wkhtmltopdf` 来将HTML文件转换为PDF。在运行程序之前，请确保您已经安装了 `wkhtmltopdf`。您可以在以下网址下载和安装 `wkhtmltopdf`：

https://wkhtmltopdf.org/downloads.html


## 用法

```bash
go run main.go -dir=<扫描的根目录> -maxCharCount=<最大字符数>
```

-dir: 指定扫描的根目录，默认值为当前目录(".").

-maxCharCount: 合并成pdf的最大字符数，默认值为800000。

例如，以下命令将扫描 "myDirectory" 目录并设置最大字符数为 1000000:

```bash
go run main.go -dir=myDirectory -maxCharCount=1000000
```