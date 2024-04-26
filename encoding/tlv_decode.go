package encoding

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func decodeInt16(in []byte) (out int16, err error) {
	if len(in) != 2 {
		err = errors.New("decodeInt16 input must be 2 bytes")
	}

	binary.Read(bytes.NewReader(in), binary.LittleEndian, &out)
	return
}

func decodeInt32(in []byte) (out int32, err error) {
	if len(in) != 4 {
		err = errors.New("decodeInt32 input must be 4 bytes")
	}
	binary.Read(bytes.NewReader(in), binary.LittleEndian, &out)
	return
}

func escapedData(src []byte, r map[byte]byte) {
	for i, c := range src {
		if r[c] > 0 {
			src[i] = r[c]
		}
	}
}

func buildEscapedMap(escaped []byte) (r map[byte]byte) {
	if len(escaped) == 6 {
		escaped = escaped[3:]
	}

	r = make(map[byte]byte)
	src := []byte{0x1A, 0x1B, 0x08}

	for i, c := range escaped {
		if c != 0x43 {
			r[c] = src[i]
		}
	}

	return
}

func escapedReplaceOne(src []byte) (thisPacket []byte, err error) {
	escapePrefixBytes := []uint8{0x27, 0x03, 0x00}
	index := bytes.Index(src, escapePrefixBytes)
	if index == -1 {
		return src, nil
	}

	// 0x27, 0x03, 0x00 后的三个字节为 替换字节
	// 在往后是文本内容 需要进行字节替换
	escaped := src[index+3 : index+6]
	escapedMap := buildEscapedMap(escaped)
	thisPacket = src[index+6:]
	escapedData(thisPacket, escapedMap)
	return thisPacket, nil
}

// 字节替换说明 字节替换+组合包 包的二进制数据如下
// 27 3 0 43 43 4 26 3 0 ....27 3 0 43 43 4 26 3 0...
// 27 3 0 表示3字节的字节替换
// 26 3 0 表示分包
// 27 3 0 在前，26 3 0在后，先根据 26 3 0 切割数组拆分包，然后获取前一段中的替换字节 进行替换

func escapedRepeatPacket(src []byte) (thisPacket []byte, err error) {
	var subPackets []*splitPackt

	escapePrefixBytes := []uint8{0x27, 0x03, 0x00} // 表示3字节的字节替换
	splitPrefixBytes := []byte{0x26, 0x03, 0x00}   // 根据0x26进行分包
	subData := bytes.Split(src, splitPrefixBytes)
	if len(subData) == 1 { // 没有分包 字节替换
		thisPacket, _ = escapedReplaceOne(src)
		return thisPacket, nil
	}

	// 查找前面的0x27 看是否有字节替换，并重新组合分包
	for i, sub := range subData {
		if i == 0 { // 第一段表头 不用解析
			continue
		}

		// 前一段数据
		preSub := subData[i-1]
		// 看前一段中是否有0x27,如果有则是本段数据中的0x27
		pack := &splitPackt{}
		index := bytes.Index(preSub, escapePrefixBytes)
		if index > -1 {
			pack.escaped = preSub[index:]
		}

		// 0x26 内容为包总数+包序号+当前包长度,第三个字节为长度
		size := int(sub[2])                             // 长度
		pack.content = sub[3 : 3+size]                  // 分隔的包中 不包含splitPrefixBytes
		pack.escapedMap = buildEscapedMap(pack.escaped) // 函数会自动过滤前三个字节（27 3 0）
		escapedData(pack.content, pack.escapedMap)
		subPackets = append(subPackets, pack)
	}

	for _, sub := range subPackets {
		thisPacket = append(thisPacket, sub.content...)
	}

	return
}

type splitPackt struct {
	escaped    []byte
	content    []byte
	escapedMap map[byte]byte
}

