package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"github.com/pborman/uuid"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
	"time"
)

var log = logging.MustGetLogger("rproxy")

const Sessioncookiename = string("RPROXYSESS")
const AdminFrontendPath = string("./web/dist/")

type ProxyConfig struct {
	Proxypass Host      `json:"proxypass"`
	Host      string    `json:"host"`
	Port      string    `json:"port"`
	Profiles  []Profile `json:"profiles"`
	Delay     int64     `json:"delay"`
}

type HeadersModel struct {
	Referer string
}

type Route struct {
	Path       string   `json:"path"`
	Parameters []string `json:"parameters"`
}

type Header struct {
	Id     string `json:"id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

type ClientRequest struct {
	Url 	string `json:"url"`
	Name 	string `json:"name"`
	Body 	string `json:"body"`
	Headers string `json:"headers"`
}

type HeaderList []Header
type OverrideList []*Override

func (list HeaderList) Len() int           { return len(list) }
func (list HeaderList) Less(i, j int) bool { return list[i].Key+list[i].Value < list[j].Key+list[j].Value }
func (list HeaderList) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }

func sortByHeaderKey(headers map[string]Header) HeaderList {
	rs := make(HeaderList, len(headers))
	i := 0
	for _, v := range headers {
		rs[i] = v
		i++
	}
	sort.Sort(rs)
	return rs
}

func (list OverrideList) Len() int           { return len(list) }
func (list OverrideList) Less(i, j int) bool { return list[i].Pattern < list[j].Pattern }
func (list OverrideList) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }

func sortOverridesByUrl(overrides map[string]*Override) OverrideList {
	rs := make(OverrideList, len(overrides))
	i := 0
	for _, v := range overrides {
		rs[i] = v
		i++
	}
	sort.Sort(rs)
	return rs
}

type Host struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Profile struct {
	Name    string   `json:"name"`
	Headers []Header `json:"headers"`
}

type Override struct {
	Pattern      string `json:"pattern"`
	Message      string `json:"message"`
	Responsecode int    `json:"responsecode"`
	Active       bool   `json:"active"`
	ContentType  string `json:"contenttype"`
	ResponseBody string `json:"responsebody"`
}

var headermap = make(map[string]Header, 0)
var proxyConfig ProxyConfig
var routes []Route
var overridemap = make(map[string]*Override)
var overridePatternMap = make(map[string]*regexp.Regexp)
const adminPath = "/rverse/"

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "rproxy.conf", "The json rproxy configuration file")
	//flag.StringVar(&adminPath, "adminpath", "/admin", "The path to the vuejs admin webapp")

	flag.Parse()
	proxyConfig.Delay = 0
	loadConfig()

	//if adminPath[len(adminPath)-1] != '/' {
	//	adminPath = adminPath + "/"
	//}

	if headermap == nil {
		headermap = make(map[string]Header, 0)
	}

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	// Service handler
	http.HandleFunc("/rservice/", serviceHandler())

	// Intercept handler
	for _, route := range routes {
		log.Debugf("Route path configured: %s\n", route.Path)
		http.HandleFunc(route.Path, interceptHandler(route))
	}

	// Main proxy handler
	http.HandleFunc("/", rproxy(client))

	// Admin frontend handler
	http.HandleFunc(adminPath, adminHandler())

	log.Debugf("Running proxy server on host: %s and port: %s\n", proxyConfig.Host, proxyConfig.Port)
	http.ListenAndServe(":"+proxyConfig.Port, nil)
}

func adminHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Admin handler for URI: %s\n", r.URL.RequestURI())
		renderFile(w, r)
	}
}

func renderFile(w http.ResponseWriter, request *http.Request) {

	uri := request.URL.Path[len(adminPath):len(request.URL.Path)]
	if len(uri) == 0 {
		uri = "index.html"
	}

	fileName := AdminFrontendPath + uri
	log.Debugf("rendering file: %s", fileName)

	src, err := os.Open(fileName)
	if err != nil {
		log.Errorf("Error loading file: %s\n", fileName)
		fmt.Fprintf(w, "Error loading file. %s\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// http://www.iana.org/assignments/media-types/media-types.xhtml
	ext := path.Ext(fileName)
	switch ext {
	case ".css":
		w.Header().Add("Content-Type", "text/css")
		break
	case ".html":
		w.Header().Add("Content-Type", "text/html")
		break
	case ".js":
		w.Header().Add("Content-Type", "application/javascript")
		break
	default:
		w.Header().Add("Content-Type", "text/plain")
	}

	io.Copy(w, src)
	return
}

func interceptHandler(route Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			renderTemplate(w, request, route)
			break

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Unsupported method, currently only HTTP GET's are supported.")
			break
		}
	}
}

type ConfigFile struct {
	HeaderMap   map[string]Header
	ProxyConfig ProxyConfig
	Routes      []Route
}

func loadConfig() {
	contents, err := ioutil.ReadFile("rproxy.conf")
	if err != nil {
		log.Errorf("Error loading config file: %s\n", err)
		setConfigDefaults()
	} else {
		var config ConfigFile
		err = json.Unmarshal(contents, &config)
		headermap = config.HeaderMap
		proxyConfig = config.ProxyConfig
		routes = config.Routes
	}
}

func setConfigDefaults() {
	proxyConfig = ProxyConfig{Host: "localhost", Port: "3333", Proxypass: Host{Host: "localhost", Port: "80"}}
	headermap = make(map[string]Header, 0)
	routes = make([]Route, 0)
}

func saveConfig() {
	configFile := ConfigFile{headermap, proxyConfig, routes}
	contents, _ := json.MarshalIndent(configFile, "", "  ")
	ioutil.WriteFile("rproxy.conf", contents, 0644)
}

func renderTemplate(w http.ResponseWriter, r *http.Request, route Route) {

	log.Debugf("Intercepting request for path: %s\n", r.URL.Path)

	values, err := url.ParseQuery(r.URL.Path)
	if err != nil {
		fmt.Fprint(w, "Error parsing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var referrer string
	for _, param := range route.Parameters {
		if val, ok := values[param]; ok {
			referrer = val[0]
			log.Debugf("Setting forward in template to: %s\n", val[0]) // TODO: val[0] could be array
		}
	}

	src, err := ioutil.ReadFile("route-template.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template. %s", err)
	}

	tmpl, err := template.New("test").Parse(string(src))
	if err != nil {
		fmt.Fprintf(w, "Error loading template. %s", err)
	}

	err = tmpl.Execute(w, HeadersModel{referrer})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error processing template: %s", err)
	}
}

func serviceHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		servicePath := r.URL.Path[len("/rservice/"):len(r.URL.Path)]

		//log.Debugf("request method: %s\n", r.Method)
		if r.Method == "OPTIONS" {
			w.Header().Add("Allow", "OPTIONS, GET, POST")
			return
		}

		switch servicePath {
		case "ping":
			fmt.Fprint(w, "{ pong }")
			break

		case "addrequest":

			creq = new Request{}

			break

		case "addheader":

			header := Header{}
			readBodyAndUnmarshal(w, r, &header)
			uuidValue := uuid.NewRandom().String()
			header.Id = uuidValue
			headermap[uuidValue] = header
			saveConfig()
			fmt.Fprint(w, `{ "ok": "true" }`)

			break

		case "updateHeader":

			header := Header{}
			readBodyAndUnmarshal(w, r, &header)
			headermap[header.Id] = header
			saveConfig()
			jsonHeader,_ := json.Marshal(header)
			fmt.Fprintf(w, "%s", jsonHeader)
			break

		case "toggleHeader":

			header := Header{}
			readBodyAndUnmarshal(w, r, &header)
			header.Active = !header.Active
			headermap[header.Id] = header
			saveConfig()
			break

		case "removeheader":

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("Error reading body: %v", err)
				http.Error(w, "Can't read body", http.StatusBadRequest)
				return
			}

			var header Header
			err = json.Unmarshal(body, &header)
			if err != nil {
				log.Errorf("Error parsing header: %s\n", err)
				fmt.Fprintf(w, "Error parsing Header: %s\n", err)
				return
			}

			log.Debugf("Removing header with id: %s\n", header.Id)
			delete(headermap, header.Id)
			break

		case "clearHeaders":
			headermap = make(map[string]Header)
			break

		case "listHeaders":

			headers := sortByHeaderKey(headermap)
			jsonString, _ := json.Marshal(headers)
			fmt.Fprintf(w, "%s", jsonString)
			break

		case "addRoute":

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("Error reading body: %v", err)
				http.Error(w, "Can't read body", http.StatusBadRequest)
				return
			}

			var route Route
			err = json.Unmarshal(body, &route)
			if err != nil {
				fmt.Fprintf(w, "Error parsing addRoute request: %s\n", err)
				return
			}

			routes = append(routes, route)
			saveConfig() // save config after changes
			break

		case "listRoutes":

			jsonString, _ := json.Marshal(routes)
			fmt.Fprintf(w, "%s", jsonString)
			break

		case "setHost":

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("Error reading body: %v", err)
				http.Error(w, "Can't read body", http.StatusBadRequest)
				return
			}

			type Host struct {
				Host string `json:"host"`
				Port string `json:"port"`
			}
			var host Host
			err = json.Unmarshal(body, &host)
			if err != nil {
				fmt.Fprintf(w, "Error parsing sethost request: %s\n", err)
				return
			}

			proxyConfig.Host = host.Host
			proxyConfig.Port = host.Port

			saveConfig()

			break

		case "setProxypass":

			host := Host{}
			readBodyAndUnmarshal(w, r, &host)
			fmt.Printf("Setting proxypass to: %s\n", host)
			proxyConfig.Proxypass = host
			saveConfig()

			break

		case "readProxypass":

			jsonResult, err := json.Marshal(proxyConfig.Proxypass)
			if err != nil {
				fmt.Fprintf(w, "Error reading config: %s\n", err)
			}
			fmt.Fprintf(w, "%s", jsonResult)
			break

		case "listProfiles":

			jsonString, _ := json.Marshal(proxyConfig.Profiles)
			fmt.Fprintf(w, "%s", jsonString)
			break

		case "addProfile":

			var profile Profile
			readBodyAndUnmarshal(w, r, profile)
			proxyConfig.Profiles = append(proxyConfig.Profiles, profile)
			saveConfig()
			break

		case "addOverride":

			override := Override{}
			readBodyAndUnmarshal(w, r, &override)

			overridemap[override.Pattern] = &override
			compiledPattern, err := regexp.Compile(override.Pattern)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "{ error: \"Invalid pattern: %s\" }\n", override.Pattern)
			}
			overridePatternMap[override.Pattern] = compiledPattern
			w.WriteHeader(http.StatusOK)
			break

		case "listOverrides":

			overrides := sortOverridesByUrl(overridemap)
			jsonString, _ := json.Marshal(overrides)
			fmt.Fprintf(w, "%s\n", jsonString)
			break

		case "rmOverride":

			override := Override{}
			readBodyAndUnmarshal(w, r, &override)
			delete(overridemap, override.Pattern)
			regex, err := regexp.Compile(override.Pattern)
			if err != nil {
				fmt.Printf("Unable to remove override from patternmap. Can't find compiled pattern: %s\n", regex.String())
			}
			delete(overridePatternMap, override.Pattern)
			if err != nil {
				log.Errorf("Unable to remove override from overridemap. Can't find pattern: %s\n", override.Pattern)
			}
			break

		case "toggleOverride":

			override := Override{}
			readBodyAndUnmarshal(w, r, override)
			overridemap[override.Pattern].Active = !overridemap[override.Pattern].Active

			w.WriteHeader(http.StatusNotImplemented)
			break

		case "setDelay":

			type Delay struct {
				Value int64 `json:"delay"`
			}

			var delay Delay
			readBodyAndUnmarshal(w, r, &delay)
			proxyConfig.Delay = delay.Value
			saveConfig()
			break

		case "readConfig":
			jsonString, _ := json.Marshal(proxyConfig)
			fmt.Fprintf(w, "%s", jsonString)
			break

		default:
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprint(w, "{ \"error\": \"unknown request\" }")
		}
	}
}

func readBodyAndUnmarshal(w http.ResponseWriter, r *http.Request, obj interface{}) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading body: %v", err)
		http.Error(w, "Can't read body", http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Fprintf(w, "Error parsing: %s because: %s\n", obj, err)
	}
}

func rproxy(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {

		requestUri := request.URL.RequestURI()

		for pattern, compiledPattern := range overridePatternMap {

			log.Debugf("Pattern: \"%s\" URI: \"%s\"\n", pattern, requestUri)
			if compiledPattern.MatchString(requestUri) {

				override := overridemap[pattern]

				if override.Active {
					log.Debugf("Override interception: %s with code: %s and type: %s\n", override.Pattern, override.Responsecode, override.ContentType)
                                        w.Header().Add("Content-Type", override.ContentType)
					w.WriteHeader(override.Responsecode)
					fmt.Fprintf(w, override.ResponseBody)
					return
				}
			}
		}

		// === handle rproxy session stuff ===
		//sessionCookie, err := request.Cookie(Sessioncookiename)
		_, err := request.Cookie(Sessioncookiename)
		if err != nil {
			expiration := time.Now().Add(24 * time.Hour)
			uid := uuid.NewRandom()
			log.Debugf("UUID generated: %s\n", uid.String())
			cookie := http.Cookie{Name: Sessioncookiename, Value: uid.String(), Expires: expiration, Path: "/"}
			http.SetCookie(w, &cookie)

		} else {
			sessionId := sessionCookie.Value
			log.Debugf("Existing session: %s\n", sessionId)
		}

		request.URL.Host = proxyConfig.Proxypass.Host + ":" + proxyConfig.Proxypass.Port
		proxyUrl := "http://" + request.URL.Host + request.URL.RequestURI()

		proxyRequest, err := http.NewRequest(request.Method, proxyUrl, request.Body)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range headermap {
			if v.Active {
				//log.Debugf("Augmenting with header: \"%s\" :  \"%s\"\n", v.Key, v.Value)
				proxyRequest.Header.Add(v.Key, v.Value)
			}
		}

		client.Jar, _ = cookiejar.New(nil)
		client.Jar.SetCookies(proxyRequest.URL, request.Cookies())
		for k, v := range request.Header {
			proxyRequest.Header.Add(k, strings.Join(v, ";"))
		}

		// delay if set
		if proxyConfig.Delay > 0 {
			numString := fmt.Sprintf("%dms", proxyConfig.Delay)
			d, err := time.ParseDuration(numString)
			if err != nil {
				log.Errorf("Error parsing delay: %s\n", proxyConfig.Delay)
			} else {
				time.Sleep(d)
			}
		}

		resp, err := client.Do(proxyRequest)
		if err != nil {
			fmt.Fprintf(w, "Proxy error: %s", err)
			return
		}

		// === Headers ===
		for k, v := range resp.Header {
			w.Header().Add(k, v[0])
		}

		// === Cookies ===
		for _, cookie := range resp.Cookies() {
			//log.Debugf("Cookie: %s\n", cookie.Name)
			/*
			// check for csrf token cookie and passthrough
			if strings.Compare(cookie.Name, CSRFTOKENCOOKIE) == 0 {
				log.Debugf("CSRF token cookie detected\n")
				request.Header.Add(CSRFTOKENCOOKIE, cookie.Value)
			}
			*/
			http.SetCookie(w, cookie)
		}

		//log.Debugf("Response code: %d: %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))

		//w.Header().Add("Content-Security-Policy", "style-src 'self'")
		w.WriteHeader(resp.StatusCode)


		data, err := ioutil.ReadAll(resp.Body)

		w.Write(data)
	}
}
