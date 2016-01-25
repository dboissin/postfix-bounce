package postfixutil

import "fmt"
import "log"
import "os"
import "bufio"
import "strings"
import "time"
import "regexp"

const (
	softDSN string = "5.2.0 5.2.1 5.2.2 5.3.1 5.4.5 5.5.3"
	queueIDPattern string = "\\]:\\s([A-Z0-9]+):"
	logPattern string = "^([A-Za-z]{3}\\s+\\d+ [0-9:]{8}) .*? .*?: ([A-Z0-9]+): to=<(.*?)>, relay=(.*?), delay=(.*?), delays=(.*?), dsn=(.*?), status=(.*?)$"
)

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

func (b *Bounce) IsHard() bool {
	return !strings.Contains(softDSN, b.DSN)
}

func FindBounces(paths *[]string) []Bounce {
	var bounces []Bounce
	bounceQueueID := make(map[string]interface{})
	var logs []string
	queueIDRegex, _ := regexp.Compile(queueIDPattern)
	logRegex, _ := regexp.Compile(logPattern)

	for _, path := range *paths {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalln(err)
			return bounces
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log := scanner.Text()
			if strings.Contains(log, "postfix/bounce") {
				queueID := queueIDRegex.FindStringSubmatch(log)[1]
				bounceQueueID[queueID] = nil
			} else if strings.Contains(log, "status=deferred") || strings.Contains(log, "status=bounced") {
				logs = append(logs, log)
			}
		}
	}

	for _, l := range logs {
		queueID := queueIDRegex.FindStringSubmatch(l)[1]
		if _, exists := bounceQueueID[queueID]; exists {
			v := logRegex.FindStringSubmatch(l)
			if len(v) == 9 {
				bounces = append(bounces, Bounce{ParseDate(v[1]), v[2], v[3], v[4], v[5], v[6], v[7], v[8]})
			} else {
				log.Fatalf("Invalid line : %s", l)
			}
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
