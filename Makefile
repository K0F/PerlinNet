# Cesty k Android NDK
#ANDROID_NDK_HOME ?= /path/to/your/android/ndk

# Název výstupního souboru
OUTPUT = PerliNet

# Výchozí cíle
.PHONY: all clean

all: linux windows macos #android-armv7 android-arm64

linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)-linux main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT).exe main.go
	zip $(OUTPUT)-windows.zip $(OUTPUT).exe
	rm $(OUTPUT).exe

macos:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT)-macos main.go

#android-armv7:
#	GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=$(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi21-clang go build -o $(OUTPUT)-armv7 main.go

#android-arm64:
#	GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC=$(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang go build -o $(OUTPUT)-arm64 main.go

clean:
	rm -f $(OUTPUT)-linux $(OUTPUT).exe $(OUTPUT)-macos #$(OUTPUT)-armv7 $(OUTPUT)-arm64
