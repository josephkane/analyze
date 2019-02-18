package app

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dre1080/recover"
	"github.com/go-openapi/loads"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"

	"github.com/supergiant/analyze/asset"
	"github.com/supergiant/analyze/pkg/analyze"
	"github.com/supergiant/analyze/pkg/api"
	"github.com/supergiant/analyze/pkg/api/handlers"
	"github.com/supergiant/analyze/pkg/api/operations"
	"github.com/supergiant/analyze/pkg/config"
	"github.com/supergiant/analyze/pkg/kube"
	"github.com/supergiant/analyze/pkg/logger"
	"github.com/supergiant/analyze/pkg/models"
	"github.com/supergiant/analyze/pkg/scheduler"
	"github.com/supergiant/analyze/pkg/storage"
	"github.com/supergiant/analyze/pkg/storage/etcd"
)

func RunCommand(cmd *cobra.Command, _ []string) error {
	configFilePaths, err := cmd.Flags().GetStringArray("config")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag value")
	}

	cfg := &analyze.Config{}

	// configFileReadError is not critical due to possibility that configuration is done by environment variables
	configFileReadError := config.ReadFromFiles(cfg, configFilePaths)

	if err = config.MergeEnv("AZ", cfg); err != nil {
		return errors.Wrap(err, "unable to merge env variables")
	}

	//TODO: try to unify APIs discovery which are hosted in k8s
	//TODO: and rewrite config population logic
	if etcdEndpoint := discoverETCDEndpoint(); etcdEndpoint != "" {
		cfg.ETCD.Endpoints = append(cfg.ETCD.Endpoints, discoverETCDEndpoint())
	}

	log := logger.NewLogger(cfg.Logging).WithField("app", "analyze-core")
	mainLogger := log.WithField("component", "main")

	mainLogger.Infof("config: %+v", cfg)
	mainLogger.Infof("config file name: %s", config.UsedFileName())
	if configFileReadError != nil {
		log.Warnf("unable to read config file, %v", configFileReadError)
	}

	if err := cfg.Validate(); err != nil {
		return errors.Wrap(err, "config validation error")
	}

	kubeClient, err := kube.NewKubeClient(log.WithField("component", "kubeClient"))
	if err != nil {
		return errors.Wrap(err, "unable to create kube client")
	}

	etcdStorage, err := etcd.NewETCDStorage(cfg.ETCD, log.WithField("component", "etcdClient"))
	if err != nil {
		return errors.Wrap(err, "unable to create ETCD client")
	}

	defer etcdStorage.Close()

	scheduler := scheduler.NewScheduler(log.WithField("component", "scheduler"))
	defer scheduler.Stop()

	watchChan := etcdStorage.WatchRange(context.Background(), models.PluginPrefix)
	log.Debug("watch stated")
	pluginController := analyze.NewPluginController(
		watchChan,
		etcdStorage,
		kubeClient,
		scheduler,
		log.WithField("component", "pluginController"),
	)

	go pluginController.Loop()

	swaggerSpec, err := loads.Analyzed(api.SwaggerJSON, "2.0")
	if err != nil {
		return errors.Wrap(err, "unable to create spec analyzed document")
	}

	//TODO: add request logging middleware
	//TODO: add metrics middleware
	analyzeAPI := operations.NewAnalyzeAPI(swaggerSpec)
	analyzeAPI.Logger = log.WithField("component", "analyzeApi").Errorf

	analyzeAPI.GetCheckResultsHandler = handlers.NewChecksResultsHandler(
		etcdStorage,
		log.WithField("handler", "CheckResultsHandler"),
	)
	analyzeAPI.GetPluginHandler = handlers.NewPluginHandler(
		etcdStorage,
		log.WithField("handler", "PluginHandler"),
	)
	analyzeAPI.GetPluginsHandler = handlers.NewPluginsHandler(
		etcdStorage,
		log.WithField("handler", "PluginsHandler"),
	)
	analyzeAPI.RegisterPluginHandler = handlers.NewRegisterPluginHandler(
		etcdStorage,
		log.WithField("handler", "RegisterPluginHandler"),
	)
	analyzeAPI.UnregisterPluginHandler = handlers.NewUnregisterPluginHandler(
		etcdStorage,
		log.WithField("handler", "UnregisterPluginHandler"),
	)

	err = analyzeAPI.Validate()
	if err != nil {
		return errors.Wrap(err, "API configuration error")
	}

	server := api.NewServer(analyzeAPI)
	server.Port = cfg.API.ServerPort
	server.Host = cfg.API.ServerHost
	server.ConfigureAPI()

	handlerWithRecovery := recover.New(&recover.Options{
		Log: logrus.Error,
	})

	//TODO fix CORS till release
	corsHandler := cors.New(cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{},
		MaxAge:         1000,
	}).Handler

	handler := alice.New(
		handlerWithRecovery,
		corsHandler,
		swaggerMiddleware,
		newProxyMiddleware(etcdStorage, log.WithField("middleware", "proxy")),
		uiMiddleware,
	).Then(analyzeAPI.Serve(nil))

	server.SetHandler(handler)

	defer server.Shutdown()

	if servingError := server.Serve(); servingError != nil {
		return errors.Wrap(servingError, "unable to serve HTTP API")
	}

	return nil
}

