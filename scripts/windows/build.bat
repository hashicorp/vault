@echo off
setlocal

set _EXITCODE=0
set _DEV_BUILD=0

if not exist %1 exit /b 1
if x%2 == xVAULT_DEV set _DEV_BUILD=1

cd %1
md bin 2>nul

:: Get the git commit
set _GIT_COMMIT_FILE=%TEMP%\vault-git_commit.txt
set _GIT_DIRTY_FILE=%TEMP%\vault-git_dirty.txt

set _NUL_CMP_FILE=%TEMP%\vault-nul_cmp.txt
type nul >%_NUL_CMP_FILE%

git rev-parse HEAD >"%_GIT_COMMIT_FILE%"
set /p _GIT_COMMIT=<"%_GIT_COMMIT_FILE%"
del /f "%_GIT_COMMIT_FILE%" 2>nul

set _GIT_DIRTY=
git status --porcelain >"%_GIT_DIRTY_FILE%"
fc "%_GIT_DIRTY_FILE%" "%_NUL_CMP_FILE%" >nul
if errorlevel 1 set _GIT_DIRTY=+CHANGES
del /f "%_GIT_DIRTY_FILE%" 2>nul
del /f "%_NUL_CMP_FILE%" 2>nul

REM Determine the arch/os combos we're building for
set _XC_ARCH=386 amd64 arm
set _XC_OS=linux darwin windows freebsd openbsd

REM Install dependencies
echo ==^> Installing dependencies...
go get ./...

REM Clean up the old binaries and packages.
echo ==^> Cleaning old builds...
rd /s /q bin pkg 2>nul
md bin 2>nul

REM If its dev mode, only build for ourself
if not %_DEV_BUILD% equ 1 goto build

:devbuild
echo ==^> Preparing for development build...
set _GO_ENV_TMP_FILE=%TEMP%\vault-go-env.txt
go env GOARCH >"%_GO_ENV_TMP_FILE%"
set /p _XC_ARCH=<"%_GO_ENV_TMP_FILE%"
del /f "%_GO_ENV_TMP_FILE%" 2>nul
go env GOOS >"%_GO_ENV_TMP_FILE%"
set /p _XC_OS=<"%_GO_ENV_TMP_FILE%"
del /f "%_GO_ENV_TMP_FILE%" 2>nul

:build
REM Build!
echo ==^> Building...
gox^
 -os="%_XC_OS%"^
 -arch="%_XC_ARCH%"^
 -ldflags "-X github.com/hashicorp/vault/version.GitCommit=%_GIT_COMMIT%%_GIT_DIRTY%"^
 -output "pkg/{{.OS}}_{{.Arch}}/vault"^
 .

if %ERRORLEVEL% equ 1 set %_EXITCODE%=1

if %_EXITCODE% equ 1 exit /b %_EXITCODE%

set _GO_ENV_TMP_FILE=%TEMP%\vault-go-env.txt

go env GOPATH >"%_GO_ENV_TMP_FILE%"
set /p _GOPATH=<"%_GO_ENV_TMP_FILE%"
del /f "%_GO_ENV_TMP_FILE%" 2>nul

go env GOARCH >"%_GO_ENV_TMP_FILE%"
set /p _GOARCH=<"%_GO_ENV_TMP_FILE%"
del /f "%_GO_ENV_TMP_FILE%" 2>nul

go env GOOS >"%_GO_ENV_TMP_FILE%"
set /p _GOOS=<"%_GO_ENV_TMP_FILE%"
del /f "%_GO_ENV_TMP_FILE%" 2>nul

REM Copy our OS/Arch to the bin/ directory
set _DEV_PLATFORM=pkg\%_GOOS%_%_GOARCH%

for /r %%f in (%_DEV_PLATFORM%) do (
	copy /b /y %%f bin\ >nul
	copy /b /y %%f %_GOPATH%\bin\ >nul
)

REM TODO(ceh): package dist

REM Done!
echo.
echo ==^> Results:
echo.
for %%A in ("bin\*") do echo %%~fA

exit /b %_EXITCODE%
