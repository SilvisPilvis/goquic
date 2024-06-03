package main

import (
	"context"
	"goquic/ECS"
	"strconv"
	"strings"

	// "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"math/rand/v2"

	// "math/rand"

	// "encoding/binary"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	// "unsafe"

	"github.com/deeean/go-vector/vector2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/quic-go/quic-go"
	// "goquic/model"
	// "goquic/vectors"
)

// var channel = make(chan string)

// type Player struct {
// 	ECS.Player
// 	id   int64
// 	name string
// }

type Player struct {
	ECS.Player
	Id   int64
	Name string
}

func (p *Player) Update(deltaTime float64) {
	// Update position using velocity, speed, and deltaTime
	p.Position.X += p.Velocity.X * p.Speed * deltaTime
	p.Position.Y += p.Velocity.Y * p.Speed * deltaTime
	// p.Position.X += (p.Velocity.X + p.Speed) * deltaTime
	// p.Position.Y += (p.Velocity.Y + p.Speed) * deltaTime
}

type Users struct {
	players map[int64]*Player
}

func MovePlayer(u *Users, id int64, delta *vector2.Vector2) {
	multi := u.players[id].Player.Transform.Velocity.Mul(delta)
	// multi := u.players[id].velocity.Mul(delta)
	u.players[id].Player.Transform.Position = *u.players[id].Player.Transform.Position.Add(multi)
}

func InitUsers() *Users {
	return &Users{
		players: make(map[int64]*Player),
	}
}

func (u *Users) AddPlayer(id int64, player *Player) {
	u.players[player.Id] = player
}

// func RemovePlayer(u *Users, id int64) {
// 	delete(u.players, id)
// }

func (u *Users) RemovePlayer(id int64) {
	delete(u.players, id)
}

const (
	JOIN = iota
	LEAVE
	MOVE
)

// type Message struct {
// 	// Action string
// 	Action int
// 	// Data   string
// 	Data map[string]interface{}
// }

var u = InitUsers()

func HandleMessage(conn quic.Connection, m *Message, u *Users, s quic.Stream) {
	log.Println("Data from client: ", m.Action)
	// handle messages
	switch m.Action {
	case JOIN:
		// add player
		// set id as the number of current player
		log.Print("[JOIN] Data from client: ", m.Action)
		u.AddPlayer(int64(len(u.players)), &Player{
			Id: int64(len(u.players)),
			// name: "Test",
			// random name generation for testing
			Name: RandString(6),
		})
		// u.AddPlayer(int64(rand.IntN(6)), &Player{})

		log.Println("Number of players:", len(u.players))
		log.Println("Players:", u.players[int64(len(u.players)-1)])
		// send user id back to client
	case LEAVE:
		// remove player
		u.RemovePlayer(m.Data["id"].(int64))
		// RemovePlayer(u, m.Data["id"].(int64))
		// sync changes with client
	case MOVE:
		// move player
		MovePlayer(u, m.Data["id"].(int64), m.Data["delta"].(*vector2.Vector2))
	default:
		log.Println("Unknown action:", m.Action)
	}
	// log.Println("Message: ", m.Action)
}

func HandleConnection(conn quic.Connection, ctx context.Context, dataChan chan string) {
	defer conn.CloseWithError(0, "closing connection")

	if dataChan == nil {
		log.Println("Channel not created")
		return
	}

	s, err := conn.AcceptStream(ctx)
	if err != nil {
		log.Println("Error accepting stream:", err)
		dataChan <- "Error accepting stream: " + err.Error()
		return
	}
	defer s.Close()

	// decoder := json.NewDecoder(s)
	// u := InitUsers()
	// creates a buffer of 1024 bytes to store recieved data
	buf := make([]byte, 1024)
	for {
		// no data in chanel because infinite loop doesn't return

		// Handle incoming data from the client
		var m Message
		// use this if you want to decode data
		// err := decoder.Decode(&m)
		n, err := s.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from stream: ", err)
				// the bellow code prints if the connection timed out
				log.Println("Action: ", buf[:n])
				// log.Println("Action: ", buf)
				dataChan <- "Error reading from stream: " + err.Error()
				// os.Exit(0)
				break
			}
		}

		// m.Action = 0
		var res strings.Builder
		for i := range len(buf) {
			// conv, err := strconv.Atoi(string(buf[i]))
			conv := string(buf[i])
			if err != nil {
				log.Fatal("Failed to decode stream data.")
			}
			res.WriteString(conv)
		}
		data := strings.Split(res.String(), "|")
		m.Action, err = strconv.Atoi(data[0])
		if err != nil {
			log.Println("Failed to convert Action: ", err)
		}
		log.Println(data)
		// m.Data = res.String()

		// decodes the data from the client
		// breaks if index is not 0
		// m.Action, err = strconv.Atoi(string(buf[0]))
		// if err != nil {
		// 	log.Println("Error converting string to int: ", err)
		// 	break
		// }

		// log.Println("Action: ", m.Action)

		// HandleMessage(conn, &m, u, s)
		s.Write([]byte("Testing"))
		return
		// fmt.Println(m.Action)

	}
}

