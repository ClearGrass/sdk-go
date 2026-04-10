import base64
from datetime import datetime

def fmtTimestamp(timestamp):
    dt = datetime.fromtimestamp(timestamp)
    return dt.strftime("%Y-%m-%d %H:%M:%S")

# Little endian byte to integer
def bytesToIntLittleEndian(byteArray):
    val = 0;
    for i in range(len(byteArray)): 
        val = val | byteArray[i] << i*8
    return val

def crc_check_with_length(data: bytes) -> bool:
    # 长度检查
    length = len(data)
    if length < 2 :
        return False
    
    # 读取小端 uint16 校验和
    crc = bytesToIntLittleEndian(data[length-2 : length])
    
    # 计算累加和
    sum_val = 0
    for v in data[:length-2]:
        sum_val += v
    
    # 返回校验结果
    return sum_val == crc

# 解析以4347开头的hex数据
def tlvUnpack(byteArray):
    cmd = byteArray[2:3].hex()
    length = bytesToIntLittleEndian(byteArray[3:5])
    payload = byteArray[5:5+length]
    productId = 0

    index = 0
    subPackList = [];
    while index < length:
        key = payload[index:index+1].hex()
        subLen = bytesToIntLittleEndian(payload[index+1:index+3])
        subPayload = payload[index+3:index+3+subLen]
        index = index+3+subLen
        subPack = {
            'key': key,
            'len': subLen,
            'payload': subPayload,
        }

        if key == '38':
            productId = subPayload[0]

        #print('key:',key,'payload:',subPayload.hex())
        subPackList.append(subPack)

    return {
        'cmd': cmd,
        'productId': productId,
        'length': length,
        'subPackList': subPackList,
    }

def decodeTHdata(byteArray):
    th = bytesToIntLittleEndian(byteArray[0:3])
    temperature = ((th >> 12) - 500) / 10
    humidity = (th & 0xFFF) / 10
    pressure = bytesToIntLittleEndian(byteArray[3:5])
    battery = byteArray[5]

    return {
        'dataType': 'data',
        'timestamp': 0,
        'time': '',
        'temperature': temperature,
        'humidity': humidity,
        'pressure': pressure,
        'battery': battery,
    }

def decodeRealTimeData(byteArray,productId = 0):
    timestamp = bytesToIntLittleEndian(byteArray[0:4])
    realtimeData = decodeTHdata(byteArray[4:])
    rssi = byteArray[10]
    if rssi >= 128:
        rssi = rssi - 256
    
    realtimeData['dataType'] = 'event'
    realtimeData['timestamp'] = timestamp
    realtimeData['time'] = fmtTimestamp(timestamp)

    realtimeData['rssi'] = rssi
    return realtimeData

def decodeHistoryData(byteArray,productId = 0):
    timestamp = bytesToIntLittleEndian(byteArray[0:4])
    duration = bytesToIntLittleEndian(byteArray[4:6])

    historyDataList = []
    packLen = 6
    index = 6
    i = 0
    while index < len(byteArray):
        historyPack = byteArray[index:index+packLen]
        historyData = decodeTHdata(historyPack)
        historyData['timestamp'] = timestamp + duration * i
        historyData['dataType'] = 'data'
        historyData['time'] = fmtTimestamp(historyData['timestamp'])

        historyDataList.append(historyData)
        i += 1
        index += packLen

    return historyDataList

def decodeSensorDataV2(byteArray):
    sensorData = {}
    timestamp = bytesToIntLittleEndian(byteArray[0:4])
    sensorData['timestamp'] = timestamp

    if byteArray[4] == 1:
        temperatureVal = bytesToIntLittleEndian(byteArray[5:7])
        humidityVal = bytesToIntLittleEndian(byteArray[7:9])
        sensorData['temperature'] = temperatureVal/10.0
        sensorData['humidity'] = humidityVal/10.0
    elif byteArray[4] == 2:
        temperatureVal = bytesToIntLittleEndian(byteArray[5:7])
        sensorData['temperature'] = temperatureVal/10.0
    elif byteArray[4] == 3:
        temperatureVal = bytesToIntLittleEndian(byteArray[5:7])
        humidityVal = bytesToIntLittleEndian(byteArray[7:9])
        pressureVal = bytesToIntLittleEndian(byteArray[9:11])
        sensorData['temperature'] = temperatureVal/10.0
        sensorData['humidity'] = humidityVal/10.0
        sensorData['pressure'] = pressureVal
    elif byteArray[4] == 4:
        temperatureVal = bytesToIntLittleEndian(byteArray[5:7])
        humidityVal = bytesToIntLittleEndian(byteArray[7:9])
        co2Val = bytesToIntLittleEndian(byteArray[9:11])
        sensorData['temperature'] = temperatureVal/10.0
        sensorData['humidity'] = humidityVal/10.0
        sensorData['co2'] = co2Val
    elif byteArray[4] == 10:
        temperatureVal = bytesToIntLittleEndian(byteArray[5:7])
        humidityVal = bytesToIntLittleEndian(byteArray[7:9])
        co2Val = bytesToIntLittleEndian(byteArray[9:11]);
        pm25Val = bytesToIntLittleEndian(byteArray[11:13]);
        pm10Val = bytesToIntLittleEndian(byteArray[13:15]);
        tvocVal = bytesToIntLittleEndian(byteArray[15:17]);
        noiseVal = bytesToIntLittleEndian(byteArray[17:19]);
        lightVal = bytesToIntLittleEndian(byteArray[19:23]);
        sensorData['temperature'] = temperatureVal/10.0
        sensorData['humidity'] = humidityVal/10.0
        sensorData['pm25'] = pm25Val
        sensorData['pm10'] = pm10Val
        sensorData['tvoc'] = tvocVal
        sensorData['noise'] = noiseVal
        sensorData['light'] = lightVal
    return sensorData

