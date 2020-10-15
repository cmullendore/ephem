package webui

import (
	"crypto/tls"
	"encoding/base64"
	"flag"
	"html/template"
	"math/rand"
	"net/http"
	"strings"

	"log"

	"github.com/cmullendore/ephem/src/config"
	"github.com/cmullendore/ephem/src/core"
)

// UIServer is the primary engine object, containing the appropriate
// configuration, cache, and log objects necessary to execute any
// given intbound request.
type UIServer struct {
	SecretsEngine core.ISecretsEngine
	Configuration *config.WebUI
	Log           *log.Logger
	fs            http.FileSystem
}

// NewUIServer creates a new request handler and initializes the
// internal objects, returning a complete and useable object.
func NewUIServer(se core.ISecretsEngine) *UIServer {
	var s UIServer
	s.SecretsEngine = se
	s.Configuration = config.LoadWebUIConfig()
	return &s
}

// Listen starts the HTTP Listener, configuring the listener based
// on the handler's configuration and using the default
// request handler.
func (s *UIServer) Listen() {

	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: *insecure,
		MinVersion:         tls.VersionTLS12,
	}

	http.HandleFunc("/get/", s.getItem)
	http.HandleFunc("/save/", s.saveItem)

	// Serves content from content.go bindata file created with
	// go-bindata.exe -fs -prefix "src/webui/content/" -o src/webui/content.go src/webui/content/...
	http.Handle("/", http.FileServer(AssetFile()))

	server := http.Server{
		Addr:      s.Configuration.Listener,
		TLSConfig: tlsConfig,
	}

	server.ListenAndServeTLS(s.Configuration.ListenerCert.TLSCertPath, s.Configuration.ListenerCert.TLSKeyPath)
}

func (s *UIServer) getItem(w http.ResponseWriter, r *http.Request) {
	pathKey := strings.TrimPrefix(r.URL.Path, "/get/")
	item, err := s.SecretsEngine.GetItem(&pathKey)

	var te = "404.html"
	var obj string
	if err != nil {
		log.Println(err)
	} else if item == nil {
		obj = ""
	} else {
		obj = string(*item)
		te = "get.html"
	}

	t := binTemplate(te)

	t.Execute(w, obj)
}

func (s *UIServer) saveItem(w http.ResponseWriter, r *http.Request) {

	pathBytes := make([]byte, s.Configuration.PathLength)

	for i := 0; i < s.Configuration.PathLength; i++ {

		r := rand.Intn(256)
		pathBytes[i] = byte(r)
	}

	path64 := base64.RawURLEncoding.EncodeToString(pathBytes)
	path := "/get/" + path64

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Println(err)
	}

	if len(r.MultipartForm.Value["saveContent"]) > 0&len(strings.TrimSpace(r.MultipartForm.Value["saveContent"][0])) {
		data := []byte(r.MultipartForm.Value["saveContent"][0])

		eng := s.SecretsEngine
		err := eng.SaveItem(&path64, &data)
		if err != nil {
			log.Println(err)
		}
	}

	var url = r.URL.Scheme + "https://" + r.Host + path
	t := binTemplate("save.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, url)
}

func binTemplate(name string) *template.Template {
	var t *template.Template

	b, _ := Asset(name)
	t, _ = template.New("").Parse(string(b))

	return t
}
