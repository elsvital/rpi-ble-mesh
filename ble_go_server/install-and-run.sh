#!/bin/bash

echo "🔧 Instalando dependências do sistema para Go + BlueZ"
sudo apt update
sudo apt install -y libglib2.0-dev libdbus-1-dev bluetooth bluez golang

echo "📦 Inicializando módulo Go (caso ainda não esteja inicializado)"
if [ ! -f "go.mod" ]; then
  go mod init ble_go_server
fi

echo "📦 Instalando pacotes Go"
go get github.com/dgrijalva/jwt-go
go get github.com/godbus/dbus/v5
go get github.com/segmentio/kafka-go
go mod tidy

echo "✅ Ambiente Go configurado com sucesso."

echo "⚙️ Compilando aplicação BLE..."
go build -o ble_server ./cmd

if [ $? -eq 0 ]; then
  echo "✅ Compilação bem-sucedida. Iniciando servidor..."
  sudo ./ble_server
else
  echo "❌ Erro na compilação. Verifique as mensagens acima."
fi