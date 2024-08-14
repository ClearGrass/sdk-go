package encoding

type AlertSetting struct {
	Valid     bool    `json:"valid,omitempty"`
	Metric    string  `json:"metric"`
	Operator  string  `json:"operator"`
	Value     float32 `json:"value"`
	StartTime int     `json:"startTime"`
	EndTime   int     `json:"endTime"`
	WorkTime  int     `json:"workTime"`
}

type IntervalSetting struct {
	ReportInterval  int `json:"reportInterval,omitempty"`
	CollectInterval int `json:"collectInterval,omitempty"`
	BleInterval     int `json:"bleInterval,omitempty"`
}

type MqttSetting struct {
	Host      string `json:"host,omitempty"`
	Port      string `json:"port,omitempty"`
	User      string `json:"user,omitempty"`
	Password  string `json:"password,omitempty"`
	ClientId  string `json:"clientId,omitempty"`
	UpTopic   string `json:"upTopic,omitempty"`
	DownTopic string `json:"downTopic,omitempty"`
	Value     string `json:"value,omitempty"`
}

type FirmwareInfo struct {
	Version       string `json:"version,omitempty"`
	McuVersion    string `json:"mcuVersion,omitempty"`
	ModuleVersion string `json:"moduleVersion,omitempty"`
}

type IntEffectiveVal struct {
	Value     int  `json:"value"`
	Effective bool `json:"effective,omitempty"`
}

type BoolEffectiveVal struct {
	Value     bool `json:"value,omitempty"`
	Effective bool `json:"effective,omitempty"`
}

type StringEffectiveVal struct {
	Value     string `json:"value,omitempty"`
	Effective bool   `json:"effective,omitempty"`
}

type DataLevelVal struct {
	Value     []float32 `json:"value,omitempty"`
	Effective bool      `json:"effective,omitempty"`
}

type SensorMetricSetting struct {
	MetricName string           `json:"-"`
	AscOpen    BoolEffectiveVal `json:"asc_open,omitempty"`
	Reset      BoolEffectiveVal `json:"reset,omitempty"`

	SensorInterval struct {
		Collect   int  `json:"collect,omitempty"`
		Effective bool `json:"effective,omitempty"`
	} `json:"sensor_interval,omitempty"`

	DataLevel struct {
		Max       float32 `json:"max"`
		Min       float32 `json:"min"`
		Effective bool    `json:"effective,omitempty"`
	} `json:"data_level,omitempty"`
}

type BatterySetting struct {
	DischargeShutdownTime IntEffectiveVal `json:"discharge_shutdown_time"` // 使用电池时的自动关机时间
	AlertValue            IntEffectiveVal `json:"alert_value"`
}

type UnitSetting struct {
	Temperature *StringEffectiveVal `json:"temperature,omitempty"`
	TvocIndex   *StringEffectiveVal `json:"tvoc_index,omitempty"`
}

type ReadingOffset struct {
	OffsetValue   float64 `json:"offsetValue,omitempty"`
	OffsetPercent float64 `json:"offsetPercent,omitempty"`
	Effective     bool    `json:"effective,omitempty"`
}

type ReadingOffsetSetting struct {
	Temperature *ReadingOffset `json:"temperature,omitempty"`
	Humidity    *ReadingOffset `json:"humidity,omitempty"`
	Co2         *ReadingOffset `json:"co2,omitempty"`
}

type WifiInfo struct {
	Ssid     string `json:"ssid,omitempty"`     // wifi名
	Password string `json:"password,omitempty"` // wifi密码
	Desc     string `json:"desc,omitempty"`     // 解析时候懒得解析 将名字和密码都放在此处
}

type NtpSetting struct { // NTP 对时服务
	Host string `json:"host"`
}

type MessagePod struct {
	CmdType              int                      `json:"command"`                        // 命令类似
	Mac                  string                   `json:"mac,omitempty"`                  // 设备mac
	ProductId            int                      `json:"productId,omitempty"`            // 产品id
	HasBuzzer            bool                     `json:"hasBuzzer,omitempty"`            // 有蜂鸣器
	NeedAck              *int                     `json:"needAck,omitempty"`              // 设备上报数据需要服务器应答
	EndFlag              *int                     `json:"endFlag,omitempty"`              // 结尾消息
	UsbPlugin            *int                     `json:"usbPlugin,omitempty"`            // usb 是否插入
	PmSn                 string                   `json:"pmSn,omitempty"`                 // pm传感器序列号
	AlertSetting         []*AlertSetting          `json:"alertSetting,omitempty"`         // 报警配置
	IntervalSetting      *IntervalSetting         `json:"intervalSetting,omitempty"`      // 数据间隔配置
	MqttSetting          *MqttSetting             `json:"mqttSetting,omitempty"`          // 联网配置
	Realtime             *SensorData              `json:"realtime,omitempty"`             // 实时数据
	History              []*SensorData            `json:"history,omitempty"`              // 历史数据
	FirmwareInfo         *FirmwareInfo            `json:"firmwareInfo,omitempty"`         // 固件信息
	Co2Setting           *SensorMetricSetting     `json:"co2Setting,omitempty"`           // 临时先用这个，等增加了其他传感器用下面那个
	SensorDataLevel      map[string]*DataLevelVal `json:"sensorDataLevel,omitempty"`      // 传感器读数标准
	BatterySetting       *BatterySetting          `json:"batterySetting,omitempty"`       // 电池设置
	UnitSetting          *UnitSetting             `json:"unitSetting,omitempty"`          // 单位设置
	RealtimeDataDuration int                      `json:"realtimeDataDuration,omitempty"` //
	WifiInfo             *WifiInfo                `json:"wifiInfo,omitempty"`             //
	Other                map[string]string        `json:"other,omitempty"`                // 其他配置
	ReadingOffsetSetting *ReadingOffsetSetting    `json:"readingOffsetSetting,omitempty"` // 读数偏移
	Debug                *int                     `json:"debug,omitempty"`
	//SensorSetting        map[string]*SensorMetricSetting `json:"sensorSetting,omitempty"`        // sensor 传感器设置
}

