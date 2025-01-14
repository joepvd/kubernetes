package filters

import (
	"net/http"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/kubernetes/openshift-kube-apiserver/filters/deprecatedapirequest"
)

// WithDeprecatedApiRequestLogging adds an http handler that logs requests to deprecated apis.
func WithDeprecatedApiRequestLogging(handler http.Handler, controller deprecatedapirequest.APIRequestLogger) http.Handler {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer handler.ServeHTTP(w, req)
		info, ok := request.RequestInfoFrom(req.Context())
		if !ok {
			return
		}
		if !controller.IsDeprecated(info.Resource, info.APIVersion, info.APIGroup) {
			return
		}
		timestamp, ok := request.ReceivedTimestampFrom(req.Context())
		if !ok {
			return
		}
		user, ok := request.UserFrom(req.Context())
		if !ok {
			return
		}
		controller.LogRequest(
			schema.GroupVersionResource{
				Group:    info.APIGroup,
				Version:  info.APIVersion,
				Resource: info.Resource,
			},
			timestamp,
			user.GetName(),
			info.Verb,
		)
	})
	return handlerFunc
}
