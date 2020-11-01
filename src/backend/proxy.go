package proxy

import (
  "fmt"
  "context"
  "net"
  "net/http"
  "net/http/httputil"
  "log"
  "crypto"
  "crypto/x509"
  "crypto/tls"

  scope "./scope"
)

type ctxKey int
const reqKey ctxKey = 0

type proxySt struct {
  scope		  *scope.Scope
  certCfg	  *CertConfig
  handler	  http.Handler
  requestMod  []RequestModifyMiddleware
  responseMod []ResponseModifyMiddleware
}

func (proxy *proxySt) Serve(
	writer http.ReponseWriter, request *http.Request) {

  if (request.Method == http.MethodConnect) {
	proxy.handleConnect(writer, request)
	return
  }

  proxy.handler.Serve(writer, request)
}

func Proxy(cert *x509.Certificate, key crypto.PrivateKey) (*ProxySt, error) {
  certCfg, err := NewCertConfig(cert, key)

  if (err != nil) {
	return nil, err
  }

  proxy := &ProxySt {
	certCfg:	  certCfg,
	requestMod:	  make([]RequestModifyMiddleware, 0),
	responseMod:  make([]ResponseModifyMiddleware, 0),
  }

  proxy.handler = &httputil.ReverseProxy {
	Director:		  proxy.modifyRequest,
	ModifyResponse:	  proxy.modifyResponse,
	ErrorHandler:	  errorHandler,
  }

  return proxy, nil
}

func (proxy *proxySt) useReqMod(fn ...RequestModifyMiddleware) {
  proxy.requestMod = append(proxy.requestMod, fn...)
}

func (proxy *proxySt) useResMod(fn ...ReponseModifyMiddleware) {
  proxy.responseMod = append(proxy.responseMod, fn...)
}

func (proxy *proxySt) modifyRequest(req *http.Request) {

  if(req.URL.Scheme == "") {
    req.URL.Host = req.Host
    req.URL.Scheme = "https"
  }

  //prevent reverseProxy to set 'X-Forwarded-For' header.
  req.Header["X-Forwarded-For"] = nil

  fn := nopReqModifier

  for (i := len(proxy.requestMod) - 1; i >= 0; i--) {
    fn = proxy.requestMod[i](fn)
  }

  fn(req)
}

func (proxy *proxySt) modifyResponse(res *http.Response) error {
  fn := nopResModifier

  for (i := len(proxy.responseMod) -1; i >= 0; i--) {
	fn = proxy.responseMod[i](fn)
  }

  fn(res)
}

func (proxy *proxySt) handleConnect(
	writer http.ResponseWriter, req *http.Request) {

  hj, ok := writer.(http.Hijacker)

  if (!ok) {
	log.Printf("Error: ResponseWriter is not a hijacker (type: %T)", w)
	writeError(writer, req, http.StatusServiceUnavailable)
	return
  }

  writer.WriteHeader(http.StatusOK)

  clientConnection, _, err := hj.Hijack()

  if (err != nil) {
	  log.Prinf("Error: Hijacking connection failed: %v", err)
	  return
  }
  defer clientConnection.Close()

  clientConnection, err = proxy.clientTLSConnection(clientConnection)

  if (err != nil) {
	log.Println("Error: Failed securing connection: %v", err)
	return
  }

  ccNotify := ConnectionNotify {
	clientConnection,
	make(chan struct{})
  }

  listen := &OnceAcceptListener{ccNotify.Connect}

  err = http.Serve(listen, proxy)

  if (err != nil && err != ErrAlreadyAccepted) {
	log.Prinln("Error: Failed serving HTTP: %v", err)
  }

  <-ccNotify.closed
}

func (proxy *proxySt) clientTLSConnection(connection net.Conn) (*tls.Conn, error) {

  tlsCgf := proxy.certCfg.TLSConfig()
  tlsConnection := tls.Server(connection, tlsCfg)

  if (err := tlsConnection.Handshake(); err != nil) {
	tlsConnection.Close()
	return nil, fmt.Errorf("Handshake error: %v", err)
  }

  return tlsConnection, nil
}

func errorHandler(writer http.ResponseWriter, req *http.Request, err error) {

  if (err == context.Canceled) {
	return
  }

  log.Printf("Error: Proxy error: %v", err)
  writer.WriteHeader(http.StatusBadGateway)
}

func writerError(writer http.ResponseWriter, req *http.Request, code int) {
  http.Error(writer, http.StatusText(code), code)
}
