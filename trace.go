package katana

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	dlTagUndefined        = "_undef"
	dlTagMysqlFailed      = "_com_mysql_failure"
	dlTagMysqlSuccess     = "_com_mysql_success"
	dlTagRedisFailed      = "_com_redis_failure"
	dlTagRedisSuccess     = "_com_redis_success"
	dlTagThriftFailed     = "_com_thrift_failure"
	dlTagThriftSuccess    = "_com_thrift_success"
	dlTagHTTPSuccess      = "_com_http_success"
	dlTagHTTPFailed       = "_com_http_failure"
	dlTagBackendRPCFailed = "_com_interactive_failure"
	dlTagRequestIn        = "_com_request_in"
	dlTagRequestOut       = "_com_request_out"

	ktnHeaderTraceID = "ktn-header-trace-id"
	ktnHeaderSpanID  = "ktn-header-span-id"
	ktnHeaderCspanID = "ktn-header-cspan-id" // caller 带过来的 cspanid, 即为本服务的 spanid, 调用下游时, 生成下游使用的 cspanid
)

// Trace ...
type Trace struct {
	TraceID string
	SpanID  string
	CspanID string
}

// GetTrace ...
func GetTrace(req *http.Request) *Trace {
	t := &Trace{
		TraceID: req.Header.Get(ktnHeaderTraceID),
		SpanID:  req.Header.Get(ktnHeaderCspanID), // caller 带过来的 cspanid, 即为本服务的 spanid, 调用下游时, 生成下游使用的 cspanid
		CspanID: "",
	}

	if t.TraceID == "" {
		t.TraceID = genTraceID()
		req.Header.Set(ktnHeaderTraceID, t.TraceID) // 防止多次生成
	}

	if t.SpanID == "" {
		t.SpanID = genSpanID()
		req.Header.Set(ktnHeaderSpanID, t.SpanID) // 防止多次生成
	}

	return t
}

// ToMap ...
func (t *Trace) ToMap() map[string]string {
	return map[string]string{
		ktnHeaderTraceID: t.TraceID,
		ktnHeaderSpanID:  t.SpanID,
	}
}

func genTraceID() string {
	ip := getLocalIP()
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()
	b := bytes.Buffer{}

	b.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))
	b.WriteString(fmt.Sprintf("%x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0")

	return b.String()
}

func genSpanID() string {
	return fmt.Sprintf("%x", rand.Int63())
}

func getLocalIP() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}
