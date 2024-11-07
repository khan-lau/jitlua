.PHONY: all build clean run

DST_DIR=dist
PWD=

BIN_FILE=jitlua_example
MAIN_PROG= cmd/demo.go

Version=0.1.1
Author=Liu Kun

DEBUG=-w -s
param=-X main.BuildVersion=${Version} -X \"main.BuildTime=${BuildDate}\" -X \"main.BuildPerson=${Author}\"



ifeq ($(OS),Windows_NT) 
	uname_S=Windows

    # 判断是否存在uname命令
	tmp_uname=$(shell where uname 2>NUL)
	ifneq ($(tmp_uname),)
		term_S=Windows_Mingw
	endif
	BuildDate=$(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd HH:mm:ss'")
	PWD=$(shell cmd /C 'chdir')
else
	uname_S=$(shell uname -s)
	BuildDate=$(shell date +"%F %T")
	PWD=$(shell pwd)
endif

# $(warning OS: $(OS) - ${uname_S} , ${param}, ${PWD})
# $(warning VAR: $(RM), $(TARGET_OS), ${term_S})

all: build


build:
ifeq ($(uname_S), Windows)
	@cmd /C 'set CGO_ENABLED=1&&set GOOS=windows&&go build -v -ldflags '${DEBUG} ${param}' -o ${DST_DIR}/${BIN_FILE}.exe  ${MAIN_PROG}'
	@powershell -Command "Write-Host \"Build windows 64bit program - ${DST_DIR}/${BIN_FILE}.exe\" -ForegroundColor green"
endif
# ifeq ($(uname_S), Linux)
# 	@export CGO_ENABLED=1; export GOOS=windows; go build -v -ldflags "${DEBUG} ${param}" -o ${DST_DIR}/${BIN_FILE}.exe  ${MAIN_PROG}
# 	@echo -e "\033[0;32m Build windows 64bit program - ${DST_DIR}/${BIN_FILE}.exe\033[0m"
# endif
	

# ifeq ($(uname_S), Windows)
# 	@cmd /C 'set CGO_ENABLED=0&&set GOOS=linux&&go build -v -ldflags '${DEBUG} ${param}' -o ${DST_DIR}/${BIN_FILE}  ${MAIN_PROG}'
# 	@powershell -Command "Write-Host \"Build linux 64bit program - ${DST_DIR}/${BIN_FILE}\" -ForegroundColor green"
# endif
ifeq ($(uname_S), Linux)
	@export CGO_ENABLED=0; export GOOS=linux; go build -v -ldflags "${DEBUG} ${param}" -o ${DST_DIR}/${BIN_FILE}  ${MAIN_PROG}
	@echo -e "\033[0;32m Build linux 64bit program - ${DST_DIR}/${BIN_FILE}\033[0m"
endif
	


# make ARGS="-v" run
run:
ifeq ($(uname_S), Windows)
	${DST_DIR}/${BIN_FILE}.exe $(ARGS)
endif
ifeq ($(uname_S), Linux)
	${DST_DIR}/${BIN_FILE} $(ARGS)
endif


check:
# 检查是否已安装 golangci-lint
ifeq ($(shell where golangci-lint 2>NUL), )
	$(error golangci-lint not found!  try run 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest')
endif

# 格式化代码
ifeq ($(uname_S), Windows)
	
ifeq ($(term_S), Windows_Mingw)
	@rm -f NUL
	@powershell -Command "Get-ChildItem -Path \"${PWD}\" -recurse *.go |ForEach-Object {go fmt \$$_.FullName}"
else
	@powershell -Command "Get-ChildItem -Path \"${PWD}\" -recurse *.go |ForEach-Object {go fmt $$_.FullName}"
endif

endif

ifeq ($(uname_S), Linux)
	@find ./ -name '*.go' -exec go fmt {} \;
endif

# 语法检查
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run


clean:
	@go clean
ifeq ($(uname_S), Windows)
	@cmd /C 'del /F /Q ${DST_DIR}\\*'
endif
ifeq ($(uname_S), Linux)
	@rm --force ${DST_DIR}/*
endif
