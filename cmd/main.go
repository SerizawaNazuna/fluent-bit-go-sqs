package main

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)
import (
	"encoding/json"
)

var url string
var svc *sqs.SQS

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(def, "sqsgo", "sqsgo")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	fmt.Println("init")
	url = output.FLBPluginConfigKey(plugin, "Url")
	if url == "" {
		return output.FLB_ERROR
	}
	config := aws.Config{
		Region: aws.String("ap-northeast-1"),
	}
	sess := session.Must(session.NewSession(&config))
	svc = sqs.New(sess)
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	var ret int
	var rawRecord map[interface{}]interface{}

	dec := output.NewDecoder(data, int(length))
	for {
		ret, _, rawRecord = output.GetRecord(dec)
		if ret != 0 {
			//取得完了した場合
			break
		}
		record := make(map[string]string)
		for k, v := range rawRecord {
			key := fmt.Sprintf("%v", k)
			value := fmt.Sprintf("%v", v)
			record[key] = value
		}
		j, _ := json.Marshal(record)
		fmt.Println(string(j))

		//send sqs
		sqsErr := sendMessage(string(j))
		if sqsErr != nil {
			fmt.Printf("error while sending to sqs: %v", sqsErr)
			return output.FLB_RETRY
		}
	}
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func sendMessage(record string) error {
	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(record),
		QueueUrl:     aws.String(url),
		DelaySeconds: aws.Int64(1),
	}
	_, err := svc.SendMessage(params)
	if err != nil {
		return err
	}
	fmt.Println("sent")
	return nil
}

func main() {
}
