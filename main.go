package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)


//Document struct
type Document struct {
	ID   string
	Name string
	Size int64
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{key}", getDocumentsById).Methods("GET")
	router.HandleFunc("/documents", saveFileByData).Methods("POST")
	router.HandleFunc("/documents/{key}", deleteDocumentsById).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":9000", router))
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var docs []Document
	fileDir := "./FileToTest/"
	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}

	for _, file := range fileInfos {
		fileHash, err := hashFileToMD5CheckSum(fileDir + "/" + file.Name())
		if err == nil {
			docs = append(docs, Document{ID: fileHash, Name: file.Name(), Size: file.Size()})
			fmt.Printf("MD5: %s  -  Name: %s  -  Size: %d Kbps \n", fileHash, file.Name(), file.Size())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}


func getDocumentsById(w http.ResponseWriter, r *http.Request) {
	var idFile = r.URL.Path[11:]
	var docs Document
	fileDir := "./FileToTest/"
	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}

	for _, file := range fileInfos {
		fileHash, err := hashFileToMD5CheckSum(fileDir + "/" + file.Name())
		if err == nil {
			if fileHash == idFile {
				docs = Document{ID: fileHash, Name: file.Name(), Size: file.Size()}
				fmt.Printf("MD5: %s  -  Name: %s  -  Size: %d Kbps \n", fileHash, file.Name(), file.Size())
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(docs)
				return
			}
		}
	}
	http.Error(w, "Forbidden", http.StatusNotFound)
}


func deleteDocumentsById(w http.ResponseWriter, r *http.Request){
	var idFile = r.URL.Path[11:]
	fileDir := "./FileToTest/"
	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}

	for _, file := range fileInfos {
		fileHash, err := hashFileToMD5CheckSum(fileDir + "/" + file.Name())
		if err == nil {
			if fileHash == idFile {
				var nameFileToRemove = file.Name()
				os.Remove(fileDir + nameFileToRemove)
				http.Error(w, "Forbidden", http.StatusOK)
				return
			}
		}
	}
	http.Error(w, "Forbidden", http.StatusNotFound)
}


func saveFileByData(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&Buf, file)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	contents := Buf.String()
	fmt.Println(contents)

	fo,_ := os.Create("./FileToTest/" + header.Filename)

	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := Buf.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
			http.Error(w, "Forbidden", http.StatusNotFound)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
			http.Error(w, "Forbidden", http.StatusNotFound)
		}
	}

	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	Buf.Reset()
	// do something else
	// etc write header
	http.Error(w, "Forbidden", http.StatusOK)
	return
}


func saveDocumentsByBodyData(w http.ResponseWriter, r *http.Request) {
	var idFile = r.URL.Path[11:]
	var docs Document
	fileDir := "./FileToTest/"
	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println("Error in accessing directory:", err)
	}

	for _, file := range fileInfos {
		fileHash, err := hashFileToMD5CheckSum(fileDir + "/" + file.Name())
		if err == nil {
			if fileHash == idFile {

				break
			}else{
				http.Error(w, "Forbidden", http.StatusNotFound)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}


func hashFileToMD5CheckSum(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}
