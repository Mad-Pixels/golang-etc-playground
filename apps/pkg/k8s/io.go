package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"

	v1 "k8s.io/api/core/v1"
)

func Read(name, ns string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	req := client.CoreV1().Pods(ns).GetLogs(name, &v1.PodLogOptions{})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Logs for %s in %s:\n%s\n", name, ns, buf.String())
	return nil
}
