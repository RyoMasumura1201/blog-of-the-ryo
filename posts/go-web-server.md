# GoでシンプルなWebサーバーを作る
date: 2024-03-10

Go の標準ライブラリだけでシンプルな HTTP サーバーを作る方法を紹介します。

## 最小構成

まず一番シンプルな Hello World サーバーから始めましょう。

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })
    http.ListenAndServe(":8080", nil)
}
```

`go run main.go` で起動し、ブラウザで `http://localhost:8080` にアクセスすると "Hello, World!" が表示されます。

## ルーティング

複数のエンドポイントを登録する例です。

```go
package main

import (
    "encoding/json"
    "net/http"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /", homeHandler)
    mux.HandleFunc("GET /users/{id}", getUserHandler)

    http.ListenAndServe(":8080", mux)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Home page"))
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    user := User{Name: "Ryo", Age: 25}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{
        "id":   id,
        "user": user,
    })
}
```

> Go 1.22 以降では `mux.HandleFunc("GET /path", ...)` のようにメソッドを指定できます。

## ミドルウェア

ログを出力するシンプルなミドルウェアです。

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

## まとめ

| 機能 | パッケージ |
|------|-----------|
| HTTP サーバー | `net/http` |
| JSON | `encoding/json` |
| テンプレート | `html/template` |
| ログ | `log` / `log/slog` |

Go の標準ライブラリは充実しているので、シンプルな API なら外部ライブラリなしで作れます。