func escapePacket(packet []byte, pos int) ([]byte, int) {
	var repeatPrefixBytes = []byte{0x27, 0x27}
	if bytes.HasPrefix(packet, repeatPrefixBytes) { // 固件错误 多写了个0x27 在此适配
		packet = packet[1:]
	}

	packet, _ = escapedRepeatPacket(packet)
	return packet, len(packet)
}

func CrcCheckWithLength(data []byte, length int) (err error) {
	size := len(data)
	if length < 2 || size < length {
		err = errors.New("length not match")
		return
	}

	var sum uint16
	crc := binary.LittleEndian.Uint16(data[length-2 : length])
	for _, v := range data[:length-2] {
		sum += uint16(v)
	}

	if sum != crc {
		err = errors.New("check sum not match")
	}

	return
}

func TlvDecode(data []byte) (msg *MessagePod, err error) {
	// 拆包
	escapeData, pos := escapePacket(data, 0)
	_ = pos

	// 解密
	if escapeData[0] != 0x43 || escapeData[1] != 0x47 { // 需要aes解密
		//key, _ := hex.DecodeString(GetDefaultDeviceSecret(""))
		key, _ := hex.DecodeString("CF64060BDCF33F15A4E9166F7778CFE4")
		decryptedBytes, ierr := AESDecrypt(escapeData, key)
		if ierr != nil {
			err = ierr
			return
		}
		escapeData = decryptedBytes
	}

	// 解析
	msg, err = tlvDecode(escapeData)
	return
}

// 先将数据拆成tlv的三元组
func tlvDecode(data []byte) (msg *MessagePod, err error) {
	originData := data

	tlvData := TlvData{}
	sop := hex.EncodeToString(data[:2])
	data = data[2:]

	cmd := int(data[0])
	data = data[1:]

	length := int(binary.LittleEndian.Uint16(data[:2]))
	data = data[2:]

	tlvData.Sop = sop
	tlvData.Cmd = cmd
	tlvData.Length = length
	tlvData.PayloadByte = data[:length]

	data = data[:length]

	err = CrcCheckWithLength(originData, length+7)
	if err != nil {
		return
	}

	var subTlv []*TlvData

	// 解析子包
	for len(data) > 0 {
		key := int(data[0])
		length := int(binary.LittleEndian.Uint16(data[1:3]))
		valuePayload := data[3 : 3+length]

		subFiled := &TlvData{
			Cmd:         key,
			Length:      length,
			PayloadByte: valuePayload,
		}

		subTlv = append(subTlv, subFiled)
		data = data[3+length:]
	}

	msg = NewMessagePod(false)
	msg.CmdType = cmd

	for _, subFiled := range subTlv {
		//fmt.Printf("%x: %x \n", subFiled.Cmd, subFiled.PayloadByte)
		tlvDecodeData(subFiled, msg)
	}

	// pheasant co2 设备将pressure转化成 co2
	if IsPheasantCo2(msg.ProductId) {
		if msg.Realtime != nil {
			sensorDataPressureToCo2(msg.Realtime)
		}

		for _, d := range msg.History {
			sensorDataPressureToCo2(d)
		}

		// setting
		if co2Level, ok := msg.SensorDataLevel["co2"]; ok {
			if msg.Co2Setting == nil {
				msg.Co2Setting = &SensorMetricSetting{}
			}

			msg.Co2Setting.DataLevel.Min = co2Level.Value[0]
			msg.Co2Setting.DataLevel.Max = co2Level.Value[1]
		}
	}

	return
}

// 有co2的无气压，上报协议中co2复用pressure的位置
func sensorDataPressureToCo2(data *SensorData) {
	if data.Pressure == nil {
		return
	}

	v := *data.Pressure * 100
	data.Co2 = &v
	data.Pressure = nil
}

