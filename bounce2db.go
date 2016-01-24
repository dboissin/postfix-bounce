package main

import "fmt"
import "log"
import "os"
import "bufio"
import "strings"
import "flag"
import "time"
import "regexp"
import "gopkg.in/mgo.v2"

type Bounce struct {
	// ID string
	Date time.Time
	QueueID string
	To string
	Relay string
	Delay string
	Delays string
	DSN string
	Status string
}

func FindBounces(path string) []Bounce {
	var bounces []Bounce
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
		return bounces
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	bounceQueueID := make(map[string]interface{})
	var logs []string
	queueIDRegex, _ := regexp.Compile("\\]:\\s([A-Z0-9]+):")
	for scanner.Scan() {
		log := scanner.Text()
		if strings.Contains(log, "postfix/bounce") {
			queueID := queueIDRegex.FindStringSubmatch(log)[1]
			bounceQueueID[queueID] = nil
		} else if strings.Contains(log, "status=deferred") || strings.Contains(log, "status=bounced") {
			logs = append(logs, log)
		}
	}
	logRegex, _ := regexp.Compile("^([A-Za-z]{3} \\d+ [0-9:]{8}) .*? .*?: ([A-Z0-9]+): to=<(.*?)>, relay=(.*?), delay=(.*?), delays=(.*?), dsn=(.*?), status=(.*?)$")
	for _, log := range logs {
		queueID := queueIDRegex.FindStringSubmatch(log)[1]
		if _, exists := bounceQueueID[queueID]; exists {
			v := logRegex.FindStringSubmatch(log)
			bounces = append(bounces, Bounce{ParseDate(v[1]), v[2], v[3], v[4], v[5], v[6], v[7], v[8]})
		}
	}
	return bounces
}

func ParseDate(Date string) time.Time {
	value := fmt.Sprintf("%v %v", time.Now().Year(), Date)
	rtime, err := time.Parse("2006 Jan 2 15:04:05", value)
	if err != nil {
		log.Fatalln("Error Parsing Time Format: ", value)
		return time.Now()
	}
	return rtime
}

func main() {
	flag.Parse()
	bounces := FindBounces(flag.Arg(0))

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("one_gridfs").C("bounces")
	for _, bounce := range bounces {
		err = c.Insert(&bounce)
		if err != nil {
			log.Fatal(err)
		}
	}

}
