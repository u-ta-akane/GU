package utils

import (
	"GU/refs"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type FileManager interface {
	Read(filename string) interface{}
	Write(filename string, data interface{})
}

type JSONFileManager struct {
	muWrite sync.Mutex
	muRead  sync.Mutex
}

var (
	JobDataSlice []refs.JobData
	JSONFM       = JSONFileManager{}
	//channelId     = map[string]string{}
)

// Write JSONデータをファイルに書き込む
func (j *JSONFileManager) Write(filename string, data interface{}) error {
	j.muWrite.Lock()
	f, err := os.Create(filename)
	if err != nil {
		j.muWrite.Unlock()
		Log(err, "", "Write")
		return err
	}
	defer func(f *os.File) {
		errFClose := f.Close()
		if errFClose != nil {
			Log(err, "", "Write")
		}
		j.muWrite.Unlock()
	}(f)
	if err != nil {
		log.Print("ファイル取得失敗: %v", err)
	}

	output, errIndent := json.MarshalIndent(data, "", "\t\t")
	if errIndent != nil {
		log.Printf("JSONエンコード失敗: %v", err)
	}

	if _, errFWrite := f.Write(output); errFWrite != nil {
		log.Printf("JSON書き込み失敗: %v", err)
	}
	return err
}

// ReadJSON JSONデータをファイルから読み込む
func (j *JSONFileManager) Read(filename string) interface{} {
	j.muRead.Lock()
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("ファイル取得失敗: %v", err)
		j.muRead.Unlock()
		return err
	}
	defer func(f *os.File) {
		errFClose := f.Close()
		if errFClose != nil {
			if IsCreatedChannel {
				Log(errFClose, "", "JSONFM.Read")
			} else {
				log.Printf(errFClose.Error())
			}
		}
		j.muRead.Unlock()
	}(f)
	decoder := json.NewDecoder(f)
	switch filename {
	case "jobData.json":
		var data []refs.JobData
		if errFilenameSwitch := decoder.Decode(&data); errFilenameSwitch != nil {
			if IsCreatedChannel {
				Log(errFilenameSwitch, "", "JSONFM.Read")
			} else {
				log.Printf("JSONデコード失敗: %v", err)
				ff, errNewCreate := os.Create("jobData.json")
				if errNewCreate != nil {
					log.Fatalf("ファイル取得失敗: %v", errNewCreate)
				}
				defer func(ff *os.File) {
					errFFClose := ff.Close()
					if errFFClose != nil {
						log.Printf(errFFClose.Error())
					}
				}(ff)

				output, errIndent := json.MarshalIndent(refs.JobData{}, "", "\t\t")
				if errIndent != nil {
					log.Fatalf("JSONエンコード失敗: %v", errIndent)
				}

				if _, err := ff.Write(output); err != nil {
					log.Fatalf("JSON書き込み失敗: %v", err)
				}
			}
		}
		fmt.Println(data)
		return data
	case "secrets.json":
		var data refs.SecretData
		if errFilenameSwitch := decoder.Decode(&data); errFilenameSwitch != nil {
			if IsCreatedChannel {
				Log(errFilenameSwitch, "", "JSONFM.Read")
			} else {
				log.Printf("JSONデコード失敗: %v", errFilenameSwitch)
				ff, errFFNewCreate := os.Create("secrets.json")
				if errFFNewCreate != nil {
					log.Fatalf("ファイル取得失敗: %v", errFFNewCreate)
				}
				defer func(ff *os.File) {
					err := ff.Close()
					if err != nil {
						log.Printf(err.Error())
					}
				}(ff)

				output, err := json.MarshalIndent(refs.Secrets, "", "\t\t")
				if err != nil {
					log.Fatalf("JSONエンコード失敗: %v", err)
				}

				if _, err := ff.Write(output); err != nil {
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
				ff, err := os.Create("config.json")
				if err != nil {
					log.Fatalf("ファイル取得失敗: %v", err)
				}
				defer func(ff *os.File) {
					err := ff.Close()
					if err != nil {
						log.Printf(err.Error())
					}
				}(ff)

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