func tlvDecodeHistoryData(data *TlvData, unitLen int) (historyData []*SensorData) {
	valuePayload := data.PayloadByte
	timestamp := int64(binary.LittleEndian.Uint32(valuePayload[:4]))
	duration := int64(binary.LittleEndian.Uint16(valuePayload[4:6]))

	valuePayload = valuePayload[6:]
	for i := 0; i < (data.Length-6)/unitLen; i++ {
		dataPayload := valuePayload[:unitLen]
		valuePayload = valuePayload[unitLen:]

		sensorData := ParseSensorData(dataPayload, "history")
		sensorData.Timestamp = timestamp + duration*int64(i)
		sensorData.Time = time.Unix(sensorData.Timestamp, 0).Format("2006-01-02 15:04:05")
		historyData = append(historyData, sensorData)
	}

	return
}

func tlvDecodeRealtimeData(data *TlvData, withRssi bool) (out *SensorData) {
	out = &SensorData{}

	valuePayload := data.PayloadByte
	length := len(valuePayload)
	timestamp := int64(binary.LittleEndian.Uint32(valuePayload[:4]))

	switch length {
	case 24, 26:
		out = ParseSensorData(valuePayload[4:], "realtime")
		withRssi = false

	case 15, 19:
		out = ParseSensorData(valuePayload[4:], "realtime")
		withRssi = false

	default:
		out = ParseSensorData(valuePayload[4:10], "realtime")
	}

	out.Timestamp = timestamp
	out.Time = time.Unix(out.Timestamp, 0).Format("2006-01-02 15:04:05")

	if withRssi { // 11-12 信号
		rssi := int(int8(valuePayload[10]))
		out.Rssi = &rssi
	}

	// 13-14 frog 外接湿度
	if length == 14 && !(valuePayload[12] == 0xFF && valuePayload[13] == 0xFF) {
		h := float64(binary.LittleEndian.Uint16(valuePayload[12:14])&0xFFFF) / 10.0
		out.ProbHumidity = &h
	}

	return out
}

func tlvDecodeSensorData(metric string, data []byte, temperatureOffset bool) float32 {
	valInt := binary.LittleEndian.Uint16(data)
	valFloat := float32(valInt)

	switch metric {
	case "temperature", "probTemperature":
		if temperatureOffset {
			valFloat -= 500.0
		}
		valFloat = valFloat / 10.0

	case "humidity":
		valFloat = valFloat / 10.0

	case "pressure":
		valFloat = valFloat / 100.0

	}

	return valFloat
}

func tlvEncodeSensorData(metric string, num float32) []byte {
	buf := new(bytes.Buffer)

	val := tlvEncodeSensorDataToUint16(metric, num, true)
	_ = binary.Write(buf, binary.LittleEndian, val)
	return buf.Bytes()
}

// 解析 报警数据
func tlvDecodeAlert(data *TlvData) (setting *AlertSetting) {
	alertSetting := alertKeyToMetric[data.Cmd]

	var thresholdByte []byte
	valuePayload := data.PayloadByte

	switch data.Length {
	case 11:
		thresholdByte = valuePayload[9:11]

	case 13:
		thresholdByte = valuePayload[9:11]
		workTime := int(binary.LittleEndian.Uint16(valuePayload[11:13]))
		// 剩下两个字节为蜂鸣器时长
		alertSetting.WorkTime = workTime

	case 12:
		thresholdByte = valuePayload[10:12]

	//case 15: // frogs
	//	thresholdByte = valuePayload[13:15]

	case 24:
		thresholdByte = valuePayload[22:24]

	case 26:
		thresholdByte = valuePayload[24:26]

	default:
		return
	}

	alertSetting.Value = tlvDecodeSensorData(alertSetting.Metric, thresholdByte, true)
	setting = &alertSetting
	return
}

