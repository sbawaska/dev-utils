# Developer Tools
This repository provides tools for riff users to develop and debug functions. These tools are bundled in a container image that is meant to be run in the development cluster.

## Using the tools
These tools can be used by running a simple pod in the k8s cluster with a configuration like this:
```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: dev-utils
spec:
  containers:
  - name: dev-utils
    image: projectriff/dev-utils
EOF
```
Then depending upon which riff runtime you are using, create the appropriate clusterrolebindings:
```bash
kubectl create clusterrolebinding dev-util-stream --clusterrole=riff-streaming-readonly-role --serviceaccount=default:default
kubectl create clusterrolebinding dev-util-core --clusterrole=riff-core-readonly-role --serviceaccount=default:default
kubectl create clusterrolebinding dev-util-knative --clusterrole=riff-knative-readonly-role --serviceaccount=default:default
```

## Included tools
1. **invoke-core:** To invoke the given core deployer.  
The command takes the form:  
    ```
    invoke-core <deployer-name> -n <namespace> -- <curl-params>
    ```
    where everything after `--` is passed as parameter to curl
1. **invoke-knative:** To invoke the given knative deployer.
The command takes the form:  
    ```
    invoke-knative <deployer-name> -n <namespace> -- <curl-params>
    ```
    where everything after `--` is passed as parameter to curl

1. **publish:** To publish an event to the given stream.
The command takes the form:
    ```
    publish <stream-name> -n <namespace> --payload <payload-as-string> --content-type <content-type> --header "<header-name>: <header-value>"
    ```
    where `stream-name`, `payload` and `content-type` are mandatory and `header` can be used multiple times.
1. **subscribe:** To subscribe for events from the given stream.
The command takes the form:
    ```
    subscribe <stream-name> --payload-as-string --from-beginning
    ```
    If the `--from-beginning` option is present, display all the events in the stream, otherwise only new events are displayed in the following json format:
    ```
    {"payload": "base64 encoded user payload","content-type": "the content type of the message","headers": {"user provided header": "while publishing"}}
    ```
    The payload will be base64 encoded unless the `--payload-as-string` flag is present, in which case it will be displayed as a string.

The namespace parameter is optional for all the commands. If not specified, the namespace of the `dev-utils` pod will be assumed.

## Examples
These tools can be invoked using kubectl exec. some examples follow:
```bash
kubectl exec dev-utils -- invoke-core upper -- -H "Content-Type:text/plain" -H "Accept:text/plain" -d test
kubectl exec dev-utils -- publish letters --payload foo
kubectl exec dev-utils -- subscribe letters --payload-as-string
```