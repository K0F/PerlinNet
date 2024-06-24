#define MINIAUDIO_IMPLEMENTATION
#include "miniaudio.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char** argv)
{
    if (argc != 2) {
        printf("Usage: %s <path to sound file>\n", argv[0]);
        return -1;
    }

    const char* filePath = argv[1];
    ma_result result;
    ma_engine engine;
    ma_sound sound;

    result = ma_engine_init(NULL, &engine);
    if (result != MA_SUCCESS) {
        return -1;
    }

    ma_result err = ma_sound_init_from_file(&engine, filePath, 0, NULL, NULL, &sound);
    if (err != MA_SUCCESS) {
        ma_engine_uninit(&engine);
        return result;
    }

    ma_sound_start(&sound);

    while(ma_sound_is_playing(&sound)){
        sleep(1);
    }

    ma_engine_uninit(&engine);

    return 0;
}
