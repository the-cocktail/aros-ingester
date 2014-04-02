package pixelwrapper

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const base64GifPixel = "R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs="

const xmlTemplate = `
<Reservation>
	<Id>%s</Id>
	<UserId>%s</UserId>
	<Total>%s</Total>
	<UserAgent>%s</UserAgent>
	<UserIP>%s</UserIP>
</Reservation>
`

func invoke_ws_wrapper(id string, body string) {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", "http://localhost:8080/reservations/"+id, strings.NewReader(body))
	req.Header.Set("Content-type", "application/xml")

	// req, err := http.NewRequest("GET", "http://www.google.es?q="+id, nil)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(resp_body))
}

func PixelHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "image/gif")
	output, _ := base64.StdEncoding.DecodeString(base64GifPixel)
	io.WriteString(res, string(output))

	user_agent := req.UserAgent()
	ip := req.RemoteAddr

	log.Println("req.URL.Path = " + req.URL.Path)
	log.Println("req.URL.RawQuery = " + req.RequestURI)

	u, err := url.Parse(req.RequestURI)
	m, _ := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic("Internal server error")
	}

	userid := m["userid"][0]
	total := m["total"][0]
	id := m["id"][0]

	log.Println("Parsed parameters")
	log.Println(" id: " + id)
	log.Println(" total: " + total)
	log.Println(" userid: " + userid)
	log.Println(" useragent: " + user_agent)
	log.Println(" ip: " + ip)

	// Y ahora hacemos una rqeuest a la propia API
	body := fmt.Sprintf(xmlTemplate, id, total, userid, user_agent, ip)

	invoke_ws_wrapper(id, body)
}
