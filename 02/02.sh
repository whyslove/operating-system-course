#!/bin/bash

set -e

sudo cgcreate -g cpu,memory:test

# Запускаем ресурсоемкий процесс
stress-ng --cpu 1  --timeout 8s &

# Получаем PID процесса
STRESS_PID=$!

echo "My PID $STRESS_PID"

# Проверяем, что процесс жив
if kill -0 $STRESS_PID 2>/dev/null; then
    echo "Process $STRESS_PID is alive."
else
    echo "Process $STRESS_PID is not running."
    exit 1
fi

# Устанавливаем ограничение по CPU
sudo cgset -r cpu.max="2000 1000000" test  # Ограничение CPU до 2000 микро (буква мю)  на 1 секунду

# Присоеднияем процесс к cgroup
sudo cgclassify -g cpu:test $STRESS_PID


# Пусть процесс поработает
sleep 5

# Посмотрим на троттлинг
cgget -r cpu.stat test


stress-ng  --vm 1 --vm-bytes 100M --timeout 8s --vm-keep &
STRESS_PID=$!
echo "My PID $STRESS_PID"


# Устанавливаем ограничение по памяти, чтобы вызывать ООМ
sudo cgset -r memory.max=1M test
sudo cgset -r memory.high=1M test
sudo cgset -r memory.swap.max=1M test
sudo cgclassify -g memory:test $STRESS_PID

# Ждем, пока процесс будет убит
while kill -0 $STRESS_PID 2>/dev/null; do sleep 1; done
echo "Finished."

# Удаляем cgroup
sudo cgdelete -g cpu,memory:test

