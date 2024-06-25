# Cesty k Android NDK
#ANDROID_NDK_HOME ?= /path/to/your/android/ndk

# Podmíněná proměnná
IS_RASPBERRY_PI := $(shell grep -q "Raspberry Pi" /proc/cpuinfo && echo true)


# Název výstupního souboru
OUTPUT = PerliNet

# Výchozí cíle
.PHONY: clean all zip

all:
	cd beep; make
	cd ..
	go mod tidy
	go test
	go build

cross: raspberrypi linux windows macos


linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT)-linux main.go

windows64:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT).exe main.go
	zip $(OUTPUT)-windows64.zip $(OUTPUT).exe
	rm $(OUTPUT).exe

windows32:
	GOOS=windows GOARCH=386 go build -o $(OUTPUT).exe main.go
	zip $(OUTPUT)-windows32.zip $(OUTPUT).exe
	rm $(OUTPUT).exe

macos:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT)-macos main.go

raspberrypi:
	GOOS=linux GOARCH=arm GOARM=6 go build -o $(OUTPUT)-armv6l main.go

#android-armv7:
#	GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=$(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi21-clang go build -o $(OUTPUT)-armv7 main.go

#android-arm64:
#	GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC=$(ANDROID_NDK_HOME)/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang go build -o $(OUTPUT)-arm64 main.go

clean:
	rm -f $(OUTPUT)-linux $(OUTPUT)-windows64.zip $(OUTPUT)-windows32.zip $(OUTPUT)-armv6l $(OUTPUT)-macos #$(OUTPUT)-armv7 $(OUTPUT)-arm64
