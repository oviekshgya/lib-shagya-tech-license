package license

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ResponseCheck struct {
	MessageService string `json:"message_service"`
	Status         bool   `json:"status"`
}

type FileLicense struct {
	Key         string `json:"key"`
	TypeService string `json:"type_service"`
}

func CheckedFaceRecog() {
	pwd, _ := os.Getwd()
	fileData, err := ioutil.ReadFile(fmt.Sprintf("%s/public/shagya-license.json", pwd))
	if err != nil {
		log.Fatal("Error invalid license")
		return
	}

	var jsonResponse FileLicense
	err = json.Unmarshal(fileData, &jsonResponse)
	if err != nil {
		log.Fatal("Error invalid license 02")
		return
	}

	client := &http.Client{}
	var data map[string]interface{}

	requestBody, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Fatal("invalid license: ", err.Error())
		return
	}
	var jsonStr = []byte(string(requestBody))

	hit, err2 := http.NewRequest("POST", "link", bytes.NewBuffer(jsonStr))
	if err2 != nil {
		return
	}
	hit.Header.Set("Content-Type", "application/json")
	hit.Header.Set("X-LICENSE-Type", "face-recognition-01")

	resp, err := client.Do(hit)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		errBody := Body.Close()
		if errBody != nil {
			fmt.Println("errBody", errBody)
		}
	}(resp.Body)

	errJson := json.NewDecoder(resp.Body).Decode(&data)
	if errJson != nil {
		return
	}

	jsonData, _ := json.Marshal(data)

	var result ResponseCheck
	json.Unmarshal(jsonData, &result)

	if !result.Status {
		log.Println(result.Status)
		return
	}

	return
}
