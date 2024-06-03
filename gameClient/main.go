package main

import (
	"gameClient/ECS"
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	ECS.Player
	Id   int64
	Name string
	// Position rl.Vector2
	// Velocity rl.Vector2
}

func FragmentShader(source string) string {
	return source
}

func (p *Player) Update(deltaTime float32) {
	// Calculate the magnitude (length) of the velocity vector
	magnitude := float32(math.Sqrt(float64(p.Velocity.X*p.Velocity.X + p.Velocity.Y*p.Velocity.Y)))

	// Normalize the velocity vector if magnitude is not zero
	if magnitude != 0 {
		p.Velocity.X /= magnitude
		p.Velocity.Y /= magnitude
	}

	// Update position using velocity, speed, and deltaTime
	p.Position.X += p.Velocity.X * p.Speed * deltaTime
	p.Position.Y += p.Velocity.Y * p.Speed * deltaTime
	rl.Vector2Clamp(p.Velocity, rl.NewVector2(-p.MaxVel, -p.MaxVel), rl.NewVector2(p.MaxVel, p.MaxVel))
	// rl.Vector2Normalize(p.Velocity)
}

func main() {
	const Scale = 2
	const targetFps = 165
	image := rl.LoadImage("./momodora.png")
	rl.ImageResizeNN(image, int32(image.Width*Scale), int32(image.Height*Scale))

	player := Player{
		Player: ECS.Player{
			Transform: ECS.Transform{
				Velocity: rl.NewVector2(0, 0),
				Position: rl.NewVector2(20, 20),
			},
			Speed:  200.0,
			MaxVel: 350.0,
			// Texture: rl.LoadTexture("./momodora.png"),
			// Texture: rl.LoadTextureFromImage(image),
		},
		Id:   1,
		Name: "testogus",
	}
	defer rl.UnloadImage(image)

	const tileCount = 6
	const animationSpeed = 6
	var currentTile int = 0
	var tileCounter int = 0
	// var animSheet[...]int use [...] for dynamic size
	// var animSheet = make([]rl.Rectangle, tileCount)

	// var animSheet = [4]int{6, 3, 5, 3}
	var scaledSprite rl.Rectangle = rl.NewRectangle(0, 0, float32(48*Scale), float32(48*Scale))
	// this is possibly dynamic
	// var scaledSprite rl.Rectangle = rl.NewRectangle(0, 0, player.Texture.Width/animSheet[0], player.Texture.Height/len(animSheet))

	rl.InitWindow(800, 450, "Game Client - basic window")
	defer rl.CloseWindow()

	player.Texture = rl.LoadTextureFromImage(image)
	rl.SetTargetFPS(targetFps)

	log.Println("Using OpenGl version", rl.GetVersion())

	// QUIC CONFIG BEGIN
	// quicConfig := &quic.Config{
	// 	Versions: []quic.Version{quic.Version2, quic.Version1},
	// }
	// // tls config
	// tlsConfig := &tls.Config{
	// 	ClientSessionCache: tls.NewLRUClientSessionCache(100),
	// 	InsecureSkipVerify: true,
	// }
	// // initialize quic context
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
	// defer cancel()
	// // initialize connection
	// // conn, err := quic.Dial(ctx, listener, addr, tlsConfig, quicConfig)
	// conn, err := quic.DialAddr(ctx, "127.0.0.1:7000", tlsConfig, quicConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Remember to close the connection when finished
	// defer conn.CloseWithError(0, "Connection terminated")

	// // open stream
	// stream, err := conn.OpenStream()
	// defer stream.Close()
	// log.Println("Using TLS v", conn.ConnectionState().TLS.Version)
	// if err != nil {
	// 	log.Fatalln("Failed to open stream: ", err)
	// }
	// log.Println("Stream opened")

	// // send data using stream
	// testVector := rl.Vector2{100, 100}
	// sendData := []byte("0|test" + testVector)
	// // td := Test{
	// // 	data: "jkDN",
	// // }
	// // sendData := []byte(td)
	// _, err = stream.Write(sendData)
	// if err != nil {
	// 	fmt.Println("Error sending data: ", err)
	// }

	// log.Println("Sent", []byte("12"))

	// // read data using stream
	// buf := make([]byte, 1024)
	// n, err := stream.Read(buf)
	// if err != nil {
	// 	log.Fatal("Error reading data: ", err)
	// }
	// log.Println("Data received")
	// // Process the received data
	// fmt.Printf("Received %d bytes from server: %s\n", n, string(buf[:n]))
	// // Implement logic to handle receiving data from the server (if applicable)
	// fmt.Println("Connection closed")
	// QUIC CONFIG END

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		tileCounter += 1
		if tileCounter >= (targetFps / animationSpeed) {
			tileCounter = 0
			currentTile += 1
			if currentTile >= tileCount-1 {
				currentTile = 0
			}
			scaledSprite.X = float32(currentTile) * scaledSprite.Width
			// scaledSprite.X = (float32(currentTile) * scaledSprite.Width)

			// currentTile = (currentTile + 1) % tileCount ?

		}

		player.Velocity = rl.Vector2{0, 0}

		// Update player velocity
		if rl.IsKeyDown(rl.KeyRight) {
			player.Velocity.X = 1
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			player.Velocity.X = -1
		}
		if rl.IsKeyDown(rl.KeyUp) {
			player.Velocity.Y = -1
		}
		if rl.IsKeyDown(rl.KeyDown) {
			player.Velocity.Y = 1
		}

		// stops player when key released
		// if rl.IsKeyReleased(rl.KeyRight) || rl.IsKeyReleased(rl.KeyLeft) {
		// 	player.Velocity.X = 0
		// }

		// if rl.IsKeyReleased(rl.KeyUp) || rl.IsKeyReleased(rl.KeyDown) {
		// 	player.Velocity.Y = 0
		// }

		// // Normalize velocity to ensure consistent diagonal movement
		// if player.Velocity.X != 0 && player.Velocity.Y != 0 {
		// 	player.Velocity = rl.Vector2Normalize(player.Velocity)
		// }

		// // Update player position
		player.Update(deltaTime)
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		if player.Velocity.X < 0 {
			scaledSprite.Width = -48 * float32(Scale)
		} else {
			scaledSprite.Width = 48 * float32(Scale)
		}

		rl.DrawTextureRec(player.Texture, scaledSprite, player.Position, rl.White)

		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.Maroon)

		rl.EndDrawing()
	}

	// rl.UnloadTexture(testTex)
}
