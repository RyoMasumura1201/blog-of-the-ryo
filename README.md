# Blog of the Ryo

GoとMarkdownで作ったシンプルなブログサイトです。

## 必要なもの

- Go 1.22 以上

## セットアップ

```bash
git clone https://github.com/ryomasumura1201/blog-of-the-ryo.git
cd blog-of-the-ryo
go mod download
```

## 起動

```bash
go run main.go
```

ブラウザで http://localhost:8080 にアクセスしてください。

## 記事の追加

`posts/` ディレクトリに `.md` ファイルを追加するだけで記事が反映されます。

```
posts/my-new-post.md
```

ファイルの先頭にタイトルと日付を記述します。

```markdown
# 記事タイトル
date: 2024-04-01

本文をここに書きます。
```

コードブロックはシンタックスハイライト付きでレンダリングされます。

## 画像の使い方

`images/` ディレクトリに画像ファイルを置き、記事中で `/images/ファイル名` を参照します。

```bash
images/my-photo.png
```

```markdown
![説明文](/images/my-photo.png)
```

## Docker

GHCRからイメージを取得して起動できます。

```bash
docker pull ghcr.io/ryomasumura1201/blog-of-the-ryo:latest
docker run -p 8080:8080 ghcr.io/ryomasumura1201/blog-of-the-ryo:latest
```

ブラウザで http://localhost:8080 にアクセスしてください。

## プロジェクト構成

```
.
├── main.go          # HTTPサーバー
├── posts/           # Markdownの記事ファイル
├── images/          # 画像ファイル
├── templates/       # HTMLテンプレート
│   ├── base.html
│   ├── index.html
│   └── post.html
└── static/
    └── style.css    # スタイルシート
```
