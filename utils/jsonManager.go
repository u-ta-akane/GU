package utils

import (
	"GU/refs"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type FileManager interface {
	Read(filename string) interface{}
	Write(filename string, data interface{})
}

type JSONFileManager struct{}

var (
	JobDataSlice []refs.JobData
	JSONFM       = JSONFileManager{}
	//channelId     = map[string]string{}
)

// Write JSONデータをファイルに書き込む
func (j *JSONFileManager) Write(filename string, data interface{}) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("ファイル取得失敗: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			Log(err, "", "Write")
		}
	}(f)

	output, err := json.MarshalIndent(data, "", "\t\t")
	if err != nil {
		log.Fatalf("JSONエンコード失敗: %v", err)
	}

	if _, err := f.Write(output); err != nil {
		log.Fatalf("JSON書き込み失敗: %v", err)
	}
}

// ReadJSON JSONデータをファイルから読み込む
func (j *JSONFileManager) Read(filename string) interface{} {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			if IsCreatedChannel {
				Log(err, "", "JSONFM.Read")
			} else {
				log.Printf(err.Error())
			}
		}
	}(f)
	decoder := json.NewDecoder(f)
	switch filename {
	case "jobData.json":
		var data []refs.JobData
		if err := decoder.Decode(&data); err != nil {
			if IsCreatedChannel {
				Log(err, "", "JSONFM.Read")
			} else {
				log.Printf("JSONデコード失敗: %v", err)
				f, err := os.Create("jobData.json")
				if err != nil {
					log.Fatalf("ファイル取得失敗: %v", err)
				}
				defer func(f *os.File) {
					err := f.Close()
					if err != nil {
						log.Printf(err.Error())
					}
				}(f)

				output, err := json.MarshalIndent(refs.JobData{}, "", "\t\t")
				if err != nil {
					log.Fatalf("JSONエンコード失敗: %v", err)
				}

				if _, err := f.Write(output); err != nil {
					log.Fatalf("JSON書き込み失敗: %v", err)
				}
			}
		}
		fmt.Println(data)
		return data
	case "secrets.json":
		var data refs.SecretData
		if err := decoder.Decode(&data); err != nil {
			if IsCreatedChannel {
				Log(err, "", "JSONFM.Read")
			} else {
				log.Printf("JSONデコード失敗: %v", err)
				f, err := os.Create("secrets.json")
				if err != nil {
					log.Fatalf("ファイル取得失敗: %v", err)
				}
				defer func(f *os.File) {
					err := f.Close()
					if err != nil {
						log.Printf(err.Error())
					}
				}(f)

				output, err := json.MarshalIndent(refs.Secrets, "", "\t\t")
				if err != nil {
					log.Fatalf("JSONエンコード失敗: %v", err)
				}

				if _, err := f.Write(output); err != nil {
					log.Fatalf("JSON書き込み失敗: %v", err)
				}
			}
		}
		fmt.Println(data)
		return data
	case "config.json":
		var data refs.GuildStructure
		if err := decoder.Decode(&data); err != nil {
			if IsCreatedChannel {
				Log(err, "", "JSONFM.Read")
			} else {
				log.Printf("JSONデコード失敗: %v", err)
				f, err := os.Create("config.json")
				if err != nil {
					log.Fatalf("ファイル取得失敗: %v", err)
				}
				defer func(f *os.File) {
					err := f.Close()
					if err != nil {
						log.Printf(err.Error())
					}
				}(f)

				output, err := json.MarshalIndent(refs.Config, "", "\t\t")
				if err != nil {
					log.Fatalf("JSONエンコード失敗: %v", err)
				}

				if _, err := f.Write(output); err != nil {
					log.Fatalf("JSON書き込み失敗: %v", err)
				}
			}
		}
		fmt.Println(data)
		return data
	}
	Log(nil, fmt.Sprintf("Error : Unknown filename \"%s\"", filename), "ReadJSON")
	return nil
}

func (j *JSONFileManager) SearchJobFromJSON(id string) []refs.JobData {
	jd := j.Read("jobData.json").([]refs.JobData)
	var jt []refs.JobData
	for _, j := range jd {
		if j.Id == id {
			jt = append(jt, j)
		}
	}
	return jt
}
