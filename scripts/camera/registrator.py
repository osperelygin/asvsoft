from picamera2 import Picamera2
from PIL import Image
from io import BytesIO

import time
import sys
import libcamera
import numpy as np
import socket
import signal
import os

class Registrator:
    def __init__(self, socket_path='/tmp/camera.sock', resolution=(160, 120)):
        self.socket_path = socket_path
        self.running = True
        self.sock = None
        
        # Регистрируем обработчики сигналов
        signal.signal(signal.SIGTERM, self.handle_signal)
        signal.signal(signal.SIGINT, self.handle_signal)

        # Инициализация камеры
        self.picam2 = Picamera2()

        # Настройка конфигурации для захвата изображения
        config = self.picam2.create_still_configuration(
            main={"size": resolution},
            # Опционально: отражение по горизонтали и вертикали
            transform=libcamera.Transform(hflip=1, vflip=1)
        )
        self.picam2.configure(config)

        # Запуск камеры
        self.picam2.start()

        # Даём камере немного времени для настройки
        time.sleep(2)
    
    def handle_signal(self, signum, frame):
        print(f"Received signal {signum}, shutting down...")
        self.running = False
        if self.sock:
            self.sock.close()

        # Останавливаем камеру
        if self.picam2:
            self.picam2.stop()
            self.picam2.close()

        sys.exit(0)
    
    def connect(self):
        """Устанавливаем соединение с сокетом"""
        while self.running:
            try:
                self.sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
                self.sock.connect(self.socket_path)
                print(f"Connected to {self.socket_path}")
                return True
            except socket.error as e:
                print(f"Connection error: {e}, retrying...")
                time.sleep(1)
        return False
    
    def send_data(self, data):
        """Отправляем данные через сокет"""
        try:
            self.sock.sendall(data)
            print(f"Sent: {data}")
        except socket.error as e:
            print(f"Send error: {e}, reconnecting...")
            self.sock.close()
            self.connect()
    
    def run(self):
        """Основной цикл отправки данных"""
        if not self.connect():
            return
        
        while self.running:
            try:
                data = self.take_monochrome_image()
                self.send_data(data)
                time.sleep(1)
            except KeyboardInterrupt:
                self.running = False
        
        self.sock.close()

    def take_monochrome_image(self):
        # Захватываем изображение в виде numpy-массива
        array = self.picam2.capture_array()

        # Преобразуем цветное изображение в чёрно-белое (grayscale)
        gray_array = np.dot(array[..., :3], [0.2989, 0.5870, 0.1140]).astype(np.uint8)

        # Кодируем изображение
        gray_image = Image.fromarray(gray_array)
        with BytesIO() as output:
            gray_image.save(output, format="JPEG")
            jpeg_bytes = output.getvalue()

        return jpeg_bytes


if __name__ == "__main__":
    sender = Registrator()
    sender.run()
