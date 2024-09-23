"""
tof_laser_range.py

This example illustrates how to read measure from TOF Laser Range 
Read more: https://www.waveshare.com/wiki/TOF_Laser_Range_Sensor_(B)

Prerequisites:

python3 -m venv .venv && source .venv/bin/activate && pip3 install pyserial

Usage:

python3 tof_laser_range.py /dev/ttyS0

"""

import logging
import serial
import sys
import time

logger = logging.getLogger('LaserDepthMeter')
logging.basicConfig(level=logging.INFO)

HEADER = b'W\x00\xff'
HEADER_LENGHT = 3
PACKET_LENGTH = 16

ID_LENGTH = 1
SYSTEM_TIME_LEGNTH = 4
DISTNACE_LENGTH = 3
SIGNAL_STATUS_LENGTH = 1
SIGNAL_STRENGTH_LENGTH = 2
SIGNAL_ACCURACY_LENGTH = 1


# readTOFData ...
def readData(ser: serial.Serial) -> bytes:
    rawData = ser.read(2*PACKET_LENGTH)
    if len(rawData) != 2*PACKET_LENGTH:
        logger.warning("Packet not red")
        return None

    idx = rawData.find(HEADER)
    if idx == -1:
        logger.warning("Header not found")
        return None

    if not verifyCheckSum(rawData[idx : idx + PACKET_LENGTH]):
        logger.warning("Wrong check sum")
        return None

    return rawData[idx + HEADER_LENGHT : idx + PACKET_LENGTH]


# verifyCheckSum ...
def verifyCheckSum(rawData: bytes) -> bool:
    payload = tuple(rawData)
    return sum(payload[:-1]) % 256 == payload[-1]


def main():
    if len(sys.argv) < 2:
        logger.error("The first arguments must be the port")
        exit(1)


    port = sys.argv[1]
    ser = serial.Serial(port, 115200, timeout=3)

    try:
        while True:
            payload = readData(ser)
            if payload == None:
                continue

            idx = 0
                
            id = int.from_bytes(payload[idx : idx + ID_LENGTH], 'little')
            idx += ID_LENGTH

            systemTime = int.from_bytes(payload[idx : idx + SYSTEM_TIME_LEGNTH], 'little')
            idx += SYSTEM_TIME_LEGNTH

            rawDistance = payload[idx : idx + DISTNACE_LENGTH]
            distance = int.from_bytes(rawDistance, 'little')
            idx += DISTNACE_LENGTH

            signalStatus = int.from_bytes(payload[idx : idx + SIGNAL_STATUS_LENGTH], 'little')
            idx += SIGNAL_STATUS_LENGTH
            
            signalStrength = int.from_bytes(payload[idx : idx + SIGNAL_STRENGTH_LENGTH], 'little')
            idx += SIGNAL_STRENGTH_LENGTH

            signalAccuracy = int.from_bytes(payload[idx : idx + SIGNAL_ACCURACY_LENGTH], 'little')
            idx += SIGNAL_ACCURACY_LENGTH

            logger.info(f'System time: {systemTime}ms, Distance: {distance}mm, Status: {signalStatus}, Strength: {signalStrength}, Accuracy: {signalAccuracy}')

            # Передаем данные если все ок
            if signalStatus == 1 and signalStrength != 0:
                pass

            time.sleep(0.1)

            # дропаем неактуальные данные
            ser.flushInput()

    except OSError as err:
        logger.error(err)
    except KeyboardInterrupt as err:
        logger.info("Successful exited")


if __name__ == "__main__":
    main()
