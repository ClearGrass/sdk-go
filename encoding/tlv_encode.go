package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
)

func encodeUint8(in interface{}) (out []byte, err error) {
	num, ok := in.(uint8)
	if !ok {
		err = errors.New("value must be uint8")
		return
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.LittleEndian, num)
	out = buf.Bytes()

	return
}

func encodeUint16(in interface{}) (out []byte, err error) {
	num, ok := in.(uint16)
	if !ok {
		err = errors.New("value must be uint16")
		return
	}
	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.LittleEndian, num)
	out = buf.Bytes()
	return
}

func encodeUint16Slice(in interface{}) (out []byte, err error) {
	nums, ok := in.([]uint16)
	if !ok {
		err = errors.New("value must be uint16 slice")
		return
	}

	out = make([]byte, 0, len(nums)*2)
	buf := &bytes.Buffer{}
	for _, num := range nums {
		err = binary.Write(buf, binary.LittleEndian, num)
		out = append(out, buf.Bytes()...)
		buf.Reset()
	}
	return
}

func encodeInt16(in interface{}) (out []byte, err error) {
	num, ok := in.(int16)
	if !ok {
		err = errors.New("value must be int16")
		return
	}
	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.LittleEndian, num)
	out = buf.Bytes()
	return
}

func encodeInt16Slice(in interface{}) (out []byte, err error) {
	nums, ok := in.([]int16)
	if !ok {
		err = errors.New("value must be int16 slice")
		return
	}

	out = make([]byte, 0, len(nums)*2)
	buf := &bytes.Buffer{}
	for _, num := range nums {
		err = binary.Write(buf, binary.LittleEndian, num)
		out = append(out, buf.Bytes()...)
		buf.Reset()
	}
	return
}

func encodeString(in interface{}) (out []byte, err error) {
	str, ok := in.(string)
	if !ok {
		err = errors.New("value must be string")
	}

	out = []byte(str)
	return
}

func encodeBytes(in interface{}) (out []byte, err error) {
	out, ok := in.([]byte)
	if !ok {
		err = errors.New("value must be bytes")
	}

	return out, nil
}

type encodeMethod func(interface{}) ([]byte, error)

const (
	BINARY_CMD_MQTT = byte(0x3A)
)

const (
	BINARY_KEY_REPORT_INTERVAL = 0x04
	BINARY_KEY_DATA_INTERVAL   = 0x05
	BINARY_KEY_BLE_INTERVAL    = 0x06
	BINARY_KEY_END_FLAG        = 0x1D
	BINARY_KEY_DEBUG           = 0x21
	BINARY_KEY_MQTT            = 0x25
	BINARY_KEY_SECRET_KEY      = 0x28
	BINARY_KEY_BLE_NAME        = 0x36

	BINARY_KEY_ALERT_TEMPERATURE_GT      = 0x07
	BINARY_KEY_ALERT_TEMPERATURE_LT      = 0x08
	BINARY_KEY_ALERT_HUMIDITY_GT         = 0x0A
	BINARY_KEY_ALERT_HUMIDITY_LT         = 0x0B
	BINARY_KEY_ALERT_PRESSURE_GT         = 0x0D
	BINARY_KEY_ALERT_PRESSURE_LT         = 0x0E
	BINARY_KEY_ALERT_PROB_TEMPERATURE_GT = 0x29
	BINARY_KEY_ALERT_PROB_TEMPERATURE_LT = 0x2A
	BINARY_KEY_ALERT_CO2_GT              = 0x39
	BINARY_KEY_ALERT_CO2_LT              = 0x3A
	BINARY_KEY_ALERT_Tvoc_GT             = 0x57
	BINARY_KEY_ALERT_Tvoc_LT             = 0x58
	BINARY_KEY_ALERT_Pm25_GT             = 0x59
	BINARY_KEY_ALERT_Pm25_LT             = 0x5A
	BINARY_KEY_ALERT_Pm10_GT             = 0x5B
	BINARY_KEY_ALERT_Pm10_LT             = 0x5C
	BINARY_KEY_ALERT_Noise_GT            = 0x5D
	BINARY_KEY_ALERT_Noise_LT            = 0x5E
	BINARY_KEY_ALERT_Lumen_GT            = 0x5F
	BINARY_KEY_ALERT_Lumen_LT            = 0x60
	BINARY_KEY_ALERT_Battery_LT          = 0x17

	BINARY_KEY_UNIT_TEMPERATURE = 0x19
	BINARY_KEY_UNIT_TVOC_INDEX  = 0x62
	BINARY_KEY_WIFI             = 0x20

	BINARY_KEY_TEMPERATURE_OFFSET = 0x2F
	BINARY_KEY_CO2_OFFSET         = 0x3F
)

