package proto

type LidarData struct {
	// Speed скорость вращения лидара в град/с
	Speed uint16
	// StartAngle начальный угол точек пакета в 0.01 град/с
	StartAngle uint16
	// Points массив точек измерения
	Points [12]Point
	// EndAngle конечный угол точек пакета в 0.01 град/с
	EndAngle uint16
	// Timestamp
	Timestamp uint16
}

type Point struct {
	Distance  uint16
	Intensity uint8
}