func tlvDecodeFrogsAlertEvent(data []byte, cmd int) (setting *AlertSetting, err error) {
	setting = &AlertSetting{}
	probSensor := data[0]

	switch data[1] {
	case 1:
		setting.Metric = "temperature"
		if probSensor == 1 {
			setting.Metric = "prob_temperature"
		}
	case 2:
		setting.Metric = "humidity"
		if probSensor == 1 {
			setting.Metric = "prob_humidity"
		}

	case 0x0B:
		if probSensor == 1 {
			setting.Metric = "co2_percent"
		}

	case 0x14:
		setting.Metric = "battery"

	default:
		return
	}

	setting.Operator = "GT"
	if cmd == 0x69 {
		setting.Operator = "LT"
	}

	value := int32(0)
	binary.Read(bytes.NewReader(data[2:6]), binary.LittleEndian, &value)
	setting.Value = float32(value) / 10.0

	return
}

func tlvDecodeFrogsAlert(data *TlvData) (setting *AlertSetting, err error) {
	if len(data.PayloadByte) != 18 {
		return
	}

	setting = &AlertSetting{}
	probSensor := data.PayloadByte[1]

	switch data.PayloadByte[2] {
	case 1:
		setting.Metric = "temperature"
		if probSensor == 1 {
			setting.Metric = "prob_temperature"
		}
	case 2:
		setting.Metric = "humidity"
		if probSensor == 1 {
			setting.Metric = "prob_humidity"
		}

	case 0x0B:
		if probSensor == 1 {
			setting.Metric = "co2_percent"
		}

	case 0x14:
		setting.Metric = "battery"

	default:
		return
	}

	setting.Operator = "GT"
	if data.Cmd == 0x69 {
		setting.Operator = "LT"
	}

	value := int32(0)
	binary.Read(bytes.NewReader(data.PayloadByte[12:16]), binary.LittleEndian, &value)
	setting.Value = float32(value) / 10.0
	return
}

func tlvDecodeSensorDataLevel(data *TlvData, metric string) (setting *DataLevelVal, err error) {
	valuePayload := data.PayloadByte

	val := make([]float32, 0, 2)
	val = append(val, tlvDecodeSensorData(metric, valuePayload[:2], false))
	val = append(val, tlvDecodeSensorData(metric, valuePayload[2:4], false))

	setting = &DataLevelVal{Value: val}
	return
}