var encodeMethodMap = map[int]encodeMethod{
	BINARY_KEY_REPORT_INTERVAL:           encodeUint16,
	BINARY_KEY_DATA_INTERVAL:             encodeUint16,
	BINARY_KEY_BLE_INTERVAL:              encodeUint16,
	BINARY_KEY_END_FLAG:                  encodeUint8,
	BINARY_KEY_DEBUG:                     encodeUint8,
	BINARY_KEY_MQTT:                      encodeString,
	BINARY_KEY_SECRET_KEY:                encodeString,
	BINARY_KEY_BLE_NAME:                  encodeString,
	BINARY_KEY_ALERT_TEMPERATURE_GT:      encodeBytes,
	BINARY_KEY_ALERT_TEMPERATURE_LT:      encodeBytes,
	BINARY_KEY_ALERT_HUMIDITY_GT:         encodeBytes,
	BINARY_KEY_ALERT_HUMIDITY_LT:         encodeBytes,
	BINARY_KEY_ALERT_PRESSURE_GT:         encodeBytes,
	BINARY_KEY_ALERT_PRESSURE_LT:         encodeBytes,
	BINARY_KEY_ALERT_PROB_TEMPERATURE_GT: encodeBytes,
	BINARY_KEY_ALERT_PROB_TEMPERATURE_LT: encodeBytes,
	BINARY_KEY_ALERT_CO2_GT:              encodeBytes,
	BINARY_KEY_ALERT_CO2_LT:              encodeBytes,
	BINARY_KEY_ALERT_Tvoc_GT:             encodeBytes,
	BINARY_KEY_ALERT_Tvoc_LT:             encodeBytes,
	BINARY_KEY_ALERT_Pm25_GT:             encodeBytes,
	BINARY_KEY_ALERT_Pm25_LT:             encodeBytes,
	BINARY_KEY_ALERT_Pm10_GT:             encodeBytes,
	BINARY_KEY_ALERT_Pm10_LT:             encodeBytes,
	BINARY_KEY_ALERT_Noise_GT:            encodeBytes,
	BINARY_KEY_ALERT_Noise_LT:            encodeBytes,
	BINARY_KEY_ALERT_Lumen_GT:            encodeBytes,
	BINARY_KEY_ALERT_Lumen_LT:            encodeBytes,
	BINARY_KEY_ALERT_Battery_LT:          encodeBytes,

	BINARY_KEY_UNIT_TEMPERATURE: encodeUint8,
	BINARY_KEY_UNIT_TVOC_INDEX:  encodeUint8,

	BINARY_KEY_TEMPERATURE_OFFSET: encodeInt16Slice,
	BINARY_KEY_CO2_OFFSET:         encodeInt16,

	BINARY_KEY_WIFI: encodeString,
	0x38:            encodeUint16,

	0x3B: encodeUint16,
	//0x3C: encodeUint16Slice,
	0x3D: encodeUint16,
	0x40: encodeUint8,
	0x41: encodeUint8,
	0x42: encodeUint16,

	0x43: encodeString,
	0x44: encodeUint8,

	0x3C: encodeUint16Slice,
	0x4F: encodeUint16Slice,
	0x50: encodeUint16Slice,
	0x51: encodeUint16Slice,
	0x52: encodeUint16Slice,
	0x53: encodeUint16Slice,
	0x54: encodeUint16Slice,
	0x55: encodeUint16Slice,
	0x56: encodeUint16Slice,
	0x63: encodeUint8,
	0x6A: encodeUint8,
	0x68: encodeBytes,
	0x69: encodeBytes,
}

