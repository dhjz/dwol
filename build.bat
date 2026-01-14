@echo off
@REM 使用windows图标:  go install github.com/akavel/rsrc@latest     rsrc -ico="ico.ico" -o="ico_windows.syso"

set APP_NAME=dwol
set MAIN_LINUX_ARCH=arm64

go env -w GOOS=linux GOARCH=amd64 GOARM=
go build -ldflags "-s -w" -o ./dist/%APP_NAME%_amd64
echo "build linux_amd64 success..."

go env -w GOOS=linux GOARCH=arm  GOARM=7 
go build -ldflags "-s -w" -o ./dist/%APP_NAME%_armv7
echo "build linux-armv7 success..."

go env -w GOOS=linux GOARCH=arm64  GOARM= 
go build -ldflags "-s -w" -o ./dist/%APP_NAME%_arm64
echo "build linux-arm64 success..."

go env -w GOOS=windows GOARCH=amd64 GOARM=
@REM go build -ldflags "-s -w -H=windowsgui" -o ./dist/
go build -ldflags "-s -w" -o ./dist/%APP_NAME%.exe
echo "build windows exe success..."

del /Q dist\*.zip
for %%f in (dist\*) do (
    if not "%%~xf"==".zip" (
        zip "dist\%%~nf.zip" "%%f"
    )
)

if not "%1"=="true" (
    ren "dist\%APP_NAME%_%MAIN_LINUX_ARCH%" "%APP_NAME%"
)

pause
endlocal