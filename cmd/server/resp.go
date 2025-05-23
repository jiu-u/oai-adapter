package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

func getContentType(header http.Header) string {
	if contentType := header.Get("Content-Type"); contentType != "" {
		return contentType
	}
	return "application/json"
}

func HandleOAIResponse(w http.ResponseWriter, req *http.Request, responseBody io.ReadCloser, respHeader http.Header) {
	defer responseBody.Close()

	//multiWriter := io.MultiWriter(w, os.Stdout)
	//for k, v := range respHeader {
	//	w.Header().Set(k, v[0])
	//}
	//
	//_, err := io.Copy(multiWriter, responseBody)
	//if err != nil {
	//	fmt.Println("error writing response:", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	for k, v := range respHeader {
		w.Header().Set(k, v[0])
	}

	//流式读取和写入响应
	reader := bufio.NewReader(responseBody)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				http.Error(w, fmt.Errorf("error writing response: %w", writeErr).Error(), http.StatusInternalServerError)
				return
			}
		}
		if err != nil {
			if err == io.EOF {
				//w.Write([]byte("data: [DONE]\n\n")) // 发送结束标记
				break
			}
			http.Error(w, fmt.Errorf("error reading response: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
	//time.Sleep(time.Second * 1)
}
