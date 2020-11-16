package notice

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	eventName = os.Getenv("IFTTT_EVENT_NAME")
	key       = os.Getenv("IFTTT_KEY")
)

var api = fmt.Sprintf("https://maker.ifttt.com/trigger/%v/with/key/%v", eventName, key)

func IftttNotice(message string) error {
	msg := fmt.Sprintf("{\"value1\":%q}", message)
	buf := bytes.NewBufferString(msg)
	fmt.Println("sending ifttt message:", buf.String())

	if eventName != "" && key != "" {
		resp, err := http.Post(api, "application/json", buf)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(result))
		return nil
	}
	return errors.New("Failed to get IFTTT event name and key from environment")
}
