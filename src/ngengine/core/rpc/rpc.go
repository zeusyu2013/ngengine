package rpc

import (
	"fmt"
	"net"
	"ngengine/logger"
	"ngengine/share"
	"runtime"
)

func GetServiceMethod(m string) string {
	return fmt.Sprintf("S2S%s", m)
}

func GetHandleMethod(m string) string {
	return fmt.Sprintf("C2S%s", m)
}

func CreateRpcService(service map[string]interface{}, handle map[string]interface{}, ch chan *RpcCall, log *logger.Log) (rpcsvr *Server, err error) {
	rpcsvr = NewServer(ch, log)
	for k, v := range service {
		err = rpcsvr.RegisterName(GetServiceMethod(k), v)
		if err != nil {
			return
		}
	}

	for k, v := range handle {
		err = rpcsvr.RegisterName(GetHandleMethod(k), v)
		if err != nil {
			return
		}
	}

	return
}

func CreateService(rs *Server, l net.Listener, log *logger.Log) {
	log.LogInfo("rpc start at:", l.Addr().String())
	for {
		conn, err := l.Accept()
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			log.LogWarn("TCP", "temporary Accept() failure - ", err.Error())
			runtime.Gosched()
			continue
		}
		if err != nil {
			log.LogWarn("rpc service quit")
			break
		}
		//启动服务
		log.LogInfo("new rpc client,", conn.RemoteAddr())
		go rs.ServeConn(conn, share.MAX_BUF_LEN)
	}
}
