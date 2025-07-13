Experement's commands:

1. Depth meter's commands:

- Sense HAT:

`asvsoft sense-hat --period=40ms --loglevel=trace --transmitting-disabled | tee sh_data_exp.log`

`fgrep "moduleID:0X41" sh_data_exp.log | sed -E 's/(.*)ts:([0-9]+)(.*)Gx:(-?[0-9]+)\sGy:(-?[0-9]+)\sGz:(-?[0-9]+)\sAx:(-?[0-9]+)\sAy:(-?[0-9]+)\sAz:(-?[0-9]+)(.*)/\2,\4,\5,\6,\7,\8,\9/g' > sh_msr_exp.log`


- TOF Laser Range (B)

`asvsoft depthmeter --port /dev/ttyUSB0 --baudrate 921600 --loglevel=trace --transmitting-disabled | tee dm_data_exp.log`

`fgrep "moduleID:0X71" dm_data_exp.log | sed -E 's/(.*)ts:([0-9]+)(.*)Distance:([0-9]+)(.*)/\2,\4/g' > dm_dist_exp.log`

Registrator commands:

- Sonar:

`asvsoft depthmeter --port /dev/ttyAMA3 --baudrate 115200 --dst-port /dev/ttySC1 | tee depthmeter.rlog`

- LiDAR

`asvsoft lidar --port /dev/ttyUSB0 --baudrate 921600 | tee lidar.rlog`

- GNSS

`asvsoft neo-m8t --port /dev/ttyS0 --baudrate 9600 --dst-port /dev/ttyAMA5 --dst-baudrate 4800 | tee gnss.rlog`

- IMU:

`asvsoft sense-hat | fgrep -v i2c | tee sensethat.rlog`

- Camera:

`python3 registrator.py`

`asvsoft camera --dst-port /dev/ttySC1 | tee camera.rlog`