func encodeKv(k int, v interface{}) (out []byte, err error) {
	encodeFunc := encodeMethodMap[k]
	if encodeFunc == nil {
		err = fmt.Errorf("filed: %d has no encode method", k)
		return
	}

	payload, err := encodeFunc(v)
	if err != nil {
		err = fmt.Errorf("filed: %d, %s", k, err.Error())
		return
	}

	size := len(payload)
	sizeByte, _ := encodeUint16(uint16(size))

	out = make([]byte, 0, 1+size+len(sizeByte))

	key := byte(k)
	out = append(out, key)
	out = append(out, sizeByte...)
	out = append(out, payload...)
	return
}

func addCrc(in []byte) (out []byte) {
	sum := 0
	for _, v := range in {
		sum += int(v)
	}

	crc, _ := encodeUint16(uint16(sum))
	out = append(in, crc...)
	return out
}

func sortMapKey(in map[int]interface{}) (keys []int) {
	keys = make([]int, 0, len(in))
	for k, _ := range in {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	return
}

func convertTimeToMinutes(time string) []byte {
	return []byte{0, 0, 0, 0}
}

var (
	alertMetricToKey = map[string]int{
		"temperature-GT":     BINARY_KEY_ALERT_TEMPERATURE_GT,
		"temperature-LT":     BINARY_KEY_ALERT_TEMPERATURE_LT,
		"probTemperature-GT": BINARY_KEY_ALERT_PROB_TEMPERATURE_GT,
		"probTemperature-LT": BINARY_KEY_ALERT_PROB_TEMPERATURE_LT,
		"humidity-GT":        BINARY_KEY_ALERT_HUMIDITY_GT,
		"humidity-LT":        BINARY_KEY_ALERT_HUMIDITY_LT,
		"pressure-GT":        BINARY_KEY_ALERT_PRESSURE_GT,
		"pressure-LT":        BINARY_KEY_ALERT_PRESSURE_LT,
		"co2-GT":             BINARY_KEY_ALERT_CO2_GT,
		"co2-LT":             BINARY_KEY_ALERT_CO2_LT,
		"tvoc-GT":            0x57,
		"tvoc-LT":            0x58,
		"pm25-GT":            0x59,
		"pm25-LT":            0x5A,
		"pm10-GT":            0x5B,
		"pm10-LT":            0x5C,
		"noise-GT":           0x5D,
		"noise-LT":           0x5E,
		"lumen-GT":           0x5F,
		"lumen-LT":           0x60,
		"battery-LT":         BINARY_KEY_ALERT_Battery_LT,
	}

	alertKeyToMetric = map[int]AlertSetting{
		BINARY_KEY_ALERT_TEMPERATURE_GT: AlertSetting{Metric: "temperature", Operator: "GT"},
		BINARY_KEY_ALERT_TEMPERATURE_LT: AlertSetting{Metric: "temperature", Operator: "LT"},
		0x29:                            AlertSetting{Metric: "probTemperature", Operator: "GT"},
		0x2A:                            AlertSetting{Metric: "probTemperature", Operator: "LT"},
		BINARY_KEY_ALERT_HUMIDITY_GT:    AlertSetting{Metric: "humidity", Operator: "GT"},
		BINARY_KEY_ALERT_HUMIDITY_LT:    AlertSetting{Metric: "humidity", Operator: "LT"},
		BINARY_KEY_ALERT_PRESSURE_GT:    AlertSetting{Metric: "pressure", Operator: "GT"},
		BINARY_KEY_ALERT_PRESSURE_LT:    AlertSetting{Metric: "pressure", Operator: "LT"},
		BINARY_KEY_ALERT_CO2_GT:         AlertSetting{Metric: "co2", Operator: "GT"},
		BINARY_KEY_ALERT_CO2_LT:         AlertSetting{Metric: "co2", Operator: "LT"},
		0x57:                            AlertSetting{Metric: "tvoc", Operator: "GT"},
		0x58:                            AlertSetting{Metric: "tvoc", Operator: "LT"},
		0x59:                            AlertSetting{Metric: "pm25", Operator: "GT"},
		0x5A:                            AlertSetting{Metric: "pm25", Operator: "LT"},
		0x5B:                            AlertSetting{Metric: "pm10", Operator: "GT"},
		0x5C:                            AlertSetting{Metric: "pm10", Operator: "LT"},
		0x5D:                            AlertSetting{Metric: "noise", Operator: "GT"},
		0x5E:                            AlertSetting{Metric: "noise", Operator: "LT"},
		0x5F:                            AlertSetting{Metric: "lumen", Operator: "GT"},
		0x60:                            AlertSetting{Metric: "lumen", Operator: "LT"},
		BINARY_KEY_ALERT_Battery_LT:     AlertSetting{Metric: "battery", Operator: "LT"},
	}

	alertKeyToMetricStr = map[int]string{
		BINARY_KEY_ALERT_TEMPERATURE_GT: "temperature-GT",
		BINARY_KEY_ALERT_TEMPERATURE_LT: "temperature-LT",
		0x29:                            "probTemperature-GT",
		0x2A:                            "probTemperature-LT",
		BINARY_KEY_ALERT_HUMIDITY_GT:    "humidity-GT",
		BINARY_KEY_ALERT_HUMIDITY_LT:    "humidity-LT",
		BINARY_KEY_ALERT_PRESSURE_GT:    "pressure-GT",
		BINARY_KEY_ALERT_PRESSURE_LT:    "pressure-LT",
		BINARY_KEY_ALERT_CO2_GT:         "co2-GT",
		BINARY_KEY_ALERT_CO2_LT:         "co2-LT",
		0x57:                            "tvoc-GT",
		0x58:                            "tvoc-LT",
		0x59:                            "pm25-GT",
		0x5A:                            "pm25-LT",
		0x5B:                            "pm10-GT",
		0x5C:                            "pm10-LT",
		0x5D:                            "noise-GT",
		0x5E:                            "noise-LT",
		0x5F:                            "lumen-GT",
		0x60:                            "lumen-LT",
		BINARY_KEY_ALERT_Battery_LT:     "battery-LT",
	}
)

func convertFrogsAlert(alert *AlertSetting) (subPacket *TlvData, err error) {
	subPacket = &TlvData{Cmd: 0x68}
	if alert.Operator == "LT" {
		subPacket.Cmd = 0x69
	}

	// 温度 0x01
	// 湿度 0x02
	// 外接co2 0x01
	// 外接温度 0x02
	// 外接湿度 0x03

	probSensor := byte(0)

	sensorType := byte(0)
	switch alert.Metric {
	case "battery":
		sensorType = 0x14

	case "temperature":
		sensorType = 1

	case "humidity":
		sensorType = 2

	case "co2_percent":
		probSensor = 1
		sensorType = 0x0B
	case "prob_temperature":
		probSensor = 1
		sensorType = 1
	case "prob_humidity":
		probSensor = 1
		sensorType = 2
	}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, int32(alert.Value*10))

	num := byte(1)
	repeat := byte(1)
	value := make([]byte, 0, 18)
	value = append(value, num)        //序号
	value = append(value, probSensor) // 0 内置传感器 1 外置
	value = append(value, sensorType) // 传感器类型
	value = append(value, repeat)     // 重复次数

	value = append(value, convertTimeToMinutes("")...) // 开始时间
	value = append(value, convertTimeToMinutes("")...) // 结束时间

	value = append(value, buf.Bytes()...) // 警保值4字节
	value = append(value, 0, 0)           // 警保值2字节
	subPacket.PayloadAny = value

	return subPacket, nil
}

