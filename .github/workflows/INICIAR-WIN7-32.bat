@echo off
cd /d "%~dp0"
if not exist "bin\NatorVM-Win7-32.exe" (echo Compile primeiro com Go 1.20. & pause & exit /b 1)
set PATH=%CD%\qemu;%PATH%
"bin\NatorVM-Win7-32.exe"
pause
