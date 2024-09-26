package web

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path"

	"github.com/giuliano-macedo/go-pdfjam/internal/core/pdfjoin"
)

type Server struct {
	outDir string
}

//go:embed static
var staticFs embed.FS

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s Server) upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4_000_000_000)
	if err != nil {
		fmt.Println("failed parsing form", err)
		w.WriteHeader(500)
		return
	}

	dir, err := os.MkdirTemp("", "go_pdfjam")
	if err != nil {
		fmt.Println("failed creating temp dir", err)
		w.WriteHeader(500)
		return
	}
	defer os.RemoveAll(dir)

	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		fmt.Println("no files", err)
		w.WriteHeader(400)
		return
	}

	inFileNames := make([]string, 0, len(files))
	buff := make([]byte, 4096)
	totalBytesWritten := int64(0)
	for i, header := range files {
		f, err := header.Open()
		if err != nil {
			fmt.Println("failed opening form file", err)
			w.WriteHeader(500)
			return
		}

		outFileName := path.Join(dir, fmt.Sprintf("%03d.pdf", i))
		outFile, err := os.Create(outFileName)
		if err != nil {
			fmt.Println("failed creating temp file", err)
			w.WriteHeader(500)
			return
		}
		defer outFile.Close()

		bytesWritten, err := io.CopyBuffer(outFile, f, buff)
		if err != nil {
			fmt.Println("failed copying temp file", err)
			w.WriteHeader(500)
			return
		}

		totalBytesWritten += bytesWritten
		inFileNames = append(inFileNames, outFileName)
	}

	outFile := fmt.Sprintf("%s.pdf", RandStringBytes(12))

	err = pdfjoin.Join(inFileNames, path.Join(s.outDir, outFile))

	if err != nil {
		fmt.Println("failed executing command", err)
		w.WriteHeader(500)
		return
	}
	fmt.Printf("Success! joined %d files %.2f MBs\n", len(files), float64(totalBytesWritten)/1_000_000)

	w.WriteHeader(200)
	fmt.Fprintf(w, "<div class=\"answer\"><a href=\"/file/%s\" download> Download (%d files joined) </a></div>", outFile, len(files))
}

func NewServer() (Server, error) {
	outDir, err := os.MkdirTemp("", "go_pdfjam_output")
	return Server{outDir: outDir}, err
}

func (s Server) Run(addr string) error {
	staticFs, err := fs.Sub(staticFs, "static")
	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.FS(staticFs)))
	http.Handle("/file/", http.StripPrefix("/file", http.FileServer(http.Dir(s.outDir))))
	http.Handle("/upload", http.HandlerFunc(s.upload))

	fmt.Println("listening...")
	defer os.RemoveAll(s.outDir)
	return http.ListenAndServe(addr, nil)
}