func convertAlert(alert *AlertSetting, hasBuzzer bool) (subPacket *TlvData, err error) {
	str := alert.Metric + "-" + alert.Operator
	key := alertMetricToKey[str]

	value := make([]byte, 0)
	value = append(value, 1) // 重复次数

	value = append(value, convertTimeToMinutes("")...) // 开始时间
	value = append(value, convertTimeToMinutes("")...) // 结束时间

	buf := new(bytes.Buffer)
	switch alert.Metric {

	case "temperature", "probTemperature":
		_ = binary.Write(buf, binary.LittleEndian, uint16(alert.Value*10+500)) //放大10倍正向偏移500
		value = append(value, buf.Bytes()...)

	case "humidity":
		_ = binary.Write(buf, binary.LittleEndian, uint16(alert.Value*10))
		value = append(value, buf.Bytes()...)

	case "pressure":
		_ = binary.Write(buf, binary.LittleEndian, uint16(alert.Value*100))
		value = append(value, buf.Bytes()...)

	default:
		_ = binary.Write(buf, binary.LittleEndian, uint16(alert.Value))
		value = append(value, buf.Bytes()...)
	}

	if hasBuzzer {
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, uint16(alert.WorkTime))
		value = append(value, buf.Bytes()...)
	}

	subPacket = &TlvData{
		Cmd:        key,
		PayloadAny: value,
	}
	return
}

