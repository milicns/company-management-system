package api

import (
	"io"
	"net/http"
)

func Authorize(protectedEndpoint http.HandlerFunc, userSvcAddress string) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		req, err := http.NewRequest("GET", userSvcAddress, nil)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		req.Header.Add("Authorization", request.Header.Get("Authorization"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(writer, string(body), http.StatusUnauthorized)
			return
		}

		protectedEndpoint(writer, request)
	})
}
