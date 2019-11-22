package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	devutil "github.com/projectriff/developer-utils/pkg"
	client "github.com/projectriff/stream-client-go"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var namespace string
var payload string
var contentType string
var header []string
var streamGVR = schema.GroupVersionResource{
	Group:    "streaming.projectriff.io",
	Version:  "v1alpha1",
	Resource: "streams",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "publish <stream-name> <payload>",
	Short:   "publish events to the given stream",
	Long:    "",
	Example: "publish letters --payload=my-value",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		k8sClient := devutil.NewK8sClient()
		topic, err := k8sClient.GetNestedString(args[0], streamGVR, "status", "address", "topic")
		if err != nil {
			fmt.Println("error while determining topic name for stream", err)
			os.Exit(1)
		}

		gateway, err := k8sClient.GetNestedString(args[0], streamGVR, "status", "address", "gateway")
		if err != nil {
			fmt.Println("error while determining gateway address for stream", err)
			os.Exit(1)
		}

		contentType, err := k8sClient.GetNestedString(args[0], streamGVR, "spec", "contentType")
		if err != nil {
			fmt.Println("error while determining contentType for stream", err)
			os.Exit(1)
		}

		m, err := getMapFromHeaders(header)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sc, err := client.NewStreamClient(gateway, topic, contentType)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = sc.Publish(context.Background(), strings.NewReader(payload), nil, contentType, m)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace of the stream")
	rootCmd.Flags().StringVarP(&payload, "payload", "p", "", "the content/payload to publish to stream")
	rootCmd.Flags().StringVarP(&contentType, "content-type", "c", "", "mime type of content")
	rootCmd.Flags().StringArrayVarP(&header, "header", "", header, "headers for the payload")
}

func getMapFromHeaders(headers []string) (map[string]string, error) {
	returnVal := map[string]string{}
	for _, h := range headers {
		splitH := strings.Split(h, ":")
		if len(splitH) != 2 {
			return nil, errors.New(fmt.Sprintf("illegal header: %s, expected form: k1:v1", h))
		}
		returnVal[splitH[0]] = splitH[1]
	}
	return returnVal, nil
}