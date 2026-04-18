package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Post struct {
	Slug      string
	Title     string
	Date      time.Time
	Content   template.HTML
	RawSource string
}

type PageData struct {
	Posts []*Post
	Post  *Post
}

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithGuessLanguage(true),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
}

// parsePost reads a markdown file and parses front matter (title, date).
// Front matter is expected as the first lines starting with "# Title" and "date: YYYY-MM-DD".
func parsePost(slug, content string) (*Post, error) {
	lines := strings.Split(content, "\n")
	title := slug
	date := time.Time{}
	bodyStart := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") && title == slug {
			title = strings.TrimPrefix(trimmed, "# ")
			bodyStart = i + 1
		} else if strings.HasPrefix(trimmed, "date:") {
			raw := strings.TrimSpace(strings.TrimPrefix(trimmed, "date:"))
			if t, err := time.Parse("2006-01-02", raw); err == nil {
				date = t
			}
			if bodyStart == i {
				bodyStart = i + 1
			}
		} else if trimmed != "" && i > 0 {
			if bodyStart == 0 {
				bodyStart = i
			}
			break
		}
	}

	body := strings.Join(lines[bodyStart:], "\n")
	var buf bytes.Buffer
	if err := md.Convert([]byte(body), &buf); err != nil {
		return nil, err
	}

	return &Post{
		Slug:      slug,
		Title:     title,
		Date:      date,
		Content:   template.HTML(buf.String()),
		RawSource: content,
	}, nil
}

func loadPosts(dir string) ([]*Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		slug := strings.TrimSuffix(e.Name(), ".md")
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			log.Printf("failed to read %s: %v", e.Name(), err)
			continue
		}
		post, err := parsePost(slug, string(data))
		if err != nil {
			log.Printf("failed to parse %s: %v", e.Name(), err)
			continue
		}
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})
	return posts, nil
}

func chromaCSS() template.CSS {
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}
	// Return empty; chroma inlines styles when using WithClasses(false) (default)
	return ""
}

func main() {
	tmplFuncs := template.FuncMap{
		"chromaCSS":  chromaCSS,
		"formatDate": func(t time.Time) string { return t.Format("2006-01-02") },
	}

	indexTmpl := template.Must(template.New("index.html").Funcs(tmplFuncs).ParseFiles(
		"templates/index.html",
		"templates/base.html",
	))
	postTmpl := template.Must(template.New("post.html").Funcs(tmplFuncs).ParseFiles(
		"templates/post.html",
		"templates/base.html",
	))

	// Static assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	// Index page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		posts, err := loadPosts("posts")
		if err != nil {
			http.Error(w, "failed to load posts", http.StatusInternalServerError)
			return
		}
		if err := indexTmpl.Execute(w, PageData{Posts: posts}); err != nil {
			log.Printf("index template error: %v", err)
		}
	})

	// Individual post page
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/posts/")
		slug = strings.TrimSuffix(slug, "/")
		if slug == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		data, err := os.ReadFile(filepath.Join("posts", slug+".md"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		post, err := parsePost(slug, string(data))
		if err != nil {
			http.Error(w, "failed to parse post", http.StatusInternalServerError)
			return
		}
		if err := postTmpl.Execute(w, PageData{Post: post}); err != nil {
			log.Printf("post template error: %v", err)
		}
	})

	addr := ":8080"
	log.Printf("starting blog server on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
