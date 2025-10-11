#!/bin/bash

# 服务器管理脚本 - 构建、启动、停止服务器

# 应用名称
APP_NAME="9zhou-server"
# 可执行文件路径
EXECUTABLE="./bin/server"
# PID 文件路径
PID_FILE="./server.pid"
# 端口
PORT=8080

# 构建应用
build() {
    echo "开始构建应用..."
    
    # 创建 bin 目录（如果不存在）
    mkdir -p ./bin
    
    # 构建服务器可执行文件
    go build -o ${EXECUTABLE} ./cmd/server
    
    if [ $? -eq 0 ]; then
        echo "构建成功: ${EXECUTABLE}"
    else
        echo "构建失败"
        exit 1
    fi
}

# 启动应用
start() {
    # 检查是否已在运行
    if [ -f ${PID_FILE} ]; then
        PID=$(cat ${PID_FILE})
        if ps -p ${PID} > /dev/null; then
            echo "应用已在运行 (PID: ${PID})"
            return 1
        else
            # PID文件存在但进程已终止，移除PID文件
            rm ${PID_FILE}
        fi
    fi
    
    echo "正在启动应用..."
    
    # 检查可执行文件是否存在
    if [ ! -f ${EXECUTABLE} ]; then
        echo "可执行文件不存在，请先构建应用"
        exit 1
    fi
    
    # 启动应用并将其放到后台运行
    nohup ${EXECUTABLE} > server.log 2>&1 &
    SERVER_PID=$!
    
    # 将PID保存到文件中
    echo ${SERVER_PID} > ${PID_FILE}
    
    # 等待几秒钟让服务器启动
    sleep 3
    
    # 检查应用是否成功启动
    if ps -p ${SERVER_PID} > /dev/null; then
        echo "应用启动成功 (PID: ${SERVER_PID})"
        echo "日志文件: server.log"
        echo "PID文件: ${PID_FILE}"
    else
        echo "应用启动失败"
        rm -f ${PID_FILE}
        exit 1
    fi
}

# 停止应用
stop() {
    if [ ! -f ${PID_FILE} ]; then
        echo "应用未在运行或PID文件不存在"
        return 1
    fi
    
    PID=$(cat ${PID_FILE})
    
    if ps -p ${PID} > /dev/null; then
        echo "正在停止应用 (PID: ${PID})..."
        kill ${PID}
        
        # 等待应用完全停止
        while ps -p ${PID} > /dev/null; do
            sleep 1
        done
        
        echo "应用已停止"
    else
        echo "进程不存在，清理PID文件"
    fi
    
    # 移除PID文件
    rm -f ${PID_FILE}
}

# 查看状态
status() {
    if [ -f ${PID_FILE} ]; then
        PID=$(cat ${PID_FILE})
        if ps -p ${PID} > /dev/null; then
            echo "应用正在运行 (PID: ${PID})"
        else
            echo "应用未在运行，但存在PID文件"
        fi
    else
        echo "应用未在运行"
    fi
}

# 显示使用说明
usage() {
    echo "使用方法: $0 {build|start|stop|status}"
    echo ""
    echo "命令:"
    echo "  build   - 构建服务器应用"
    echo "  start   - 启动服务器应用"
    echo "  stop    - 停止服务器应用"
    echo "  status  - 查看服务器状态"
    echo ""
}

# 主程序入口
case "$1" in
    build)
        build
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    *)
        usage
        exit 1
        ;;
esac

exit 0