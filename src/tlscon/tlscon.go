package tlscon

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
)

func OpenTLSClient(ipPort string) (*tls.Conn, error) {

	// Note this loads standard x509 certificates, test keys can be
	// generated with makecert.sh
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	// InsecureSkipVerify required for unsigned certs with Go1 and later.
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ipPort, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	// This shows the public key of the server, we will accept any key, but
	// we could terminate the connection based on the public key if desired.
	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println("Client: Server public key is:")
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
	}
	// Lets verify behind the doubt that both ends of the connection
	// have completed the handshake and negotiated a SSL connection
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	// All TLS handling has completed, now to pass the connection off to
	// go-rpcgen/protobuf/AddService
	return conn, err
}

func handleClient(conn net.Conn) {
	tlscon, ok := conn.(*tls.Conn)
	if ok {
		log.Print("server: conn: type assert to TLS succeedded")
		err := tlscon.Handshake()
		if err != nil {
			log.Fatalf("server: handshake failed: %s", err)
		} else {
			log.Print("server: conn: Handshake completed")
		}
		state := tlscon.ConnectionState()
		// Note we could reject clients if we don't like their public key.      
		log.Println("Server: client public key is:")
		for _, v := range state.PeerCertificates {
			log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
		}
		// Now that we have completed SSL/TLS 
		addservice.ServeAddService(tlscon, Add{})
	}
}

func ServerTLSListen(service string) {

	// Load x509 certificates for our private/public key, makecert.sh will
	// generate them for you.

	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	// Note if we don't tls.RequireAnyClientCert client side certs are ignored.
	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	// Keep this loop simple/fast as to be able to handle new connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		// Fire off go routing to handle rest of connection.
		go handleClient(conn)
	}
}