func tlvDecodeData(data *TlvData, msg *MessagePod) {
	valuePayload := data.PayloadByte
	switch data.Cmd {
	case 0x01: //设备ID

	case 0x02: //设备SN

	case 0x03: // 历史数据
		unitLen := 6
		if IsRobb(msg.ProductId) {
			unitLen = 18

			if (len(data.PayloadByte)-6)%20 == 0 {
				unitLen = 20
			}
		}

		if IsFrogS(msg.ProductId) {
			unitLen = 9
			if (len(data.PayloadByte)-6)%13 == 0 {
				unitLen = 13
			}
		}

		msg.History = tlvDecodeHistoryData(data, unitLen)

	case 0x04: // 上报频率
		if msg.IntervalSetting == nil {
			msg.IntervalSetting = &IntervalSetting{}
		}

		reportInterval := int(binary.LittleEndian.Uint16(valuePayload)) * 60
		msg.IntervalSetting.ReportInterval = reportInterval

	case 0x05: // 采集频率
		if msg.IntervalSetting == nil {
			msg.IntervalSetting = &IntervalSetting{}
		}

		collectInterval := binary.LittleEndian.Uint16(valuePayload)
		msg.IntervalSetting.CollectInterval = int(collectInterval)

	case 0x06: // 蓝牙广播时间
		if msg.IntervalSetting == nil {
			msg.IntervalSetting = &IntervalSetting{}
		}

		bleInterval := binary.LittleEndian.Uint16(valuePayload)
		msg.IntervalSetting.BleInterval = int(bleInterval)

	case 0x1B: // 持续报警时间
		alertInterval := binary.LittleEndian.Uint16(valuePayload)
		_ = alertInterval
	case 0x11: // 固件版本号
		if msg.FirmwareInfo == nil {
			msg.FirmwareInfo = &FirmwareInfo{}
		}

		version := string(data.PayloadByte)
		msg.FirmwareInfo.Version = version

	case 0x12: // 固件url
		value := string(data.PayloadByte)
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-firmware"
		msg.Other[cmdStr] = value

	case 0x14: // 实时数据
		msg.Realtime = tlvDecodeRealtimeData(data, true)

	case 0x07, 0x08, 0x0A, 0x0B, 0x0D, 0x0E, 0x29, 0x2A, 0x39, 0x3A, 0x17: // 报警配置 或者报警事件
		alertSetting := tlvDecodeAlert(data)
		msg.AlertSetting = append(msg.AlertSetting, alertSetting)

		if msg.CmdType == 0x34 { // 57 为上报的配置
			msg.Realtime = tlvDecodeRealtimeData(data, false)
		}

	case 0x57, 0x58, 0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60: // 报警配置 或者报警事件
		alertSetting := tlvDecodeAlert(data)
		msg.AlertSetting = append(msg.AlertSetting, alertSetting)

		if msg.CmdType == 0x34 { // 57 为上报的配置
			msg.Realtime = tlvDecodeRealtimeData(data, false)
		}

	case 0x15: // 时间戳

	case 0x16: // SIM卡号

	case 0x19: // 设置温度单位
		value := hex.EncodeToString(data.PayloadByte)
		cmdStr := fmt.Sprintf("%#X", data.Cmd)
		cmdStr += "-temperatureUnit"
		msg.Other[cmdStr] = value

	case 0x1a: // 硬件版本号

	case 0x1d: // 断开标志位
		endFlag := 0
		for i, v := range valuePayload {
			endFlag += int(v) << (8 * i)
		}

		msg.EndFlag = &endFlag

	case 0x21:
		value := hex.EncodeToString(data.PayloadByte)
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-debug"
		msg.Other[cmdStr] = value

	case BINARY_KEY_MQTT: // mqtt 链接信息
		info := string(data.PayloadByte)
		msg.MqttSetting = &MqttSetting{Value: info}

	case 0x28: // 数据加密 秘钥
		value := hex.EncodeToString(data.PayloadByte)
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-cert"
		msg.Other[cmdStr] = value

	case 0x2c: // usb 插拔
		usbIn := 0
		for i, v := range valuePayload {
			usbIn += int(v) << (8 * i)
		}

		msg.UsbPlugin = &usbIn

	case BINARY_KEY_WIFI: // wifi 账号密码
		wifi := string(data.PayloadByte)
		wifi = strings.ReplaceAll(wifi, `"`, "")
		msg.WifiInfo = &WifiInfo{Desc: wifi}

	case 0x2F: // 温湿度 偏移量
		if msg.ReadingOffsetSetting == nil {
			msg.ReadingOffsetSetting = &ReadingOffsetSetting{}
		}

		tOffset, _ := decodeInt16(data.PayloadByte[:2])
		hOffset, _ := decodeInt16(data.PayloadByte[2:])
		msg.ReadingOffsetSetting.Temperature = &ReadingOffset{OffsetValue: float64(tOffset) / 10}
		msg.ReadingOffsetSetting.Humidity = &ReadingOffset{OffsetValue: float64(hOffset) / 10}

	case 0x30: // 外接温湿度 偏移量

	case 0x31: // 气压  偏移量

	case 0x32: // 外接温湿度

	case 0x33: // 历史数据
		unitLen := 8
		msg.History = tlvDecodeHistoryData(data, unitLen)

	case 0x34: // 模块版本号
		if msg.FirmwareInfo == nil {
			msg.FirmwareInfo = &FirmwareInfo{}
		}

		modelVersion := string(data.PayloadByte)
		msg.FirmwareInfo.ModuleVersion = modelVersion

	case 0x35: // mcu版本
		if msg.FirmwareInfo == nil {
			msg.FirmwareInfo = &FirmwareInfo{}
		}
		mcuVersion := string(data.PayloadByte)
		msg.FirmwareInfo.McuVersion = mcuVersion

	case 0x38: // 产品id
		deviceModel := binary.LittleEndian.Uint16(data.PayloadByte)
		msg.ProductId = int(deviceModel)

	case 0x3B: // co2 传感器采集间隔
		co2Setting := msg.Co2Setting
		if co2Setting == nil {
			co2Setting = &SensorMetricSetting{}
		}

		co2Interval := int(binary.LittleEndian.Uint16(valuePayload)) * 60
		co2Setting.SensorInterval.Collect = co2Interval
		msg.Co2Setting = co2Setting

	case 0x3D: // 关机时间
		shutdownTime := int(binary.LittleEndian.Uint16(valuePayload))

		if msg.BatterySetting == nil {
			msg.BatterySetting = &BatterySetting{}
		}

		msg.BatterySetting.DischargeShutdownTime.Value = shutdownTime * 60

	case 0x3F:
		if msg.ReadingOffsetSetting == nil {
			msg.ReadingOffsetSetting = &ReadingOffsetSetting{}
		}

		cOffset, _ := decodeInt16(data.PayloadByte)
		msg.ReadingOffsetSetting.Co2 = &ReadingOffset{OffsetPercent: float64(cOffset) / 10}

	case 0x40: // 自动校准开关
		co2Setting := msg.Co2Setting
		if co2Setting == nil {
			co2Setting = &SensorMetricSetting{}
		}

		ascOpen := byte(valuePayload[0])
		co2Setting.AscOpen.Value = ascOpen > 0
		msg.Co2Setting = co2Setting

	case 0x41: // 手动校准
		co2Setting := msg.Co2Setting
		if co2Setting == nil {
			co2Setting = &SensorMetricSetting{}
		}

		co2Setting.Reset.Value = true
		msg.Co2Setting = co2Setting

	case 0x42: // 临时上报数据
		value := int(binary.LittleEndian.Uint16(valuePayload))
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-tmpData"
		msg.Other[cmdStr] = strconv.Itoa(value)

	case 0x43: // sntp 服务
		value := string(data.PayloadByte)
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-sntpHost"
		msg.Other[cmdStr] = value

	case 0x44: // 设置SNTP校时开关
		value := hex.EncodeToString(data.PayloadByte)
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-sntpOpen"
		msg.Other[cmdStr] = value

	case 0x45: // CO2偏移校准
		value := int(binary.LittleEndian.Uint16(valuePayload))
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-co2Offset"
		msg.Other[cmdStr] = strconv.Itoa(value)

	case 0x46: // 温度偏移校准
		valueInt := int16(0)
		binary.Read(bytes.NewReader(valuePayload), binary.LittleEndian, &valueInt)
		value := float64(valueInt) / 10.0

		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-temperatureOffset"
		msg.Other[cmdStr] = fmt.Sprintf("%.1f", value)

	case 0x47: // 温度百分校准
		value := float32(binary.LittleEndian.Uint16(valuePayload)) / 10.0
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-temperatureOffset"
		msg.Other[cmdStr] = fmt.Sprintf("%.1f", value) + "%"

	case 0x48: // 湿度偏移校准
		value := float32(binary.LittleEndian.Uint16(valuePayload)) / 10.0
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-humidityOffset"
		msg.Other[cmdStr] = fmt.Sprintf("%.1f", value)

	case 0x49: // 湿度百分校准
		value := float32(binary.LittleEndian.Uint16(valuePayload)) / 10.0
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-humidityOffset"
		msg.Other[cmdStr] = fmt.Sprintf("%.1f", value) + "%"

	case 0x4A:
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		cmdStr += "-co2Status"
		msg.Other[cmdStr] = strconv.Itoa(int(data.PayloadByte[0]))

	case 0x3C, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56: // 读数分级标准
		metric := tlvGetMetricByCmd(data)

		levelSetting, _ := tlvDecodeSensorDataLevel(data, metric)
		msg.SensorDataLevel[metric] = levelSetting

	case 0x61:
		msg.PmSn = string(data.PayloadByte)

	case 0x63:
		metric := "light"
		val := float32(valuePayload[0])
		levelSetting := &DataLevelVal{Value: []float32{val}}
		msg.SensorDataLevel[metric] = levelSetting
	case 0x6A:
		needAck := int(valuePayload[0])
		msg.NeedAck = &needAck

	case 0x68, 0x69:
		if msg.CmdType == 0x34 { // frogs 报警数据
			alertSetting, _ := tlvDecodeFrogsAlertEvent(data.PayloadByte[17:], data.Cmd)
			msg.AlertSetting = append(msg.AlertSetting, alertSetting)

			msg.Realtime = parseFrogSensorData(data.PayloadByte[4:17])
			timestamp := int64(binary.LittleEndian.Uint32(valuePayload[:4]))
			msg.Realtime.Timestamp = timestamp
			msg.Realtime.Time = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		} else { // 配置
			alertSetting, _ := tlvDecodeFrogsAlert(data)
			msg.AlertSetting = append(msg.AlertSetting, alertSetting)
		}

	default:
		cmdStr := fmt.Sprintf("%#x", data.Cmd)
		msg.Other[cmdStr] = hex.EncodeToString(data.PayloadByte)
	}
}