func convertCo2SensorSetting(metric string, setting *SensorMetricSetting) (packets []*TlvData, err error) {
	// 采集间隔
	{
		sensorCollect := uint16(setting.SensorInterval.Collect / 60)
		subPacket := &TlvData{
			Cmd:        0x3B,
			PayloadAny: sensorCollect,
		}

		packets = append(packets, subPacket)
	}

	// 读数分级
	{
		dataLevel := []uint16{uint16(setting.DataLevel.Min), uint16(setting.DataLevel.Max)}
		subPacket := &TlvData{
			Cmd:        0x3C,
			PayloadAny: dataLevel,
		}

		packets = append(packets, subPacket)
	}

	// 自动校准
	if setting.AscOpen.Value {
		ascOpen := byte(1)
		subPacket := &TlvData{
			Cmd:        0x40,
			PayloadAny: ascOpen,
		}

		packets = append(packets, subPacket)
	}

	// 充值传感器
	if setting.Reset.Value {
		subPacket := &TlvData{
			Cmd:        0x41,
			PayloadAny: byte(1),
		}

		packets = append(packets, subPacket)
	}

	return
}

func tlvEncodeSensorDataToUint16(metric string, num float32, temperatureOffset bool) (val uint16) {
	switch metric {
	case "temperature", "probTemperature":
		val = uint16(num * 10)

		if temperatureOffset {
			val += 500
		}

	case "humidity":
		val = uint16(num * 10)

	case "pressure":
		val = uint16(num * 100)

	default:
		val = uint16(num)
	}

	return val
}

func convertSensorSetting(metric string, setting *SensorMetricSetting) (packets []*TlvData, err error) {
	switch metric {
	case "co2":
		return convertCo2SensorSetting(metric, setting)
	}

	// 只处理读数分级
	l1 := tlvEncodeSensorDataToUint16(metric, setting.DataLevel.Min, false)
	l2 := tlvEncodeSensorDataToUint16(metric, setting.DataLevel.Max, false)

	dataLevel := []uint16{l1, l2}
	subPacket := &TlvData{
		Cmd:        tlvGetSensorDataLevelCmdByMetric(metric),
		PayloadAny: dataLevel,
	}

	packets = append(packets, subPacket)

	return
}

func convertDataLevelSetting(metric string, setting *DataLevelVal) (subPacket *TlvData, err error) {
	if metric == "light" {
		lightVal := byte(0)

		if len(setting.Value) > 0 && setting.Value[0] > 0 {
			lightVal = 1
		}

		subPacket = &TlvData{
			Cmd:        0x63,
			PayloadAny: lightVal,
		}

		return
	}

	dataLevel := make([]uint16, 0, len(setting.Value))
	for _, val := range setting.Value {
		dataLevel = append(dataLevel, tlvEncodeSensorDataToUint16(metric, val, false))
	}

	subPacket = &TlvData{
		Cmd:        tlvGetSensorDataLevelCmdByMetric(metric),
		PayloadAny: dataLevel,
	}

	return subPacket, nil
}

func convertBatterySetting(setting *BatterySetting) (subPacket *TlvData, err error) {
	shutdownTime := uint16(setting.DischargeShutdownTime.Value / 60)
	subPacket = &TlvData{
		Cmd:        0x3D,
		PayloadAny: shutdownTime,
	}
	return
}

