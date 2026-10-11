@echo off
setlocal
where go >nul 2>nul || (echo Instale Go 1.20.x para Windows 32 bits antes de compilar. & pause & exit /b 1)
for /f "tokens=3" %%v in ('go version') do set GOVERSION=%%v
echo Versao: %GOVERSION%
echo IMPORTANTE: apenas Go 1.20.x e compativel com Windows 7.
echo Compilacao recomendada usando Go 1.20.x no Windows 7 ou ambiente compativel.
set GOOS=windows
set GOARCH=386
set CGO_ENABLED=0
cd source
go build -trimpath -o ..\bin\NatorVM-Win7-32.exe .
if errorlevel 1 (echo Falha na compilacao & pause & exit /b 1)
echo Concluido. Teste antes de usar discos com dados importantes.
pause
