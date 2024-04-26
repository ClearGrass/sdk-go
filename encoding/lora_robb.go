package encoding

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

func DecodeLoraRobbData(dataBytes []byte) (out *MessagePod, err error) {
	err = robbLoraCrcCheck(dataBytes)
	if err != nil {
		return
	}

	// 第一个字节不管，第二个字节是命令第三个字节是长度
	out = &MessagePod{}
	cmd := dataBytes[1]
	length := int(dataBytes[2])
	payload := dataBytes[3 : 3+length]
	sensorDataLen := 20

	switch cmd {
	case 0x41:
		if payload[0] == 1 { // 0历史数据 1 实时数据
			sensorData := parseRobbLoraSensorData(payload[5 : 5+sensorDataLen])
			sensorData.Timestamp = int64(binary.BigEndian.Uint32(payload[1:5]))
			out.Realtime = sensorData

			usbPlugin := int(payload[25])
			pmSn := strings.ToUpper(hex.EncodeToString(payload[26:30]))

			out.UsbPlugin = &usbPlugin
			out.PmSn = pmSn

			//co2Status := int(payload[30])
			//out.Co2Setting.AscOpen = co2Status
			out.FirmwareInfo = &FirmwareInfo{Version: string(payload[31:36])}
		}

		if payload[0] == 0 {
			timestamp := int64(binary.BigEndian.Uint32(payload[1:5]))
			interval := int64(binary.BigEndian.Uint16(payload[5:7]))

			for i, start := 0, 7; start < length; start += sensorDataLen {
				sensorData := parseRobbLoraSensorData(payload[start : start+sensorDataLen])
				sensorData.Timestamp = timestamp + int64(i)*interval
				out.History = append(out.History, sensorData)
			}
		}
	}

	return
}

func parseRobbLoraSensorData(data []byte) (out *SensorData) {
	out = &SensorData{}
	th := []byte{00, data[0], data[1], data[2]}
	value := binary.BigEndian.Uint32(th)
	temperature := (float64(value>>12) - 500.0) / 10.0 //高12bit为温度
	humidity := float64(value&0x000FFF) / 10.0         //低12位为湿度

	out.Temperature = &temperature
	out.Humidity = &humidity

	//pressure := float64(binary.LittleEndian.Uint16(data[3:5])) / 100.0
	//out.Pressure = &pressure

	co2 := float64(binary.BigEndian.Uint16(data[5:7]))
	out.Co2 = &co2

	pm25 := float64(binary.BigEndian.Uint16(data[7:9]))
	out.Pm25 = &pm25

	pm10 := float64(binary.BigEndian.Uint16(data[9:11]))
	out.Pm10 = &pm10

	tvoc := float64(binary.BigEndian.Uint16(data[11:13]))
	out.Tvoc = &tvoc

	noise := float64(binary.BigEndian.Uint16(data[13:15]))
	out.Noise = &noise

	lumen := float64(binary.BigEndian.Uint32(data[15:19]))
	out.Lumen = &lumen
	battery := int64(data[19]) // 255为直接供电 0-100是电池盒
	out.Battery = &battery
	return out
}

func robbLoraCrcCheck(data []byte) (err error) {
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
