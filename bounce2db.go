package main

import "log"
import "os"
import "flag"
import "./postfixutil"
import "encoding/json"
import "gopkg.in/mgo.v2"

const (
	BounceCollection string = "bounces"
)

type Config struct {
	DBHost string
	DBName string
}

func ReadConfig(configFilePath string) Config {
	file, _ := os.Open(configFilePath)
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "conf", "conf.json", "Configuration file")
	flag.Parse()
	filePaths := flag.Args()

	config := ReadConfig(configPath)

	session, err := mgo.Dial(config.DBHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.DBName).C(BounceCollection)

	bounces := postfixutil.FindBounces(&filePaths)
	for _, bounce := range bounces {
		if bounce.IsHard() {
			err = c.Insert(&bounce)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
