package main

import (
	"io"
	"log"
	"net"
	"sync"
        "strconv"
	"strings"
	"fmt"
	"github.com/mafl97/is105sem03/mycrypt"
	"github.com/mafl97/funtemps/conv"
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:16")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}

					krypterMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {
  				        case "ping":
						_, err = c.Write([]byte("pong"))
					  case "Kjevik":
						fields := strings.Split(msg, ";")
						lastField := fields[len(fields)-1]
						celsius, err := strconv.ParseFloat(lastField, 64)
						if err != nil {
							log.Println(err)
							continue
						}
						fahrenheit := conv.CelsiusToFarhenheit(celsius)
						if err != nil {
							log.Println(err)
							continue
						}
						msg := fmt.Sprintf("%.2f", fahrenheit)
						_, err = c.Write([]byte(msg))
					default:
						_, err = c.Write([]byte(string(krypterMelding)))
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for l  kke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
