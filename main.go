package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var www = "./www"
var addr = "0.0.0.0:8080"

const configPostfix = ".wdv.json"

var absRoot = ""

func main() {

	www = getEnv("MDVIEWER_WWW", www)
	addr = getEnv("MDVIEWER_ADDRESS", addr)

	abs, err := filepath.Abs(www)
	if err != nil {
		fmt.Println("Failed to get abs path of files folder")
		return
	}
	absRoot = abs

	http.HandleFunc("/ls", ls)
	http.Handle("/", http.FileServer(http.Dir(www)))

	fmt.Println("Server started at http://" + addr)
	http.ListenAndServe(addr, nil)
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ls(w http.ResponseWriter, r *http.Request) {

	var docsFolder = filepath.Join(absRoot, "docs")
	_, folders, err := scanDirectory(docsFolder, ".md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var response = Response{folders, getExtraConfig(docsFolder, "config"+configPostfix)}
	json.NewEncoder(w).Encode(response)
}

type Response struct {
	Docs   FileNode        `json:"docs"`
	Config json.RawMessage `json:"config"`
}

type FileNode struct {
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	IsFolder    bool            `json:"isFolder"`
	Children    []FileNode      `json:"children,omitempty"`
	ExtraConfig json.RawMessage `json:"extra_config,omitempty"`
}

func scanDirectory(path string, filterSuffix string) (map[string]string, FileNode, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, FileNode{}, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, FileNode{}, err
	}

	if !fileInfo.IsDir() {
		return nil, FileNode{}, errors.New("not a folder")
	}

	relPath, err := filepath.Rel(absRoot, path)
	if err != nil {
		return nil, FileNode{}, err
	}

	files := make(map[string]string)
	node := FileNode{
		Name:        fileInfo.Name(),
		Path:        filepath.ToSlash(relPath),
		IsFolder:    fileInfo.IsDir(),
		ExtraConfig: getExtraConfig(path, configPostfix),
	}

	children, err := file.Readdir(-1)
	if err != nil {
		return nil, FileNode{}, err
	}

	for _, childInfo := range children {
		childPath := filepath.Join(path, childInfo.Name())
		if !childInfo.IsDir() {
			if filterSuffix == "" || strings.HasSuffix(childPath, filterSuffix) {
				relPath, err := filepath.Rel(absRoot, childPath)
				if err != nil {
					return nil, FileNode{}, err
				}

				childNode := FileNode{
					Name:        childInfo.Name(),
					Path:        filepath.ToSlash(relPath),
					IsFolder:    childInfo.IsDir(),
					ExtraConfig: getExtraConfig(path, childInfo.Name()+configPostfix),
				}
				files[filepath.ToSlash(relPath)] = childPath
				node.Children = append(node.Children, childNode)
			}
		} else {
			subfiles, childNode, err := scanDirectory(childPath, filterSuffix)
			if err != nil {
				continue // or handle the error as you prefer
			}
			if len(subfiles) > 0 {
				maps.Copy(files, subfiles)
				node.Children = append(node.Children, childNode)
			}
		}

	}

	return files, node, nil
}

func getExtraConfig(path string, configName string) []byte {
	var filename = filepath.Join(path, configName)
	if _, err := os.Stat(filename); err == nil {
		// File exists, read its contents
		contents, err := os.ReadFile(filename)
		if err == nil {
			return contents
		}
	}
	return nil
}
