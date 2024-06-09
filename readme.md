# 减小编译体积

go build -ldflags="-s -w" -o server main.go

-s：忽略符号表和调试信息。  
-w：忽略DWARFv3调试信息，使用该选项后将无法使用gdb进行调试。

带壳压缩  
upx -9 server

