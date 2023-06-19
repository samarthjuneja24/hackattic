package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

type Solution struct {
	Int               int32   `json:"int"`
	Uint              uint32  `json:"uint"`
	Short             int16   `json:"short"`
	Float             float64 `json:"float"`
	Double            float64 `json:"double"`
	Big_endian_double float64 `json:"big_endian_double"`
}

func main() {
	client := &http.Client{}
	accessToken := os.Getenv("HACKATTIC_ACCESS_TOKEN")

	req, err := http.NewRequest("GET", "https://hackattic.com/challenges/help_me_unpack/problem?access_token="+accessToken, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	jsonResponse := make(map[string]string)
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(jsonResponse["bytes"])
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return
	}

	int32Number := int32(binary.LittleEndian.Uint32(decodedBytes[0:4]))
	uIntNumber := uint32(binary.LittleEndian.Uint32(decodedBytes[4:8]))
	shortNumber := int16(binary.LittleEndian.Uint32(decodedBytes[8:12]))
	floatNumber := float64(binary.LittleEndian.Uint32(decodedBytes[12:16]))
	doubleNumber := float64(binary.LittleEndian.Uint32(decodedBytes[16:24]))
	bigEndianDoubleNumber := math.Float64frombits(binary.BigEndian.Uint64(decodedBytes[24:32]))

	solution := Solution{int32Number, uIntNumber, shortNumber, floatNumber, doubleNumber, bigEndianDoubleNumber}
	jsonSolution, err := json.Marshal(solution)
	if err != nil {
		fmt.Println("Error parsing struct to bytes:", err)
	}
	req, err = http.NewRequest("POST", "https://hackattic.com/challenges/help_me_unpack/problem?access_token="+accessToken, bytes.NewReader(jsonSolution))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}
