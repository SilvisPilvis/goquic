#include "raylib.h"
#include "raymath.h"
#include <stdio.h>
#include <iostream>
#include <sys/resource.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>

// #include <stdio.h>
// #include <unistd.h>
// #include <glm/glm.hpp>

typedef struct Player {
    int64_t id;
    Texture2D texture;
    Vector2 position{200.0f, 200.0f};
    Vector2 velocity{0.0f, 0.0f};
    float maxVelocity = 5000.0f;
    float speed;
    const int vOffset = 2;
    
    // constructor below
    Player(Texture2D texture, float spd) : texture(texture), speed(spd) {}
} Player;

typedef struct Options {
    bool telemetry = true;
    bool lockedScroll = false;
} Options;

typedef enum State {
    STATE_PAUSED = 0,
    STATE_PLAYING = 1,
    STATE_GAMEOVER = 2
} State;

struct MemoryInfo {
  long residentSetSize;  // Resident Set Size (in bytes)
};

MemoryInfo getMemoryUsage() {
  MemoryInfo info;
  struct rusage usage;
  if (getrusage(RUSAGE_SELF, &usage) == -1) {
    std::cerr << "Error: Failed to call getrusage" << std::endl;
    return info;
  }

  info.residentSetSize = usage.ru_maxrss;  // Maximum resident set size
  return info;
//   std::string memInfo = std::to_string(info.residentSetSize);
//   return memInfo;
}

// Module functions declaration
void UpdatePlayer(Player *player, float delta);
void UpdateCamera2D(Camera2D *camera, Player *player, int width, int height);

// const std::string formatString = "MEM: %ldB";



int main(void){

    // int lsquic_engine_packet_in (lsquic_engine_t *,
    // const unsigned char *udp_payload, size_t sz,
    // const struct sockaddr *sa_local,
    // const struct sockaddr *sa_peer,
    // void *peer_ctx, int ecn);

    // make variables for game window
    const int screenWidth = 800;
    const int screenHeight = 450;

    // init variables for fps
    double lastFrameTime = 0.0; // Variable to store time of the last frame
    int frames = 0; // Counter for the number of frames rendered

    // allow resizing
    SetConfigFlags(FLAG_WINDOW_RESIZABLE);

    // init game window
    InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window");

    // initialize options & state
    Options options;
    State state = STATE_PLAYING;

    // make a player & load texture
    const float scale = 2.0f;
    Image playerTex = LoadImage("momodora-sheet.png");
    ImageResizeNN(&playerTex, playerTex.width*scale, playerTex.height*scale);
    Player player{LoadTextureFromImage(playerTex), 7.0f};
    // unload image from RAM into VRAM
    UnloadImage(playerTex);

    // create a rectangle the size of Spritesheet 1 tile and scale it
    Rectangle scaledBounds = (Rectangle){0, 0, 48*scale, 48*scale}; // size of one animation tile scaled

    // initialize camera
    Camera2D camera = {0};
    camera.target = (Vector2){player.position.x/2.0f, player.position.y/2.0f}; // player.position;
    camera.offset = (Vector2){screenWidth/2.0f, screenHeight/2.0f};
    camera.rotation = 0.0f;
    camera.zoom = 1.0f;

    // set target fps
    // int targetFps = 60;
    int targetFps = 165;
    SetTargetFPS(targetFps);
    int gamepad = 0;

    while (!WindowShouldClose())
    {
        float delta = GetFrameTime();
        // if(IsGamepadAvailable(gamepad)) printf("Gamepad: %d\n", gamepad);

        MemoryInfo memInfo = getMemoryUsage();

        if(IsKeyPressed(KEY_T)) options.telemetry = !options.telemetry;

        // update player
        UpdatePlayer(&player, delta);
        if(!options.lockedScroll){
            camera.zoom += ((float)GetMouseWheelMove()*0.05f);
        }
        UpdateCamera2D(&camera, &player, screenWidth, screenHeight);

        // Draw here
        BeginDrawing();
            BeginMode2D(camera);
            // clear bg
            ClearBackground(RAYWHITE);
            // draw floor
            DrawRectangle(0, GetScreenHeight() - 64, GetScreenWidth(), (GetScreenHeight() - (GetScreenHeight() - 64)), MAROON);
            // draw player
            if (player.velocity.x < 0){
                scaledBounds.width = -48*scale;
            }else if (player.velocity.x > 0){
                scaledBounds.width = 48*scale;
            }
            DrawTextureRec(player.texture, scaledBounds, player.position, WHITE);
            EndMode2D();
 

        // display fps
        frames++;
        if(delta >= 1.0f / (float)targetFps){
            int fps = (int)(1.0 / delta);
            char fpsText[10];
            sprintf(fpsText, "%d", fps);
            if (options.telemetry) DrawText(fpsText, 10, 10, 14, MAROON);
            frames = 0;
        }

        std::string text = "MEM: " + std::to_string(memInfo.residentSetSize) + " B";
        // if (options.telemetry) DrawText(text.c_str(), 10, 24, 14, MAROON);
        std::string velX = "VELOCITY: " + std::to_string(player.velocity.x) + " " + std::to_string(player.velocity.y) + "px/s";
        // if (options.telemetry) DrawText(velX.c_str(), 10, 38, 14, MAROON);
        std::string posX = "X: " + std::to_string(player.position.x);
        // if (options.telemetry) DrawText(posX.c_str(), 10, 52, 14, MAROON);
        std::string posY = "Y: " + std::to_string(player.position.y);
        // if (options.telemetry) DrawText(posY.c_str(), 10, 68, 14, MAROON);
        std::string gameState = "STATE: " + std::to_string(state);
        if(options.telemetry){
            DrawText(text.c_str(), 10, 24, 14, MAROON);
            DrawText(velX.c_str(), 10, 38, 14, MAROON);
            DrawText(posX.c_str(), 10, 52, 14, MAROON);
            DrawText(posY.c_str(), 10, 68, 14, MAROON);
            DrawText(gameState.c_str(), 10, 82, 14, MAROON);
        }
        // if (options.telemetry) DrawText(gameState.c_str(), 10, 82, 14, MAROON);
        EndDrawing();
    }

    // De-Initialization
    UnloadTexture(player.texture);
    CloseWindow();

    return 0;
}

