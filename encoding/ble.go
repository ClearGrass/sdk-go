package encoding

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func DecodeBleData(hexData string) (sensorData *SensorData, err error) {
	sensorData = &SensorData{}
	bs, _ := hex.DecodeString(hexData)

	packlen := len(bs)
	index := 0

	var dataPack []byte
	for index < packlen-1 {
		lenth := int(bs[index])

		subPack := bs[index+1 : index+1+lenth]
		if subPack[0] == 0x16 {
			dataPack = subPack
			break
		}
		index += index + lenth + 1
	}

	if len(dataPack) == 0 {
		return nil, errors.New("not vaild data")
	}

	uuid := hex.EncodeToString(dataPack[1:3])

	if uuid != "cdfd" && uuid != "f9ff" {
		return nil, errors.New("not vaild data")
	}

	productId := dataPack[4]
	_ = productId

	macByte := make([]byte, 6)
	copy(macByte, dataPack[5:11])
	reverseBytes(macByte)

	mac := hex.EncodeToString(macByte)
	sensorData.Mac = strings.ToUpper(mac)

	sensorDataPack := dataPack[11:]
	sIndex := 0
	for sIndex < len(sensorDataPack) {
		sKey := sensorDataPack[sIndex]
		sLen := int(sensorDataPack[sIndex+1])

		switch sKey {
		case 0x01:
			if sLen == 4 {
				temperatureInt16 := int16(binary.LittleEndian.Uint16(sensorDataPack[sIndex+2 : sIndex+4]))
				humidityInt16 := int16(binary.LittleEndian.Uint16(sensorDataPack[sIndex+4 : sIndex+6]))

				temperature := float64(temperatureInt16) / 10.0
				humidity := float64(humidityInt16) / 10.0

				sensorData.Temperature = &temperature
				sensorData.Humidity = &humidity
			}
		case 0x02:
			if sLen == 1 {
				battery := int64(sensorDataPack[sIndex+2])
				sensorData.Battery = &battery
			}
		}

		sIndex += 2 + sLen
	}

	return
}
