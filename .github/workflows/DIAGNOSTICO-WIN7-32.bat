@echo off
setlocal
echo === Diagnostico NatorVM Legacy 10.3 ===
ver
wmic os get OSArchitecture /value 2>nul
if exist "qemu\qemu-system-i386.exe" (echo QEMU 32 bits: presente na pasta qemu) else (echo QEMU 32 bits: nao incluido / nao encontrado na pasta qemu)
if exist "qemu\qemu-img.exe" (echo qemu-img: presente) else (echo qemu-img: nao incluído / nao encontrado na pasta qemu)
if exist "bin\NatorVM-Win7-32.exe" (echo Executavel legado: presente) else (echo Executavel legado: precisa ser compilado com Go 1.20)
echo ALERTA: Windows 7 nao recebe correcoes regulares de seguranca.
pause