func RandString(length int) string {
	letterBytes := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int64()%int64(len(letterBytes))]
	}
	return string(b)
}

// var src = rand.Source(time.Now().UnixNano())
// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// const (
//     letterIdxBits = 6                    // 6 bits to represent a letter index
//     letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
//     letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
// )
// func RandStringBytesMaskImprSrcUnsafe(n int) string {
//     b := make([]byte, n)
//     // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
//     for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
//         if remain == 0 {
//             cache, remain = src.Int63(), letterIdxMax
//         }
//         if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
//             b[i] = letterBytes[idx]
//             i--
//         }
//         cache >>= letterIdxBits
//         remain--
//     }

//     return *(*string)(unsafe.Pointer(&b))
// }

func main() {
	// log.Println(test)
	// Load certificate and key from files
	cert, err := loadCertificate("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Certificate loaded")

	// Load private key
	key, err := loadPrivateKey("key.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Key loaded")

	// Create tls.Config with desired certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert},
			PrivateKey:  key,
		}},
		InsecureSkipVerify: true, // Remove for production!
	}

	var timeOut time.Duration = 3 * time.Second

	// Create quic.Config with desired versions
	quicConfig := &quic.Config{
		Versions:             []quic.Version{quic.Version2, quic.Version1},
		HandshakeIdleTimeout: timeOut,
		MaxIdleTimeout:       timeOut,
	}

	// Create quic.Listener & listen on UDP port 7000
	listener, err := quic.ListenAddr("127.0.0.1:7000", tlsConfig, quicConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port 7000")

	// make a channel for communication between goroutines
	channel := make(chan string)
	userMap := map[string]string{}

	for {
		ctx := context.Background()
		conn, err := listener.Accept(ctx)
		if err != nil {
			log.Println("Error accepting connection:", err)
			break
		}
		// generates a nanoid when a connection is accepted
		nanoId, err := gonanoid.New(6)
		if err != nil {
			log.Fatal(err)
		}
		// adds the connection to the map to make communication more secure using nanoId
		userMap[conn.RemoteAddr().String()] = nanoId
		log.Println(userMap)
		log.Println(conn.RemoteAddr(), " -> ", nanoId)
		// defer conn.CloseWithError(0, "closing connection")

		go HandleConnection(conn, ctx, channel)

		// HandleConnection(conn, channel)

		// channel <- fmt.Sprint("Testing channel...")
		// fmt.Println("Channel len:", len(channel))

		select {
		case m := <-channel:
			log.Println(m)
		default:
			// log.Println("default case")
			// log.Println("Players: ", len(u.players))
			// break
			// no message continue server
		}
		// for res := range channel {
		// 	log.Println(res)
		// 	fmt.Print(res)
		// }

		// this exits the loop when HandleConnection returns
		// return
	}

	// log.Println("Players: ", len(u.players))
	log.Print("Server closed")
}

func loadCertificate(filename string) ([]byte, error) {
	pemBytes, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	cert, _ := pem.Decode(pemBytes)
	if cert == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	return cert.Bytes, nil
}

func loadPrivateKey(filename string) (any, error) {
	pemBytes, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	key, _ := pem.Decode(pemBytes)
	if key == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	// privateKey, err := x509.ParsePKCS1PrivateKey(key.Bytes)
	privateKey, err := x509.ParsePKCS8PrivateKey(key.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func readFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return data, nil
}
