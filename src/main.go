package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
    "io"
    "io/ioutil"
	"golang.org/x/net/http2"
)

func uploadDownloadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("File Upload Endpoint Hit")

    // Parse our multipart form, 10 << 20 specifies a maximum
    // upload of 10 MB files.
    r.ParseMultipartForm(2000)
    // FormFile returns the first file for the given key `file`
    // it also returns the FileHeader so we can get the Filename,
    // the Header and the size of the file
    uploadFile, handler, err := r.FormFile("file")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
    defer uploadFile.Close()
    fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    fmt.Printf("File Size: %+v\n", handler.Size)
    fmt.Printf("MIME Header: %+v\n", handler.Header)

    // read all of the contents of our uploaded file into a
    // byte array
    fileBytes, err := ioutil.ReadAll(uploadFile)
    if err != nil {
        fmt.Println(err)
    }

    f, err := os.Create("./file")

    if err != nil {
        fmt.Println(err)
    }

    defer f.Close()
    // write this byte array to file
    f.Write(fileBytes)

    // return that we have successfully uploaded our file!
    fmt.Fprintf(w, "Successfully Uploaded File\n")

    downloadFile, err := os.Open("./file")
    if err != nil {
        fmt.Println("error reading file")
        return
    }
    defer downloadFile.Close()
    w.Header().Set("Transfer-Encoding", "chunked")
    io.Copy(w, downloadFile)
}


func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This request is served over %s protocol.", r.Proto)
    http.ServeFile(w, r, "./index.html")
}


func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	addr := net.JoinHostPort("", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %s", addr)

	server := http2.Server{}

    http.HandleFunc("/", handler)
	http.HandleFunc("/uploadDownload", uploadDownloadFile)

	opts := &http2.ServeConnOpts{
		Handler: http.DefaultServeMux,
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
		}
		go server.ServeConn(conn, opts)
	}
}