func discoverETCDEndpoint() string {
	etcdHost, hostExists := os.LookupEnv("ETCD_SERVICE_HOST")
	etcdPort, portExists := os.LookupEnv("ETCD_SERVICE_PORT")
	if !hostExists || !portExists {
		return ""
	}
	return etcdHost + ":" + etcdPort
}

func swaggerMiddleware(handler http.Handler) http.Handler {
	var staticServer = http.FileServer(asset.Assets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/api/v1/swagger-ui" || r.URL.Path == "/api/v1/help" {
			http.Redirect(w, r, "/api/v1/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.HasPrefix(r.URL.Path, "/api/v1/swagger-ui/") {
			url := strings.TrimPrefix(r.URL.Path, "/api/v1/swagger-ui/")
			r.URL.Path = "/swagger/" + url
			staticServer.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func uiMiddleware(handler http.Handler) http.Handler {
	var staticServer = http.FileServer(asset.Assets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.URL.Path, "/api/v1") {
			r.URL.Path = "/ui" + r.URL.Path
			staticServer.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func newProxyMiddleware(storage storage.Interface, logger logrus.FieldLogger) func(handler http.Handler) http.Handler {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic("can't get kube config")
	}

	tr, err := rest.TransportFor(config)
	if err != nil {
		panic("can't get transport")
	}

	var proxies = make(map[string]*httputil.ReverseProxy)

	return func(handler http.Handler) http.Handler {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		pluginsRaw, err := storage.GetAll(ctx, models.PluginPrefix)
		if err != nil {
			panic("can't get plugins")
		}

		for _, rawPlugin := range pluginsRaw {
			p := &models.Plugin{}
			err := p.UnmarshalBinary(rawPlugin.Payload())
			if err != nil {
				panic("can't unmarshal plugin")
			}
			_, exists := proxies[p.ID]
			if !exists {
				url, err := url.Parse("http://" + p.ServiceEndpoint)
				if err != nil {
					panic("can't parse host")
				}

				logger.Debugf("create proxy for url: %+v", *url)
				reverseProxy := httputil.NewSingleHostReverseProxy(url)
				reverseProxy.Transport = tr
				reverseProxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
					logger.Errorf("reverse proxy error params: %+v", *req.URL)
					logger.Errorf("reverse proxy error: %v", err)
					rw.WriteHeader(http.StatusBadGateway)
				}
				proxies[p.ID] = reverseProxy
			}
		}

		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			var targetProxy *httputil.ReverseProxy
			for id, proxy := range proxies {
				if strings.Contains(req.URL.Path, id) {
					targetProxy = proxy
				}
			}
			if targetProxy == nil {
				handler.ServeHTTP(res, req)
				return
			}

			logger.Debugf("got proxy request at: %v, request: %+v", time.Now(), req.URL)
			defer logger.Debugf("proxy request processing finished at: %v, request: %+v", time.Now(), req.URL)

			// Update the headers to allow for SSL redirection
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

			// Note that ServeHttp is non blocking and uses a go routine under the hood
			targetProxy.ServeHTTP(res, req)
		})
	}
}
