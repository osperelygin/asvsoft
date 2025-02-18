Experement's commands:

1. Depth meter's commands:

- Sense HAT:

`asvsoft sense-hat --period=40ms --loglevel=trace --transmitting-disabled | tee sh_data_exp.log`

`fgrep "moduleID:0X41" sh_data_exp.log | sed -E 's/(.*)ts:([0-9]+)(.*)Gx:(-?[0-9]+)\sGy:(-?[0-9]+)\sGz:(-?[0-9]+)\sAx:(-?[0-9]+)\sAy:(-?[0-9]+)\sAz:(-?[0-9]+)(.*)/\2,\4,\5,\6,\7,\8,\9/g' > sh_msr_exp.log`


- TOF Laser Range (B)

`asvsoft depthmeter --port /dev/ttyUSB0 --baudrate 921600 --loglevel=trace --transmitting-disabled | tee dm_data_exp.log`

`fgrep "moduleID:0X71" dm_data_exp.log | sed -E 's/(.*)ts:([0-9]+)(.*)Distance:([0-9]+)(.*)/\2,\4/g' > dm_dist_exp.log`
