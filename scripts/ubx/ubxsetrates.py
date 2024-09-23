"""
ubxsetrates.py

This example illustrates how to send UBX commands to a receiver
(in this case a series of CFG-MSG commands)

Prerequisites:

python3 -m venv .venv && source .venv/bin/activate && pip3 install pyubx2 pyserial

Usage:

python3 ubxsetrates.py port="/dev/ttyS0" baudrate=9600 timout=0.1 rate=1 target-msgs="NAV-POSLLH,NAV-VELNED" enable-reading=false

"""

from sys import argv
from time import sleep
from serial import Serial
from pyubx2 import SET, UBX_MSGIDS, UBX_PROTOCOL, UBXMessage, UBXReader


def read_target_messages(serial):
    ubxreader = UBXReader(serial, protfilter=UBX_PROTOCOL)

    while True:
        try:
            if not serial.in_waiting:
                sleep(1)
                continue

            (_, parsed_data) = ubxreader.read()
            if parsed_data:
                print(parsed_data)

        except Exception as err:
            print(f"\n\nSomething went wrong {err}\n\n")


def main(**kwargs):
    port = kwargs.get("port", "/dev/ttyS0")
    baudrate = int(kwargs.get("baudrate", 9600))
    timeout = float(kwargs.get("timeout", 1))
    rate = int(kwargs.get("rate", 1))
    target_msgs = kwargs.get("target-msgs", "NAV-POSLLH,NAV-VELNED").split(",")
    enable_reading = bool(kwargs.get("enable-reading", False))

    with Serial(port, baudrate, timeout=timeout) as serial:
        # set the UART message rate for target UBX-NAV message via a CFG-MSG command
        print("\nSending CFG-MSG message rate configuration messages...\n")

        for msgid, msgname in UBX_MSGIDS.items():
            if not msgname in target_msgs:
                continue
            msg = UBXMessage(
                "CFG",
                "CFG-MSG",
                SET,
                msgClass=msgid[0],
                msgID=msgid[1],
                rateUART1=rate,
            )
            print(f"Setting message rate for {msgname} message type to {rate}...\n")

            serial.write(msg.serialize())

            sleep(1)

        if enable_reading:
            read_target_messages(serial)


if __name__ == "__main__":

    main(**dict(arg.split("=") for arg in argv[1:]))