func convertMqttSetting(setting *MqttSetting) (subPacket *TlvData, err error) {
	mqttStr := fmt.Sprintf("%s %s %s %s %s %s %s",
		setting.Host, setting.Port, setting.User, setting.Password, setting.ClientId, setting.DownTopic, setting.UpTopic)

	subPacket = &TlvData{
		Cmd:        BINARY_KEY_MQTT,
		PayloadAny: mqttStr,
	}
	return
}

func convertUnitSetting(setting *UnitSetting) (subPackets []*TlvData, err error) {
	if setting.Temperature != nil {
		unitValue := byte(0)

		if setting.Temperature.Value == "F" {
			unitValue = 1
		}

		subPacket := &TlvData{
			Cmd:        BINARY_KEY_UNIT_TEMPERATURE,
			PayloadAny: unitValue,
		}

		subPackets = append(subPackets, subPacket)
	}

	if setting.TvocIndex != nil {
		unitValue := byte(99)
		switch setting.TvocIndex.Value {
		case "index":
			unitValue = 1
		case "mg/m³":
			unitValue = 3
		case "ppb":
			unitValue = 4
		}

		if unitValue != 99 {
			subPacket := &TlvData{
				Cmd:        BINARY_KEY_UNIT_TVOC_INDEX,
				PayloadAny: unitValue,
			}

			subPackets = append(subPackets, subPacket)
		}
	}

	return
}

func convertReadingOffset(setting *ReadingOffsetSetting) (subPackets []*TlvData, err error) {
	thOffset := []int16{0, 0}
	thOffsetHasSet := false
	if setting.Temperature != nil {
		offset := int16(setting.Temperature.OffsetValue * 10)
		thOffset[0] = offset
		thOffsetHasSet = true
	}

	if setting.Humidity != nil {
		offset := int16(setting.Humidity.OffsetValue * 10)
		thOffset[1] = offset
		thOffsetHasSet = true
	}

	if thOffsetHasSet {
		subPacket := &TlvData{
			Cmd:        0x2F,
			PayloadAny: thOffset,
		}
		subPackets = append(subPackets, subPacket)
	}

	if setting.Co2 != nil {
		offset := int16(setting.Co2.OffsetPercent * 10)
		subPacket := &TlvData{
			Cmd:        0x3F,
			PayloadAny: offset,
		}
		subPackets = append(subPackets, subPacket)
	}

	return
}

func convertWifiInfo(wifiInfo *WifiInfo) (subPacket *TlvData, err error) {
	value := fmt.Sprintf(`"%s","%s"`, wifiInfo.Ssid, wifiInfo.Password)
	subPacket = &TlvData{
		Cmd:        BINARY_KEY_WIFI,
		PayloadAny: value,
	}
	return
}

func tlvEncodeDict(cmd byte, in map[int]interface{}) (out []byte, err error) {
	size := 0
	encodeList := make([][]byte, 0, len(in))

	keys := sortMapKey(in)
	for _, k := range keys {
		v := in[k]

		encodePart, ierr := encodeKv(k, v)
		if ierr != nil {
			err = ierr
			return
		}

		size += len(encodePart)
		encodeList = append(encodeList, encodePart)
	}

	sop := []byte{0x43, 0x47} // 固定头
	sizeByte, _ := encodeUint16(uint16(size))

	out = make([]byte, 0, 2+1+len(sizeByte)+size)
	out = append(out, sop...)
	out = append(out, cmd)
	out = append(out, sizeByte...)

	for _, encodePart := range encodeList {
		out = append(out, encodePart...)
	}

	out = addCrc(out)
	return
}