def escapePacket(byteArray):
    return packetBytesReplace(byteArray)

def packetBytesReplace(byteArray): 
    if byteArray is None or len(byteArray) < 3:
        return byteArray

    # 判断前3个字节是否匹配
    is_match = (
        byteArray[0] == 0x27 and
        byteArray[1] == 0x03 and
        byteArray[2] == 0x00
    )
    if not is_match:
        return byteArray

    replace_bytes = [0x1A, 0x1B, 0x08]

    # 构建替换映射
    replace_map = {}
    for i in range(3):
        source_byte = byteArray[3 + i]  # 第4-6位
        if source_byte == 0x43:
            continue
        target_byte = replace_bytes[i]
        replace_map[source_byte] = target_byte

    # 处理剩余字节（从第7位开始）
    processed_bytes = bytearray(len(byteArray) - 6)

    for i in range(6, len(byteArray)):
        current_byte = byteArray[i]
        processed_bytes[i - 6] = replace_map.get(current_byte, current_byte)

    return bytes(processed_bytes)


def tlvDecode(byteArray):
    print("origin:", byteArray.hex())
    byteArray = escapePacket(byteArray)

    if (not crc_check_with_length(byteArray)):
        raise("crc check faild")

    unpackData = tlvUnpack(byteArray)
    print("escape:", byteArray.hex(), " ", unpackData["length"]+7 )
    
    outData = {'cmd': unpackData['cmd'], 'productId': unpackData['productId']}
    dataList = []
    
    for subPack in unpackData['subPackList']: 
        print({"key":subPack['key'],"value":subPack['payload'].hex()})

        # if subPack['key'] == '14':
        #     realtimeData = decodeRealTimeData(subPack['payload'], unpackData['productId'])
        #     outData['sensorData'] = [realtimeData];
        
        # if subPack['key'] == '03':
        #     historyeData = decodeHistoryData(subPack['payload'], unpackData['productId'])
        #     outData['sensorData'] = historyeData;

        # if subPack['key'] == '11':
        #     outData['version'] = subPack['payload'].decode("utf-8")
        
        # if subPack['key'] == '34':
        #     outData['versionModel'] = subPack['payload'].decode("utf-8")

        # if subPack['key'] == '35':
        #     outData['versionMcu'] = subPack['payload'].decode("utf-8")

        # if subPack['key'] == '04':
        #     reportInterval = bytesToIntLittleEndian(subPack['payload'])
        #     outData['reportInterval'] = reportInterval
        # if subPack['key'] == '05':
        #     collectInterval = bytesToIntLittleEndian(subPack['payload'])
        #     outData['collectInterval'] = collectInterval
        # if subPack['key'] == '85':
        #     sensorData = decodeSensorDataV2(subPack['payload'])
        #     dataList.append(sensorData)

    if len(dataList) > 0:
        outData['sensorData'] = dataList
    
    return outData

# test
if __name__ == '__main__':
    # decode hex data

    srcList = [
        '4347311c0038020032001d0100011410007124a35ee200ffff2f02e01ccd000000f607',
        '2703004343034347311c0038020032001d010001141000ae06a45ec200ffff9f00fb1ccd0000007f03',
        '2703004343034347311c0038020032001d010001141000fd13a45ec800ffff9f00ec1ccd000000d203',
        '2703004303434347311c0038020032001d0100011410000703a45ecc00ffff9f00b31ccd000000af07'
    ]

    for src in srcList:
        bs = bytes.fromhex(src)
        out = tlvDecode(bs)
        print(out)