func ParseSensorData(data []byte, dataType string) (out *SensorData) {
	out = &SensorData{}
	length := len(data)

	if length < 6 {
		return out
	}

	if length >= 18 {
		return parseRobbSensorData(data, dataType)
	}

	if length >= 9 && length <= 15 {
		return parseFrogSensorData(data)
	}

	th := []byte{data[0], data[1], data[2]}

	value := binary.LittleEndian.Uint32(append(th, 00))

	temperature := (float64(value>>12) - 500.0) / 10.0 //高12bit为温度
	humidity := float64(value&0x000FFF) / 10.0         //低12位为湿度

	out.Temperature = &temperature
	out.Humidity = &humidity

	flagByte := data[4]
	if (flagByte & 0xF0) == 0xF0 {
		probValue := int(data[4]&0x0F)*256 + int(data[3])

		if (probValue & 0x0FFF) != 0x0FFF {
			f := float64(probValue-500) / 10.0
			out.ProbTemperature = &f
		}
	} else {
		pressure := float64(binary.LittleEndian.Uint16(data[3:5])) / 100.0
		out.Pressure = &pressure
	}

	battery := int64(data[5])
	out.Battery = &battery

	// 这里是历史数据中的解析 外接湿度 与 事件中的不一样
	if length == 8 && !(data[6] == 0xFF && data[7] == 0xFF) {
		h := float64(binary.LittleEndian.Uint16(data[6:8])&0xFFFF) / 10.0
		out.ProbHumidity = &h
	}

	return out
}