func TlvEncode(msg *MessagePod) (out []byte, err error) {
	var packets []*TlvData

	if msg.NeedAck != nil {
		subPacket := &TlvData{
			Cmd:        0x6A,
			PayloadAny: uint8(msg.GetNeedAck()),
		}
		packets = append(packets, subPacket)
	}

	if msg.Debug != nil {
		subPacket := &TlvData{
			Cmd:        0x21,
			PayloadAny: uint8(msg.GetDebug()),
		}
		packets = append(packets, subPacket)
	}

	if interval := msg.IntervalSetting; interval != nil {
		if interval.CollectInterval != 0 {
			subPacket := &TlvData{
				Cmd:        BINARY_KEY_DATA_INTERVAL,
				PayloadAny: uint16(interval.CollectInterval),
			}
			packets = append(packets, subPacket)
		}

		if interval.ReportInterval != 0 {
			subPacket := &TlvData{
				Cmd:        BINARY_KEY_REPORT_INTERVAL,
				PayloadAny: uint16(interval.ReportInterval / 60),
			}
			packets = append(packets, subPacket)
		}

		if interval.BleInterval != 0 {
			subPacket := &TlvData{
				Cmd:        BINARY_KEY_BLE_INTERVAL,
				PayloadAny: uint16(interval.BleInterval),
			}
			packets = append(packets, subPacket)
		}
	}

	if msg.ProductId > 0 {
		subPacket := &TlvData{
			Cmd:        0x38,
			PayloadAny: uint16(msg.ProductId),
		}
		packets = append(packets, subPacket)
	}

	if msg.EndFlag != nil {
		subPacket := &TlvData{
			Cmd:        BINARY_KEY_END_FLAG,
			PayloadAny: uint8(*(msg.EndFlag)),
		}
		packets = append(packets, subPacket)
	}

	for _, alert := range msg.AlertSetting {
		if IsFrogS(msg.ProductId) {
			subPacket, _ := convertFrogsAlert(alert)
			packets = append(packets, subPacket)
			continue
		}

		subPacket, _ := convertAlert(alert, msg.HasBuzzer)
		packets = append(packets, subPacket)
	}

	if msg.Co2Setting != nil {
		subPackets, _ := convertSensorSetting("co2", msg.Co2Setting)
		packets = append(packets, subPackets...)
	}

	/*for metric, setting := range msg.SensorSetting {
		subPackets, _ := convertSensorSetting(metric, setting)
		packets = append(packets, subPackets...)
	}*/

	// 读数分级标准
	for metric, setting := range msg.SensorDataLevel {
		subPacket, _ := convertDataLevelSetting(metric, setting)
		packets = append(packets, subPacket)
	}

	if msg.BatterySetting != nil {
		subPacket, _ := convertBatterySetting(msg.BatterySetting)
		packets = append(packets, subPacket)
	}

	if msg.RealtimeDataDuration > 0 {
		subPacket := &TlvData{
			Cmd:        0x42,
			PayloadAny: uint16(msg.RealtimeDataDuration),
		}
		packets = append(packets, subPacket)
	}

	if msg.MqttSetting != nil {
		subPacket, _ := convertMqttSetting(msg.MqttSetting)
		packets = append(packets, subPacket)
	}

	if msg.CmdType == 0 {
		msg.CmdType = 50
	}

	if msg.UnitSetting != nil {
		subPackets, _ := convertUnitSetting(msg.UnitSetting)
		packets = append(packets, subPackets...)
	}

	if msg.ReadingOffsetSetting != nil {
		subPackets, _ := convertReadingOffset(msg.ReadingOffsetSetting)
		packets = append(packets, subPackets...)
	}

	if msg.WifiInfo != nil {
		subPacket, _ := convertWifiInfo(msg.WifiInfo)
		packets = append(packets, subPacket)
	}

	return tlvEncodePackets(byte(msg.CmdType), packets)
}

func tlvEncodePackets(cmd byte, packets []*TlvData) (out []byte, err error) {
	size := 0
	preEncode := make([]byte, 0)

	for _, subPacket := range packets {
		key := subPacket.Cmd
		value := subPacket.PayloadAny

		encodePart, ierr := encodeKv(key, value)
		if ierr != nil {
			fmt.Printf("encodeKv failed: %d,%v,%s\n", key, value, ierr)
			continue
		}

		size += len(encodePart)
		preEncode = append(preEncode, encodePart...)
	}

	sop := []byte{0x43, 0x47} // 固定头
	sizeByte, _ := encodeUint16(uint16(size))

	out = make([]byte, 0, 2+1+len(sizeByte)+size)
	out = append(out, sop...)
	out = append(out, cmd)
	out = append(out, sizeByte...)
	out = append(out, preEncode...)
	out = addCrc(out)
	return
}