func NewMessagePod(fieldInit bool) *MessagePod {
	out := &MessagePod{
		Other:           make(map[string]string),
		SensorDataLevel: make(map[string]*DataLevelVal),
		//SensorSetting:   make(map[string]*SensorMetricSetting),
	}

	if fieldInit {
		out.IntervalSetting = &IntervalSetting{}
	}

	return out
}

func (m *MessagePod) SetEndFlag(val int) {
	m.EndFlag = &val
}

func (m *MessagePod) SetDebug(val int) {
	m.Debug = &val
}

func (m *MessagePod) GetDebug() int {
	if m.Debug == nil {
		return 0 // 无效值
	}
	return *m.Debug
}

func (m *MessagePod) GetEndFlag() int {
	if m.EndFlag == nil {
		return 99999 // 无效值
	}
	return *m.EndFlag
}

func (m *MessagePod) SetNeedAck(val int) {
	m.NeedAck = &val
}

func (m *MessagePod) GetNeedAck() int {
	if m.NeedAck == nil {
		return 9 // 无效值
	}
	return *m.NeedAck
}

type SensorData struct {
	Mac             string   `json:"mac,omitempty"`              //
	Type            string   `json:"type,omitempty"`             //
	Battery         *int64   `json:"battery,omitempty"`          //
	Temperature     *float64 `json:"temperature,omitempty"`      //
	ProbTemperature *float64 `json:"prob_temperature,omitempty"` //
	Humidity        *float64 `json:"humidity,omitempty"`         //
	ProbHumidity    *float64 `json:"prob_humidity,omitempty"`    //
	Pressure        *float64 `json:"pressure,omitempty"`         //
	Co2             *float64 `json:"co2,omitempty"`              //
	Co2Percent      *float64 `json:"co2_percent,omitempty"`      //
	Pm25            *float64 `json:"pm25,omitempty"`             // pm2.5
	Pm10            *float64 `json:"pm10,omitempty"`             // pm10
	Tvoc            *float64 `json:"tvoc,omitempty"`             // tvoc
	Noise           *float64 `json:"noise,omitempty"`            // 噪音
	Lumen           *float64 `json:"lumen,omitempty"`            // 光照
	Rssi            *int     `json:"rssi,omitempty"`             // 信号
	Timestamp       int64    `json:"timestamp,omitempty"`        // 时间戳
	Time            string   `json:"time,omitempty"`             // 字符集类型的时间
}

type TlvData struct {
	Sop         string
	Cmd         int
	Length      int
	Payload     string
	PayloadByte []byte
	PayloadAny  interface{}
	Checksum    int
}

// GetDeviceSecret 根据mac 获取秘钥
type GetDeviceSecret func(string) string

func GetDefaultDeviceSecret(mac string) string {
	return "CF64060BDCF33F15A4E9166F7778CFE4"
}

func IsRobb(productId int) bool {
	switch productId {
	case 0x34:
		return true
	case 0x35:
		return true
	case 0x3a:
		return true
	case 0x3b:
		return true
	}

	return false
}

func IsFrogS(productId int) bool {
	switch productId {
	case 0x3C:
		return true
	case 0x3D:
		return true
	case 0x3E:
		return true
	}

	return false
}

func IsPheasantCo2(productId int) bool {
	switch productId {
	case 0x33:
		return true
	case 0x36:
		return true
	case 0x37:
		return true
	}

	return false
}

func tlvGetMetricByCmd(data *TlvData) string {
	keyToMetric := map[int]string{
		0x3C: "co2",
		0x4F: "temperature",
		0x50: "humidity",
		0x51: "pm25",
		0x52: "pm10",
		0x53: "tvoc",
		0x54: "noise",
		0x55: "lumen",
		0x56: "pressure",
	}

	return keyToMetric[data.Cmd]
}

func tlvGetSensorDataLevelCmdByMetric(metric string) int {
	metricToCmd := map[string]int{
		"co2":         0x3C,
		"temperature": 0x4F,
		"humidity":    0x50,
		"pm25":        0x51,
		"pm10":        0x52,
		"tvoc":        0x53,
		"tvoc_index":  0x53,
		"noise":       0x54,
		"lumen":       0x55,
		"pressure":    0x56,
	}

	return metricToCmd[metric]
}