void UpdatePlayer(Player *player, float delta)
{
    // player movement
    if (IsKeyDown(KEY_LEFT) || IsKeyDown(KEY_A) || IsGamepadButtonDown(0, GAMEPAD_BUTTON_LEFT_FACE_LEFT)) player->velocity.x -= player->speed;
    if (IsKeyDown(KEY_RIGHT) || IsKeyDown(KEY_D) || IsGamepadButtonDown(0, GAMEPAD_BUTTON_LEFT_FACE_RIGHT)) player->velocity.x += player->speed;
    if (IsKeyDown(KEY_UP) || IsKeyDown(KEY_W) || IsGamepadButtonDown(0, GAMEPAD_BUTTON_LEFT_FACE_UP)) player->velocity.y -= player->speed;
    if (IsKeyDown(KEY_DOWN) || IsKeyDown(KEY_S) || IsGamepadButtonDown(0, GAMEPAD_BUTTON_LEFT_FACE_DOWN)) player->velocity.y += player->speed;


    // stop player if no input
    if(IsKeyReleased(KEY_LEFT) || IsKeyReleased(KEY_RIGHT) || IsKeyReleased(KEY_A) || IsKeyReleased(KEY_D) || IsKeyReleased(KEY_W) || IsKeyReleased(KEY_S) || IsKeyReleased(KEY_UP) || IsKeyReleased(KEY_DOWN) || IsGamepadButtonReleased(0, GAMEPAD_BUTTON_LEFT_FACE_LEFT) || IsGamepadButtonReleased(0, GAMEPAD_BUTTON_LEFT_FACE_RIGHT)) player->velocity = (Vector2){0.0f, 0.0f};
    // if(IsKeyReleased(KEY_SPACE) && player->grounded) player->speed = 0.0f;

    // add a rectangle on the collision
    // DrawRectangle(0, GetScreenHeight() - 64, GetScreenWidth(), (GetScreenHeight() - (GetScreenHeight() - 64)), MAROON);

    // defines gravity
    float gravity = 10.0f * delta;

    // working clamping
    float maxSpeed = player->maxVelocity;
    player->velocity = { Clamp(player->velocity.x, -maxSpeed, maxSpeed), Clamp(player->velocity.y, -maxSpeed, maxSpeed) };

    // Normalize the velocity to ensure consistent movement speed regardless of frame rate
    // float velocityMagnitude = sqrtf(player->velocity.x * player->velocity.x + player->velocity.y * player->velocity.y);
    // if (velocityMagnitude > 0.0f) {
    //     player->velocity = { player->velocity.x / velocityMagnitude * maxSpeed, player->velocity.y / velocityMagnitude * maxSpeed };
    // }

    // // adds velocity to player / moves player
    // // Update player's position based on velocity
    // player->position.x += player->velocity.x * delta;
    // player->position.y += player->velocity.y * delta;
    player->position = Vector2Add(player->position, Vector2Multiply(player->velocity, (Vector2){delta, delta}));


    // player->position = Vector2Normalize(Vector2Add(player->position, Vector2Multiply(player->velocity, (Vector2){delta, delta})));
}

void UpdateCamera2D(Camera2D *camera, Player *player, int width, int height)
{
    camera->offset = (Vector2){width/2.0f, height/2.0f};
    camera->target = player->position;
}
