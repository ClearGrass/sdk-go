package encoding

import (
	"encoding/binary"
	"errors"
)

func DecodeLoraPheasantCo2Data(dataBytes []byte) (out *MessagePod, err error) {
	err = pheasantCo2LoraCrcCheck(dataBytes)
	if err != nil {
		return
	}

	// 第一个字节不管，第二个字节是命令第三个字节是长度
	out = &MessagePod{}
	cmd := dataBytes[1]
	length := int(dataBytes[2])
	payload := dataBytes[3 : 3+length]
	sensorDataLen := 6

	switch cmd {
	case 0x41:
		if payload[0] == 1 { // 0历史数据 1 实时数据
			sensorData := parsePheasantCo2LoraSensorData(payload[5 : 5+sensorDataLen])
			sensorData.Timestamp = int64(binary.BigEndian.Uint32(payload[1:5]))
			out.Realtime = sensorData
		}

		if payload[0] == 0 {
			timestamp := int64(binary.BigEndian.Uint32(payload[1:5]))
			interval := int64(binary.BigEndian.Uint16(payload[5:7]))

			for i, start := 0, 7; start < length; start += sensorDataLen {
				sensorData := parsePheasantCo2LoraSensorData(payload[start : start+sensorDataLen])
				sensorData.Timestamp = timestamp + int64(i)*interval
				out.History = append(out.History, sensorData)
			}
		}
	}

	return
}

func parsePheasantCo2LoraSensorData(data []byte) (out *SensorData) {
	out = &SensorData{}
	th := []byte{00, data[0], data[1], data[2]}
	value := binary.BigEndian.Uint32(th)
	temperature := (float64(value>>12) - 500.0) / 10.0 //高12bit为温度
	humidity := float64(value&0x000FFF) / 10.0         //低12位为湿度

	out.Temperature = &temperature
	out.Humidity = &humidity

	co2 := float64(binary.BigEndian.Uint16(data[3:5]))
	out.Co2 = &co2

	battery := int64(data[5]) // 255为直接供电 0-100是电池盒
	out.Battery = &battery

	return out
}

func pheasantCo2LoraCrcCheck(data []byte) (err error) {
	size := len(data)

	var crcReg, check = uint(0xffff), uint(0)
	for i := 0; i < size; i++ {
		crcReg = crcReg ^ uint(data[i])
		for j := 0; j < 8; j++ {
			check = crcReg & 0x0001
			crcReg >>= 1
			if check == 0x0001 {
				crcReg ^= 0xA001
			}
		}
	}

	crc := crcReg>>8 | crcReg<<8
	if crc != 0 {
		err = errors.New("crc check failed")
	}
	return
}
