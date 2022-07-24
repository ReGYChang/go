package debezium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"nexdata/pkg/config"

	"github.com/buger/jsonparser"
)

func SetupConnector() error {
	url := "http://" + config.Debezium.Connector.ConnectAddr + "/connectors/" + config.Debezium.Connector.ConnectName + "/config"

	// 將字串分割 (EX) 10.13.1.156:3306 分割為 10.13.1.156 & 3306
	//splitDBAddr := strings.Split(config.Debezium.Connector.DBAddr, ":")

	// 參數設置 refer https://debezium.io/documentation/reference/1.8/connectors/mysql.html
	reqBody := config.Debezium.Connector

	// 將 struct 轉為 byte
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println(err)
		return err
	}

	payload := bytes.NewReader(jsonBytes)
	client := &http.Client{}

	// 創建req
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// 發送http請求
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(res.Body)

	// 回傳結果讀取
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	// 取出回傳結果 json 字串中的 name 參數值
	n, err := jsonparser.GetString(body, "name")
	if err != nil || n != config.Debezium.Connector.ConnectName {
		fmt.Println("Set Debezium Connect fail")
		return err
	}

	return nil
}
