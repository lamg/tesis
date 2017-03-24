package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/lamg/tesis"
	"io"
	"io/ioutil"
	"log"
	h "net/http"
)

type PortalUser struct {
	client *h.Client
	ck     *h.Cookie
	auInf  string
	index  string
}

func NewPortalUser(url string) (p *PortalUser) {
	cfg := &tls.Config{InsecureSkipVerify: true}
	tr := &h.Transport{TLSClientConfig: cfg}
	cl := &h.Client{Transport: tr}
	ai := fmt.Sprintf("https://%s%s", url, dashP)
	in := fmt.Sprintf("https://%s", url)
	p = &PortalUser{client: cl, auInf: ai, index: in}
	return
}

func (p *PortalUser) Auth(user, pass string) (a bool, e error) {
	a = false
	var r *h.Response
	var cr *tesis.Credentials
	var bs []byte
	cr = &tesis.Credentials{User: user, Pass: pass}
	bs, e = json.Marshal(cr)
	if e == nil {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		r, e = p.client.Post(p.index+authP, "application/json", rd)
	} else {
		log.Panicf("Incorrect program: %s", e.Error())
	}
	if e == nil {
		if r.StatusCode == 200 {
			var cs []*h.Cookie
			cs = r.Cookies()
			if len(cs) == 1 {
				p.ck = cs[0]
				//token string stored
				a = true
			} else {
				bs, e = ioutil.ReadAll(r.Body)
				if e == nil {
					e = fmt.Errorf("Cantidad de cookies %d ≠ 1, %s", len(cs), string(bs))
				}
			}
		} else {
			e = fmt.Errorf("Status = %s", r.Status)
		}
	}
	return
}

func (p *PortalUser) Info() (s string, e error) {
	var q *h.Request
	if p.ck == nil {
		e = fmt.Errorf("Auth failed")
	}
	if e == nil {
		q, e = h.NewRequest("GET", p.index+dashP, nil)
	}
	if e == nil {
		q.AddCookie(p.ck)
		var rp *h.Response
		//create a request with appropiate header
		rp, e = p.client.Do(q)
		var b []byte
		if e == nil && rp.StatusCode == 200 {
			b, e = ioutil.ReadAll(rp.Body)
		} else if e == nil {
			e = fmt.Errorf("Status = %s", rp.Status)
		}
		if e == nil {
			s = string(b)
			e = rp.Body.Close()
		}
	}
	return
}

func (p *PortalUser) Index() (s string, e error) {
	var r *h.Response
	var b []byte

	r, e = p.client.Get(p.index)
	// { responseGet.(p.index).r ≡ e = nil }
	if e == nil {
		b, e = ioutil.ReadAll(r.Body)
	}
	// { read.(r.Body).b ≡ e = nil }
	if e == nil {
		s = string(b)
	}
	// { s = string.respGetIndex ≡ e = nil }
	return
}

func (p *PortalUser) Sync() (s string, e error) {
	//marshall Info to JSON and send to server
	//to use syncH handler
	var r *h.Response
	var b []byte
	var acs []tesis.AccMatch
	var rd io.Reader
	acs = []tesis.AccMatch{
		tesis.AccMatch{
			DBId:    "1",
			ADId:    "3",
			ADName:  "LUIS",
			SrcName: "Luis",
			SrcDB:   "ASET",
		},
	}
	b, e = json.Marshal(acs)
	rd = bytes.NewReader(b)
	var q *h.Request
	if p.ck == nil {
		e = fmt.Errorf("Auth failed")
	}
	if e == nil {
		q, e = h.NewRequest("POST", p.index+syncP, rd)
	}
	if e == nil {
		q.AddCookie(p.ck)
		q.Header.Set("content-type", "application/json")
		r, e = p.client.Do(q)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
	}
	if e == nil {
		s = string(bs)
		e = r.Body.Close()
	}
	// { syncRespStr.s ≡ e = nil }
	return
}