func parseFrogSensorData(data []byte) (out *SensorData) {
	length := len(data)

	out = &SensorData{}
	th := []byte{data[0], data[1], data[2]}

	value := binary.LittleEndian.Uint32(append(th, 00))

	temperature := (float64(value>>12) - 500.0) / 10.0 //高12bit为温度
	humidity := float64(value&0x000FFF) / 10.0         //低12位为湿度

	out.Temperature = &temperature
	out.Humidity = &humidity

	probValue := int32(0)
	binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &probValue)

	switch data[3] {
	case 1:
		fValue := float64(probValue) / 10.0
		out.Co2Percent = &fValue
	case 2, 3:
		fValue := float64(probValue) / 10.0
		out.ProbTemperature = &fValue

	case 4, 6:
		fValue := float64(probValue) / 10.0
		out.Co2Percent = &fValue

		// 外置温湿度 新固件+4个字节
		if length == 13 || length == 15 {
			var t int16
			var h uint16
			binary.Read(bytes.NewReader(data[8:10]), binary.LittleEndian, &t)
			binary.Read(bytes.NewReader(data[10:12]), binary.LittleEndian, &h)

			ft := float64(t) / 10.0
			fh := float64(h) / 10.0
			out.ProbTemperature = &ft
			out.ProbHumidity = &fh
		} else { // 老固件 内置变外置
			out.ProbTemperature = out.Temperature
			out.ProbHumidity = out.Humidity
			out.Temperature = nil
			out.Humidity = nil
		}

	case 5:
		// 2字节温度+ 2字节湿度
		var t int16
		var h uint16
		binary.Read(bytes.NewReader(data[4:6]), binary.LittleEndian, &t)
		binary.Read(bytes.NewReader(data[6:8]), binary.LittleEndian, &h)
		ft := float64(t) / 10.0
		fh := float64(h) / 10.0
		out.ProbTemperature = &ft
		out.ProbHumidity = &fh
	}

	// 电量 新固件 多4个字节
	battery := int64(data[8])
	if length >= 13 {
		battery = int64(data[12])
	}
	out.Battery = &battery

	// 信号
	if len(data) == 11 {
		rssi := int(int8(data[9]))
		out.Rssi = &rssi
	}

	if len(data) == 15 {
		rssi := int(int8(data[13]))
		out.Rssi = &rssi
	}

	return out
}

