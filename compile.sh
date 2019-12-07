rsrc -manifest main.manifest -ico ./res/favicon.ico -o rsrc.syso
go build  -ldflags="-s -w -H windowsgui"
tool/upx.exe -9 DVWASolver.exe