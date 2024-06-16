# ベースイメージとしてGoの公式イメージを使用
FROM golang:1.22

# 作業ディレクトリを設定
WORKDIR /E-Commerce-Chat-Microservice

# bashをインストール
RUN apt-get update && apt-get install -y bash

# go.mod と go.sum をコピー
COPY ./src/go.mod ./src/go.sum ./

# 依存関係をダウンロードして go.sum を更新
RUN go mod tidy
RUN go mod download

# アプリケーションのソースコードをコピー
COPY ./src .

# 依存関係のチェックを含むビルドステップを追加
RUN go build -o /dev/null ./...
