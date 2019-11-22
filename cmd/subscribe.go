package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	devutil "github.com/projectriff/developer-utils/pkg"
	client "github.com/projectriff/stream-client-go"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var payloadAsString bool
var offset int32
var namespace string
var streamGVR = schema.GroupVersionResource{
	Group:    "streaming.projectriff.io",
	Version:  "v1alpha1",
	Resource: "streams",
}

type Event struct {
	Payload     string `json:"payload"`
	ContentType string `json:"contentType"`
	Headers     map[string]string `json:"headers"`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var eventHandler = func(ctx context.Context, payload io.Reader, contentType string, headers map[string]string) error {
	bytes, err := ioutil.ReadAll(payload)
	if err != nil {
		return err
	}
	var payloadStr string
	if payloadAsString {
		payloadStr = string(bytes)
	} else {
		payloadStr = base64.StdEncoding.EncodeToString(bytes)
	}
	evt := Event{
		Payload:     payloadStr,
		ContentType: contentType,
		Headers:     headers,
	}
	marshaledEvt, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", marshaledEvt)
	return nil
}

var rootCmd = &cobra.Command{
	Use:     "subscribe-stream <stream-name>",
	Short:   "subscribe for events from the given stream",
	Long:    "",
	Example: "subscribe-stream letters --payload-as-string",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args", args)
		k8sClient := devutil.NewK8sClient()
		topic, err := k8sClient.GetNestedString(args[0], streamGVR, "status", "address", "topic")
		if err != nil {
			fmt.Println("error while determining topic name for stream", err)
		}
		fmt.Println("topic:", topic)
		gateway, err := k8sClient.GetNestedString(args[0], streamGVR, "status", "address", "gateway")
		if err != nil {
			fmt.Println("error while determining gateway address for stream", err)
		}
		fmt.Println("gateway:", gateway)
		contentType, err := k8sClient.GetNestedString(args[0], streamGVR, "spec", "contentType")
		if err != nil {
			fmt.Println("error while determining contentType for stream", err)
		}
		fmt.Println("contentType", contentType)

		sc, err := client.NewStreamClient(gateway, topic, contentType)
		if err != nil {
			panic(err)
		}

		var eventErrHandler client.EventErrHandler
		eventErrHandler = func(cancel context.CancelFunc, err error) {
			fmt.Printf("ERROR: %v\n", err)
		}
		_, err = sc.Subscribe(context.Background(), "monitor", 0, eventHandler, eventErrHandler)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&payloadAsString, "payload-as-string", "", false,
		"display the payload as string rather than base64 encoded string")
	rootCmd.Flags().Int32P("offset", "", offset, "subscribe to events after the given offset")
	rootCmd.Flags().StringP("namespace", "n", namespace, "namespace of the stream")
}
