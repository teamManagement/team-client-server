#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o teamClientServer main.go > /dev/null 2>&1
CC=x86_64-w64-mingw32-gcc-posix CXX=x86_64-w64-mingw32-g++-posix GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o teamClientServer.exe main.go > /dev/null 2>&1

externalProgramDir="externalProgram"

processDir="process/"
processLinuxDir="$processDir""linux/"
processWinDir="$processDir""win/"

linuxMkcertFilePath="$processLinuxDir""mkcert"
winMkcertFilePath="$processWinDir""mkcert.exe"

clientProcessName="teamClientServer"
clientWinProcessName="teamClientServer".exe


rm -rf $processDir

cp -rf $externalProgramDir $processDir

mv $clientProcessName $processLinuxDir
mv $clientWinProcessName $processWinDir


linuxClientServerSha512Sum="$(sha512sum $processLinuxDir/$clientProcessName | cut -d ' ' -f1)"
winClientServerSha512Sum="$(sha512sum $processWinDir/$clientWinProcessName | cut -d ' ' -f1)"

linuxMkcertSha512Sum="$(sha512sum $linuxMkcertFilePath | cut -d ' ' -f1)"
winMkCertSha512Sum="$(sha512sum $winMkcertFilePath | cut -d ' ' -f1)"


jsonResult=$(cat <<- EOF
 "signature": {
    "mkcert": {
      "linux": "$linuxMkcertSha512Sum",
      "win32": "$winMkCertSha512Sum"
    },
    "localServer": {
      "linux": "$linuxClientServerSha512Sum",
      "win": "$winClientServerSha512Sum"
    }
  },
EOF
)

echo "$jsonResult"



