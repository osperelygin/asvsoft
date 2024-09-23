package sensehat

const (
	QMI8658RegisterCtrl1 = 0x02 // QMI8658RegisterCtrl1 SPI Interface and Sensor Enable
	QMI8658RegisterCtrl2 = 0x03 // QMI8658RegisterCtrl2 Accelerometer control.
	QMI8658RegisterCtrl3 = 0x04 // QMI8658RegisterCtrl3 Gyroscope control.
	QMI8658RegisterCtrl5 = 0x06 // QMI8658RegisterCtrl5 Data processing settings.
	QMI8658RegisterCtrl7 = 0x08 // QMI8658RegisterCtrl7 Sensor enabled status.
)

const (
	I2CAddImuQMI8658 = byte(0x6B)
	I2CAddImuAK09918 = byte(0x0C)
)

const (
	QMI8658RegisterAxL = byte(0x35 + iota)
	QMI8658RegisterAyL
	QMI8658RegisterAzL
	QMI8658RegisterGxL
	QMI8658RegisterGyL
	QMI8658RegisterGzL
)

const (
	QMI8658_CTRL7_ACC_ENABLE = 0x01
	QMI8658_CTRL7_GYR_ENABLE = 0x02
)

const (
	QMI8658AccRange2g  = iota << 4 // QMI8658AccRange2g +/- 2g range
	QMI8658AccRange4g              // QMI8658AccRange4g +/- 4g range
	QMI8658AccRange8g              // QMI8658AccRange8g +/- 8g range
	QMI8658AccRange16g             // QMI8658AccRange16g +/- 16g range
)

const (
	QMI8658GyrRange16dps   = iota << 4 // QMI8658GyrRange16dps +-16 degrees per second.
	QMI8658GyrRange32dps               // QMI8658GyrRange32dps +-32 degrees per second.
	QMI8658GyrRange64dps               // QMI8658GyrRange64dps +-64 degrees per second.
	QMI8658GyrRange128dps              // QMI8658GyrRange128dps +-128 degrees per second.
	QMI8658GyrRange256dps              // QMI8658GyrRange256dps +-256 degrees per second.
	QMI8658GyrRange512dps              // QMI8658GyrRange512dps +-512 degrees per second.
	QMI8658GyrRange1024dps             // QMI8658GyrRange1024dps +-1024 degrees per second.
	QMI8658GyrRange2048dps             // QMI8658GyrRange2048dps +-2048 degrees per second.
)

const (
	QMI8658GyrOdr8000Hz  = iota // QMI8658GyrOdr8000Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr4000Hz         // QMI8658GyrOdr4000Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr2000Hz         // QMI8658GyrOdr2000Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr1000Hz         // QMI8658GyrOdr1000Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr500Hz          // QMI8658GyrOdr500Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr250Hz          // QMI8658GyrOdr250Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr125Hz          // QMI8658GyrOdr125Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr62_5Hz         // QMI8658GyrOdr62_5Hz High resolution 8000Hz output rate.
	QMI8658GyrOdr31_25Hz        // QMI8658GyrOdr31_25Hz High resolution 8000Hz output rate.
)

const (
	QMI8658AccOdr_8000Hz  = iota // QMI8658AccOdr_8000Hz High resolution 8000Hz output rate.
	QMI8658AccOdr_4000Hz  = 0x01 // QMI8658AccOdr_4000Hz High resolution 4000Hz output rate.
	QMI8658AccOdr_2000Hz  = 0x02 // QMI8658AccOdr_2000Hz High resolution 2000Hz output rate.
	QMI8658AccOdr_1000Hz  = 0x03 // QMI8658AccOdr_1000Hz High resolution 1000Hz output rate.
	QMI8658AccOdr_500Hz   = 0x04 // QMI8658AccOdr_500Hz High resolution 500Hz output rate.
	QMI8658AccOdr_250Hz   = 0x05 // QMI8658AccOdr_250Hz High resolution 250Hz output rate.
	QMI8658AccOdr_125Hz   = 0x06 // QMI8658AccOdr_125Hz High resolution 125Hz output rate.
	QMI8658AccOdr_62_5Hz  = 0x07 // QMI8658AccOdr_62_5Hz High resolution 62.5Hz output rate.
	QMI8658AccOdr_31_25Hz = 0x08 // QMI8658AccOdr_31_25Hz High resolution 31.25Hz output rate.
)

const (
	QMI8658AccOdr_LowPower_128Hz = (iota + 0x0C) // QMI8658AccOdr_LowPower_128Hz Low power 128Hz output rate.
	QMI8658AccOdr_LowPower_21Hz                  // QMI8658AccOdr_LowPower_21Hz Low power 21Hz output rate.
	QMI8658AccOdr_LowPower_11Hz                  // QMI8658AccOdr_LowPower_11Hz Low power 11Hz output rate.
	QMI8658AccOdr_LowPower_3Hz                   // QMI8658AccOdr_LowPower_3Hz Low power 3Hz output rate.
)

const (
	AK09918_WIA2 = 0x01 // AK09918_WIA2 Device ID
)

const (
	AK09918_CONTINUOUS_10HZ  = 0x02
	AK09918_CONTINUOUS_20HZ  = 0x04
	AK09918_CONTINUOUS_50HZ  = 0x06
	AK09918_CONTINUOUS_100HZ = 0x08
)

const (
	AK09918_CNTL2 = 0x31 // AK09918_CNTL2 Control settings
	AK09918_CNTL3 = 0x32 // AK09918_CNTL3 Control settings
)

const (
	AK09918_SRST_BIT = 0x01 // AK09918_SRST_BIT Soft Reset
	AK09918_ST1      = 0x10 // AK09918_ST1 DataStatus 1
)

const (
	AK09918_HXL = 0x11
	AK09918_HXH = 0x12
	AK09918_HYL = 0x13
	AK09918_HYH = 0x14
	AK09918_HZL = 0x15
	AK09918_HZH = 0x16
)