func parseRobbSensorData(data []byte, dataType string) (out *SensorData) {
	out = &SensorData{}
	th := []byte{data[0], data[1], data[2]}
	value := binary.LittleEndian.Uint32(append(th, 00))
	temperature := (float64(value>>12) - 500.0) / 10.0 //高12bit为温度
	humidity := float64(value&0x000FFF) / 10.0         //低12位为湿度

	out.Temperature = &temperature
	out.Humidity = &humidity

	pressure := float64(binary.LittleEndian.Uint16(data[3:5])) / 100.0
	out.Pressure = &pressure

	co2 := float64(binary.LittleEndian.Uint16(data[5:7]))
	out.Co2 = &co2

	pm25 := float64(binary.LittleEndian.Uint16(data[7:9]))
	out.Pm25 = &pm25

	pm10 := float64(binary.LittleEndian.Uint16(data[9:11]))
	out.Pm10 = &pm10

	tvoc := float64(binary.LittleEndian.Uint16(data[11:13]))
	out.Tvoc = &tvoc

	noise := float64(binary.LittleEndian.Uint16(data[13:15]))
	out.Noise = &noise

	if dataType == "history" {
		if len(data) == 18 {
			lumen := float64(binary.LittleEndian.Uint16(data[15:17]))
			out.Lumen = &lumen
			battery := int64(data[17])
			out.Battery = &battery

		} else {
			lumen := float64(binary.LittleEndian.Uint32(data[15:19]))
			out.Lumen = &lumen
			battery := int64(data[19])
			out.Battery = &battery
		}
	}

	if dataType == "realtime" {
		if len(data) == 20 {
			lumen := float64(binary.LittleEndian.Uint16(data[15:17]))
			out.Lumen = &lumen
			battery := int64(data[17])
			out.Battery = &battery

			rssi := int(int8(data[18]))
			out.Rssi = &rssi

		} else {
			lumen := float64(binary.LittleEndian.Uint32(data[15:19]))
			out.Lumen = &lumen
			battery := int64(data[19])
			out.Battery = &battery
			rssi := int(int8(data[20]))
			out.Rssi = &rssi
		}
	}

	return out
}

func CrcCheck(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	var sum uint16
	crc := binary.LittleEndian.Uint16(data[len(data)-2:])
	for _, v := range data[:len(data)-2] {
		sum += uint16(v)
	}

	return sum == crc
}
